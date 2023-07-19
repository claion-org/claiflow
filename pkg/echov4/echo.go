package echov4

import "github.com/labstack/echo/v4"

// HttpError
func HttpError(err error, code int) error {
	if err == nil {
		return nil
	}
	return echo.NewHTTPError(code).SetInternal(err)
}

// OK
func OK() interface{} {
	return "OK"
}
