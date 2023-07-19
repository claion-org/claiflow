package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/claion-org/claiflow/pkg/client/credstor"
	"github.com/claion-org/claiflow/pkg/client/helm"
	"github.com/claion-org/claiflow/pkg/client/internal"
	"github.com/claion-org/claiflow/pkg/client/internal/fetcher"
	"github.com/claion-org/claiflow/pkg/client/internal/scheduler"
	"github.com/claion-org/claiflow/pkg/client/internal/server"
	"github.com/claion-org/claiflow/pkg/client/k8s"
	"github.com/claion-org/claiflow/pkg/client/log"
	"github.com/claion-org/claiflow/pkg/client/p8s"
	apiclient "github.com/claion-org/claiflow/pkg/server/api/client"
)

const (
	localTargetServer  = "localhost:18099"
	defaultMaxMsgSize  = 1024 * 1024 * 1024 // 1GB
	defaultAuthTimeout = 10 * time.Second
	defaultLogLevel    = "info"
)

type Client struct {
	conn        *grpc.ClientConn
	fetcher     *fetcher.Fetcher
	scheduler   *scheduler.Scheduler
	clusterId   string
	bearerToken string
	version     string
}

type ClientOptions struct {
	TargetServer    string
	ClusterUuid     string
	BearerToken     string
	LogLevel        string
	WorkersCount    int
	FixedFetchLimit int
	Version         string
	ConnOptions     ConnOptions
}

type ConnOptions struct {
	TLS                          *tls.Config
	MaxMsgSize                   int
	EnableKeepAlive              bool
	KeepAliveTime                time.Duration
	KeepAliveTimeout             time.Duration
	KeepAlivePermitWithoutStream bool
	UserDialOptions              []grpc.DialOption
}

func New(opts ClientOptions) (*Client, error) {
	if opts.TargetServer == "" {
		opts.TargetServer = localTargetServer
	}

	if opts.ClusterUuid == "" {
		return nil, fmt.Errorf("cluster_id is required")
	}
	if opts.BearerToken == "" {
		return nil, fmt.Errorf("bearer_token is required")
	}

	version := "dev"
	if opts.Version != "" {
		version = opts.Version
	}

	var credDialOpts grpc.DialOption
	if opts.ConnOptions.TLS != nil {
		credDialOpts = grpc.WithTransportCredentials(credentials.NewTLS(opts.ConnOptions.TLS))
	} else {
		credDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	maxMsgSize := defaultMaxMsgSize
	if opts.ConnOptions.MaxMsgSize != 0 {
		maxMsgSize = opts.ConnOptions.MaxMsgSize
	}

	var dialOpts []grpc.DialOption

	dialOpts = append(dialOpts, credDialOpts)
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxMsgSize)))
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))

	if opts.ConnOptions.EnableKeepAlive {
		dialOpts = append(dialOpts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                opts.ConnOptions.KeepAliveTime,
			Timeout:             opts.ConnOptions.KeepAliveTimeout,
			PermitWithoutStream: opts.ConnOptions.KeepAlivePermitWithoutStream,
		}))
	}

	dialOpts = append(dialOpts, opts.ConnOptions.UserDialOptions...)

	conn, err := grpc.Dial(opts.TargetServer, dialOpts...)
	if err != nil {
		return nil, err
	}

	logLevel := defaultLogLevel
	if logLevel != "" {
		logLevel = opts.LogLevel
	}

	log.New(logLevel)

	sch := scheduler.NewScheduler(opts.WorkersCount)

	fetcher, err := fetcher.NewFetcher(opts.BearerToken, opts.ClusterUuid, sch, server.NewServerAPI(apiclient.NewClientServiceClient(conn)), opts.FixedFetchLimit, version)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:        conn,
		scheduler:   sch,
		fetcher:     fetcher,
		clusterId:   opts.ClusterUuid,
		bearerToken: opts.BearerToken,
		version:     version,
	}

	// register built-in credstor command list
	credstorCmdFuncs, err := credstor.InitCommandFuncs()
	if err != nil {
		return nil, err
	}
	for name, cf := range credstorCmdFuncs {
		if err := client.RegisterCommand(name, cf); err != nil {
			return nil, err
		}
	}

	// register built-in kubernetes command list
	k8sCmdFuncs, err := k8s.InitCommandFuncs()
	if err != nil {
		return nil, err
	}
	for name, cf := range k8sCmdFuncs {
		if err := client.RegisterCommand(name, cf); err != nil {
			return nil, err
		}
	}

	// register built-in prometheus command list
	p8sCmdFuncs, err := p8s.InitCommandFuncs()
	if err != nil {
		return nil, err
	}
	for name, cf := range p8sCmdFuncs {
		if err := client.RegisterCommand(name, cf); err != nil {
			return nil, err
		}
	}

	// register built-in helm command list
	helmCmdFuncs, err := helm.InitCommandFuncs()
	if err != nil {
		return nil, err
	}
	for name, cf := range helmCmdFuncs {
		if err := client.RegisterCommand(name, cf); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) RegisterCommand(name string, cmdFunc interface{}) error {
	return internal.RegisterCommand(name, cmdFunc)
}

func (c *Client) Run() error {
	c.scheduler.Start()

	if err := c.fetcher.HandShake(); err != nil {
		return fmt.Errorf("failed to handshake: error=%v", err)
	}
	log.Infof("successed to handshake\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// polling
	log.Infof("polling start\n")
	c.fetcher.Polling(fCtx)

	select {
	case <-ctx.Done():
		log.Infof("received signal: SIGINT or SIGTERM\n")

		// fetcher polling stop
		cancel()

		// clean up the remaining services before termination
		for {
			<-time.After(time.Second * 3)
			services := c.fetcher.RemainServices()
			if len(services) == 0 {
				break
			}

			buf := bytes.Buffer{}
			buf.WriteString("remain services:")
			for uuid, status := range services {
				buf.WriteString(fmt.Sprintf("\n\tuuid: %s, status: %s", uuid, status.String()))
			}
			log.Infof(buf.String() + "\n")
		}
	case <-c.fetcher.Done():
		log.Infof("received fetcher done")
	}

	return nil
}
