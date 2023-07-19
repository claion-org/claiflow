package echov4

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo/v4"
)

func PathParam(echo_ echo.Context) map[string]string {
	m := map[string]string{}
	for _, name := range echo_.ParamNames() {
		m[name] = echo_.Param(name)
	}

	return m
}

func PathParams(echo_ echo.Context) map[string]string {
	m := map[string]string{}
	for _, name := range echo_.ParamNames() {
		m[name] = echo_.Param(name)
	}

	return m
}

func ParamString(echo_ echo.Context) string {
	param := PathParam(echo_)
	s := make([]string, 0, len(param))
	for key := range param {
		s = append(s, fmt.Sprintf("%s:%s", key, param[key]))
	}
	return strings.Join(s, ",")
}

func QueryParam(echo_ echo.Context) map[string]string {
	m := map[string]string{}
	for key := range echo_.QueryParams() {
		m[key] = echo_.QueryParam(key)
	}

	return m
}

func QueryParamString(echo_ echo.Context) string {
	return echo_.QueryString()
}

func FormParam(echo_ echo.Context) map[string]string {
	m := map[string]string{}
	formdatas, err := echo_.FormParams()
	if err != nil {
		return m
	}

	for key := range formdatas {
		m[key] = echo_.FormValue(key)
	}

	return m
}

func FormParams(echo_ echo.Context, key string) ([]string, bool) {
	formdatas, err := echo_.FormParams()
	if err != nil {
		return nil, false
	}

	vals, ok := formdatas[key]

	return vals, ok
}

func FormParamString(echo_ echo.Context) string {
	formparam := FormParam(echo_)
	s := make([]string, 0, len(formparam))
	for key := range formparam {
		s = append(s, fmt.Sprintf("%s=%s", key, formparam[key]))
	}
	return strings.Join(s, "&")
}

func Drain(ctx echo.Context) ([]byte, error) {
	return io.ReadAll(ctx.Request().Body)
}

func Body(ctx echo.Context) ([]byte, error) {
	var err error

	// create clone
	var a, b bytes.Buffer
	w := io.MultiWriter(&a, &b)
	_, err = io.Copy(w, ctx.Request().Body)
	// check error
	if err != nil {
		return nil, err
	}
	// restore A to preserve
	ctx.Request().Body = ioutil.NopCloser(&a)
	// read all
	return io.ReadAll(&b)
}

func Bind[T any](ctx echo.Context, v *T) error {
	var err error

	// create clone
	var a, b bytes.Buffer
	w := io.MultiWriter(&a, &b)
	_, err = io.Copy(w, ctx.Request().Body)
	// check error
	if err != nil {
		return err
	}
	// restore A to read
	ctx.Request().Body = ioutil.NopCloser(&a)
	// bind
	err = ctx.Bind(v)
	// restore B to preserve
	ctx.Request().Body = ioutil.NopCloser(&b)

	return err
}
