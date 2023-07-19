// @title                      CLAIFLOW
// @version                    0.0.1
// @description                this is a claiflow server.
// @contact.url                https://claion.co.kr
// @contact.email              jaehoon@claion.co.kr
// @securityDefinitions.apikey ClientAuthorization
// @in                         header
// @name                       Authorization
// @description                Bearer token for client api

package route

import (
	"fmt"
	"os"
	"sync/atomic"

	middlewareX "github.com/claion-org/claiflow/pkg/echov4/middleware"
	pprof "github.com/claion-org/claiflow/pkg/echov4/pprof"
	"github.com/claion-org/claiflow/pkg/server/api/v1/client"
	"github.com/claion-org/claiflow/pkg/server/api/v1/cluster"
	"github.com/claion-org/claiflow/pkg/server/api/v1/cluster_client_session"
	"github.com/claion-org/claiflow/pkg/server/api/v1/cluster_client_token"
	"github.com/claion-org/claiflow/pkg/server/api/v1/global_variables"
	"github.com/claion-org/claiflow/pkg/server/api/v1/service"
	"github.com/claion-org/claiflow/pkg/server/api/v1/webhook"
	"github.com/claion-org/claiflow/pkg/server/config"
	"github.com/claion-org/claiflow/pkg/server/route/docs"
	"github.com/claion-org/claiflow/pkg/version"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	//swago docs version
	docs.SwaggerInfo.Version = version.Version
}

func Route(cfg config.Config) *echo.Echo {
	e := echo.New()

	// logger
	if true {
		e.Use(middlewareX.ServiceLogger(os.Stdout))
	}

	// CORS
	e.Use(middlewareX.SetCORS(cfg.HttpService.CORS))
	// "X-Request-Id": string
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() func() string {
			var (
				id uint64
			)
			return func() string {
				id := atomic.AddUint64(&id, 1)
				return fmt.Sprintf("%d", id)
			}
		}(),
	}))
	// "Content-Encoding": "gzip"
	e.Use(middleware.Decompress())

	// echo error handler
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		middlewareX.ErrorResponder(err, ctx)
		// middlewareX.ErrorLogger(err, ctx, logger.Logger())
	}
	// echo recover
	e.Use(middlewareX.Recover())

	// prometheus
	p := prometheus.NewPrometheus("echo", nil)
	e.Use(p.HandlerFunc)
	promGroup := e.Group(cfg.HttpService.URLPrefix)
	if false {
		promGroup.Use(
			middleware.Gzip(), // 'Content-Encoding' handler
		)
	}
	promGroup.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// swaggo
	docs.SwaggerInfo.BasePath = cfg.HttpService.URLPrefix
	swaggoGroup := e.Group(cfg.HttpService.URLPrefix)
	swaggoGroup.Use(
		middleware.Gzip(), // 'Content-Encoding' handler
	)
	swaggoGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	// pprof
	pprofGroup := e.Group(cfg.HttpService.URLPrefix)
	pprofGroup.Use(
		middleware.Gzip(), // 'Content-Encoding' handler
	)
	pprof.WrapGroup(pprofGroup)

	// service
	rootGroup := e.Group(cfg.HttpService.URLPrefix)
	rootGroup.Use(
		middleware.Gzip(), // 'Content-Encoding' handler
	)

	apiV1 := rootGroup.Group("/api/v1")

	// global_variables
	apiV1.PUT("/global_variables/:uuid", global_variables.UpdateGlobalVariableValue)
	apiV1.GET("/global_variables", global_variables.FindGlobalVariables)
	apiV1.GET("/global_variables/:uuid", global_variables.GetGlobalVariable)

	// cluster
	apiV1.POST("/cluster", cluster.CreateCluster)
	apiV1.PUT("/cluster/:uuid", cluster.UpdateCluster)
	apiV1.DELETE("/cluster/:uuid", cluster.DeleteCluster)
	apiV1.GET("/cluster", cluster.FindCluster)
	apiV1.GET("/cluster/:uuid", cluster.GetCluster)
	apiV1.GET("/cluster/:cluster_uuid/session/alive", cluster.GetClusterClientSessionAlive)

	// cluster_token
	apiV1.POST("/cluster_token", cluster_client_token.CreateClusterClientToken)
	apiV1.PUT("/cluster_token/:uuid", cluster_client_token.UpdateClusterClientToken)
	apiV1.PUT("/cluster_token/:uuid/refresh", cluster_client_token.RefreshClusterClientToken)
	apiV1.PUT("/cluster_token/:uuid/expire", cluster_client_token.ExpireClusterClientToken)
	apiV1.DELETE("/cluster_token/:uuid", cluster_client_token.DeleteClusterClinetToken)
	apiV1.GET("/cluster_token", cluster_client_token.FindClusterClientToken)
	apiV1.GET("/cluster_token/:uuid", cluster_client_token.GetClusterClientToken)

	// session
	apiV1.DELETE("/session/:uuid", cluster_client_session.DeleteClusterClientSession)
	apiV1.GET("/session", cluster_client_session.FindClusterClientSession)
	apiV1.GET("/session/:uuid", cluster_client_session.GetClusterClientSession)
	apiV1.GET("/session/cluster/:cluster_uuid/alive", cluster_client_session.GetClusterClientSessionAlive)

	// service
	apiV1.POST("/service", service.CreateServiceMultiClusters)
	apiV1.POST("/cluster/:cluster_uuid/service", service.CreateService)
	apiV1.GET("/service", service.FindService)
	apiV1.GET("/cluster/:cluster_uuid/service/:uuid", service.GetService)
	apiV1.GET("/cluster/:cluster_uuid/service/:uuid/result", service.GetServiceResult)

	// client ECHO
	apiV1.POST("/client/auth", client.Echo{}.Auth)
	apiV1.PUT("/client/service", client.Echo{}.UpdateService)
	apiV1.GET("/client/service", client.Echo{}.PollService)

	// webhook
	apiV1.POST("/webhook", webhook.CreateWebhook)
	apiV1.POST("/webhook/:uuid/publish", webhook.Publish)
	apiV1.PUT("/webhook/:uuid", webhook.UpdateWebhook)
	apiV1.DELETE("/webhook/:uuid", webhook.DeleteWebhook)
	apiV1.GET("/webhook", webhook.FindWebhook)
	apiV1.GET("/webhook/:uuid", webhook.GetWebhook)

	return e
}
