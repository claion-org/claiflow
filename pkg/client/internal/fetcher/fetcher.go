package fetcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/panta/machineid"

	"github.com/claion-org/claiflow/pkg/client/internal/config"
	"github.com/claion-org/claiflow/pkg/client/internal/flow"
	"github.com/claion-org/claiflow/pkg/client/internal/scheduler"
	"github.com/claion-org/claiflow/pkg/client/internal/server"
	"github.com/claion-org/claiflow/pkg/client/internal/service"
	"github.com/claion-org/claiflow/pkg/client/log"
	"github.com/claion-org/claiflow/pkg/server/model"
)

type Fetcher struct {
	bearerToken   string
	machineID     string
	clusterId     string
	serverAPI     server.ServerAPIInterface
	scheduler     *scheduler.Scheduler
	done          chan struct{}
	fetchLimit    int
	clientVersion string
}

func NewFetcher(bearerToken, clusterId string, scheduler *scheduler.Scheduler, api server.ServerAPIInterface, fetchLimit int, clientVersion string) (*Fetcher, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	id = strings.ReplaceAll(id, "-", "")

	if api == nil {
		return nil, fmt.Errorf("server api client is required")
	}

	config.ClusterUuid = clusterId

	return &Fetcher{
		bearerToken:   bearerToken,
		machineID:     id,
		clusterId:     clusterId,
		serverAPI:     api,
		scheduler:     scheduler,
		done:          make(chan struct{}),
		fetchLimit:    fetchLimit,
		clientVersion: clientVersion}, nil
}

func (f *Fetcher) Done() <-chan struct{} {
	return f.done
}

func (f *Fetcher) Cancel() {
	close(f.done)
}

func (f *Fetcher) Polling(ctx context.Context) error {
	go f.UpdateServiceProcess()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugf("polling context done\n")
				return
			case <-f.Done():
				log.Debugf("polling context done\n")
				return
			default:
				if err := f.longPoll(ctx); err != nil {
					log.Debugf("failed to fetcher.run. error=%v", err)
					f.RetryHandshake()
				}
			}
		}
	}()

	return nil
}

func (f *Fetcher) longPoll(ctx context.Context) error {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-cctx.Done():
			return cctx.Err()
		default:
		}

		//  if existing service's status is ServiceStatusSuccess or ServiceStatusFailed, delete in statusMap
		f.scheduler.CleanupRemainingServices()

		var reqChs []chan *service.Service
		var getServicesLimit int
		if f.fetchLimit > 0 {
			reqChs = append(reqChs, f.scheduler.GetAvailableWorkerChannel())
			getServicesLimit = f.fetchLimit
		} else {
			// block, when until get available worker channel
			reqChs = f.scheduler.GetAvailableWorkerChannelList()
			getServicesLimit = len(reqChs)
		}

		// get services on long poll
		respData, err := f.serverAPI.GetServices(cctx, getServicesLimit)
		if err != nil {
			log.Errorf("failed to polling: error=%s\n", err.Error())
			if err == context.DeadlineExceeded {
				for _, ch := range reqChs {
					ch <- nil
				}
				continue
			}
			for _, ch := range reqChs {
				ch <- nil
			}
			return err
		}

		log.Debugf("received services from server: service_count=%d\n", len(respData))

		if len(respData) == 0 {
			for _, ch := range reqChs {
				ch <- nil
			}
			continue
		}

		// respData -> services
		recvServices, failed := service.ConvertServiceListServerToClient(respData)

		if len(failed) > 0 {
			log.Debugf("failed to convert service: failed_service_count=%d\n", len(failed))
			f.UpdateFailedToConvertServices(failed)
		}

		// catch client system service
		if ok := f.CatchClientSystemService(recvServices); ok {
			for _, ch := range reqChs {
				ch <- nil
			}
			continue
		}

		// Register new services.
		f.scheduler.PushServices(reqChs, recvServices)
	}
}

func (f *Fetcher) ChangeClientConfigFromToken() {
	claims := new(model.ClusterClientSessionClaim)
	jwt_token, _, err := jwt.NewParser().ParseUnverified(f.serverAPI.GetToken(), claims)
	if _, ok := jwt_token.Claims.(*model.ClusterClientSessionClaim); !ok || err != nil {
		if err == nil {
			err = fmt.Errorf("unable to convert token.claims to *sessionv1.ClientSessionPayload")
		}
		log.Warnf("failed to bind payload: error=%v\n", err)
		return
	}
}

func (f *Fetcher) UpdateServiceProcess() {
	for update := range f.scheduler.NotifyServiceUpdate() {
		<-time.After(time.Millisecond * 100)

		go func(up *service.UpdateService) {
			if err := f.serverAPI.UpdateServices(context.Background(), up); err != nil {
				log.Errorf("failed to update service on server: service_uuid=%s, error=%v\n", up.GetId(), err)
			}

			f.scheduler.UpdateServiceStatus(up)
		}(update)
	}
}

func (f *Fetcher) CatchClientSystemService(services []*service.Service) bool {
	exist := false

	for _, svc := range services {
		for _, step := range svc.Flow {
			switch sstep := step.(type) {
			case *flow.CommandStep:
				if sstep.Command != "" {
					method := sstep.Command

					switch method {
					case "claiflow.client_pod.rebounce":
						exist = true
						f.RebounceClientPod(svc.Id)
					case "claiflow.client.upgrade":
						exist = true
						f.UpgradeClient(svc.Id, sstep.Inputs.GetInputs())
					}
					if exist {
						return exist
					}
				}
			default:
			}
		}
	}

	return exist
}

func (f *Fetcher) RemainServices() map[string]service.ServiceStatus {
	return f.scheduler.CleanupRemainingServices()
}

func (f *Fetcher) UpdateFailedToConvertServices(failed []*service.UpdateService) {
	for _, d := range failed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := f.serverAPI.UpdateServices(ctx, d); err != nil {
			log.Errorf("failed to update service on server: service_uuid=%s, error=%v\n", d.GetId(), err)
		}
	}
}
