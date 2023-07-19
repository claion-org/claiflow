package p8s

import (
	"fmt"
	"strings"
)

var supportedMethodList = []string{
	"prometheus.query.v1",
	"prometheus.query_range.v1",
	"prometheus.alerts.v1",
	"prometheus.rules.v1",
	"prometheus.alertmanagers.v1",
	"prometheus.targets.v1",
	"prometheus.targets/metadata.v1",
}

func InitCommandFuncs() (map[string]func(map[string]interface{}) ([]byte, error), error) {
	funcList := make(map[string]func(map[string]interface{}) ([]byte, error))

	for _, method := range supportedMethodList {
		mlist := strings.SplitN(method, ".", 3)

		if len(mlist) != 3 {
			return nil, fmt.Errorf("there is not enough method(%s) for p8s. want(3) but got(%d)", method, len(mlist))
		}

		api := mlist[1]
		apiVersion := mlist[2]

		fn := func(api, apiVersion string) func(map[string]interface{}) ([]byte, error) {
			return func(args map[string]interface{}) ([]byte, error) {
				url, ok := args["url"]
				if !ok {
					return nil, fmt.Errorf("prometheus url is empty")
				}

				urlStr, ok := url.(string)
				if !ok {
					return nil, fmt.Errorf("url type must be string, not %T", url)
				}
				if len(urlStr) == 0 {
					return nil, fmt.Errorf("prometheus url is empty")
				}

				c, err := NewClient(urlStr)
				if err != nil {
					return nil, err
				}
				res, err := c.ApiRequest(apiVersion, api, args)
				return []byte(res), err
			}
		}

		funcList[method] = fn(api, apiVersion)
	}

	return funcList, nil
}
