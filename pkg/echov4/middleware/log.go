package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/claion-org/claiflow/pkg/echov4"
	"github.com/claion-org/claiflow/pkg/server/macro/logs"
	"github.com/go-logr/logr"
	"github.com/labstack/echo/v4"
)

func ServiceLogger(w io.Writer) echo.MiddlewareFunc {
	//echo logger
	format := fmt.Sprintf("{%v}\n",
		strings.Join([]string{
			`"time":${time_rfc3339_nano}`,
			`"id":${id}`,
			`"remote_ip":${remote_ip}`,
			`"host":${host}`,
			`"method":${method}`,
			`"uri":${uri}`,
			`"status":${status}`,
			`"error":${error}`,
			`"latency":${latency}`,
			`"latency_human":${latency_human}`,
			`"bytes_in":${bytes_in}`,
			`"bytes_out":${bytes_out}`,
		}, ","))

	logconfig := DefaultLoggerConfig
	logconfig.Output = w
	logconfig.Format = format

	return LoggerWithConfig(logconfig)
}

func ErrorResponder(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if httperr, ok := err.(*echo.HTTPError); ok {
		code = httperr.Code
		if httperr.Internal != nil {
			err = httperr.Internal
		}
	}

	ctx.JSON(code, map[string]interface{}{
		"code": code,
		// "status":     http.StatusText(code),
		"message": err.Error(),
	})
}

func ErrorLogger(err error, ctx echo.Context, logger logr.Logger) {
	nullstring := func(p *string) (s string) {
		s = fmt.Sprintf("%v", p)
		if p != nil {
			s = *p
		}
		return
	}

	code := http.StatusInternalServerError
	if httperr, ok := err.(*echo.HTTPError); ok {
		code = httperr.Code
		if httperr.Internal != nil {
			err = httperr.Internal
		}
	}

	var stack *string
	//stack for surface
	logs.StackIter(err, func(s string) {
		stack = &s
	})
	//stack for internal
	logs.CauseIter(err, func(err error) {
		logs.StackIter(err, func(s string) {
			stack = &s
		})
	})

	id := ctx.Response().Header().Get(echo.HeaderXRequestID)

	reqbody, _ := echov4.Body(ctx)

	if stack == nil {
		logger.Error(err, "request with error",
			"id", id,
			"code", code,
			"method", ctx.Request().Method,
			"url", ctx.Request().RequestURI,
			"reqbody", reqbody,
			"bytes_in", ctx.Request().ContentLength,
		)
	} else {
		logger.Error(err, "request with error",
			"id", id,
			"code", code,
			"method", ctx.Request().Method,
			"url", ctx.Request().RequestURI,
			"reqbody", reqbody,
			"bytes_in", ctx.Request().ContentLength,
			"stack", nullstring(stack),
		)
	}
}
