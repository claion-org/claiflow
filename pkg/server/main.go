package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/claion-org/claiflow/pkg/cryptography"
	"github.com/claion-org/claiflow/pkg/enigma"
	"github.com/claion-org/claiflow/pkg/logger"
	"github.com/claion-org/claiflow/pkg/server/config"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/database"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	_ "github.com/claion-org/claiflow/pkg/server/database/vanilla/excute/dialects/mysql"
	"github.com/claion-org/claiflow/pkg/server/route"
	"github.com/claion-org/claiflow/pkg/version"
	"github.com/go-logr/logr"
	"github.com/jinzhu/configor"
	"github.com/spf13/cobra"
	"github.com/zenazn/goji/graceful"
)

var rootCmd = &cobra.Command{
	Use:     "claiflow [flags] YAML, [YAML]...",
	Short:   "claiflow is an open-source project for Kubernetes automation.",
	Long:    ``,
	Version: version.BuildVersion("claiflow-API-server"),
	Args:    cobra.MinimumNArgs(1),
	Run:     Run,
}

func init() {
	rootCmd.SetVersionTemplate("{{printf .Version}}\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Run(cmd *cobra.Command, args []string) {
	// log := logger.Logger()

	time.Local = time.UTC //set timezone UTC
	logger.Info("init time",
		"now", time.Now(),
		"unix", time.Now().Unix(),
		"local", time.Local)

	configFiles := args

	_config := config.Config{}
	if err := configor.Load(&_config, configFiles...); err != nil {
		logger.Error(err, "load config",
			"files", configFiles)
		os.Exit(1)
	}

	logger.Logger = func() logr.Logger {
		return logger.NewZapr(
			func(config *logger.ZaprConfig) { config.Verbose = _config.Logger.Verbose },
			func(config *logger.ZaprConfig) { config.DisableCaller = _config.Logger.DisableCaller },
			func(config *logger.ZaprConfig) { config.DisableStacktrace = _config.Logger.DisableStacktrace },
		)
	}

	// log = logger.Logger()

	logger.Info("init logger",
		"logger", _config.Logger)

	// enigma
	machine, err := enigma.NewMachine(_config.Enigma.ToOption())
	if err != nil {
		logger.Error(err, "new enigma machine")
		os.Exit(1)
	}

	if err := TestCryptography(machine); err != nil {
		logger.Error(err, "test enigma machine")
		os.Exit(1)
	}

	cryptography.Cipher = machine

	// database schema migration
	if err := Migrate(_config); err != nil {
		logger.Error(err, "database schema migration")
		return
	}

	// new database connection
	db, err := database.New(_config.Database)
	if err != nil {
		logger.Error(err, "new database connection")
		os.Exit(1)
	}
	defer db.Close()

	control.Database = func() *sql.DB { return db }
	control.Driver = func() excute.SqlExcutor { return excute.GetSqlExcutor(_config.Database.Type) }

	// cron for global variables
	cronGVClose, err := Cron_GlobalVariables(db, excute.GetSqlExcutor(_config.Database.Type))
	if err != nil {
		logger.Error(err, "init cron global variables")
		os.Exit(1)
	}
	defer cronGVClose()

	NewHttpServe := HttpServeFactory(_config)

	grpcL, err := NewGrpcListener(_config)
	if err != nil {
		logger.Error(err, "new listener for grpc service")
		os.Exit(1)
	}
	defer grpcL.Close()

	fmt.Printf("⚡ grpc server listen on %v\n", grpcL.Addr())

	httpS := route.Route(_config)

	grpcS, err := NewGrpcServer(_config, logger.Logger())
	if err != nil {
		logger.Error(err, "new grpc server")
		os.Exit(1)
	}

	graceful.HandleSignals()
	graceful.AddSignal(syscall.SIGTERM)

	graceful.PreHook(func() {
		logger.Info("Server received signal, gracefully stopping.")
	})
	graceful.PreHook(func() {
		// stop grpc service
		grpcS.GracefulStop()
	})
	graceful.PreHook(func() {
		// stop http service
		httpS.Shutdown(context.Background())
	})

	graceful.PostHook(func() {
		logger.Info("Server stopped")
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := grpcS.Serve(grpcL); err != nil {
			logger.Error(err, "grpc service was close with error")
			return
		}

		logger.Info("grpc service was close")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := NewHttpServe(httpS); err != nil {
			logger.Error(err, "http service was close with error")
			return
		}

		logger.Info("htto service was close")
	}()

	graceful.Wait()

	wg.Wait()
}

func WaitSignal(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
}

func TestCryptography(machine *enigma.Machine) error {
	const quickbrownfox = `the quick brown fox jumps over the lazy dog`

	encripted, err := machine.Encode([]byte(quickbrownfox))
	if err != nil {
		return fmt.Errorf("%w: enigma test: encode", err)
	}

	plain, err := machine.Decode(encripted)
	if err != nil {
		return fmt.Errorf("%w: enigma test: decode", err)
	}

	if strings.Compare(quickbrownfox, string(plain)) != 0 {
		return fmt.Errorf("%w: enigma test: diff result", err)
	}

	return nil
}
