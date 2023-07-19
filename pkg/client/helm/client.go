package helm

import (
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/claion-org/claiflow/pkg/client/log"
)

const defaultMaxHistory = 20

type Client struct {
	settings *cli.EnvSettings
}

func NewClient() (*Client, error) {
	settings := cli.New()

	return &Client{settings: settings}, nil
}

func (c *Client) getActionConfig() (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(c.settings.RESTClientGetter(), c.settings.Namespace(), os.Getenv("HELM_DRIVER"), log.HelmDebugf); err != nil {
		return nil, err
	}

	return actionConfig, nil
}
