package helm

type HandlerFunc = func(args map[string]interface{}) (string, error)

func InitCommandFuncs() (map[string]func(map[string]interface{}) ([]byte, error), error) {
	return map[string]func(map[string]interface{}) ([]byte, error){
		"helm.install": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.Install
		}),
		"helm.uninstall": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.Uninstall
		}),
		"helm.upgrade": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.Upgrade
		}),
		"helm.get_values": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.GetValues
		}),
		"helm.repo_add": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.RepoAdd
		}),
		"helm.repo_list": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.RepoList
		}),
		"helm.repo_update": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.RepoUpdate
		}),
		"helm.history": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.History
		}),
		"helm.rollback": WrapFunc(NewClient, func(cc *Client) HandlerFunc {
			return cc.Rollback
		}),
	}, nil
}

func WrapFunc(newClient func() (*Client, error), fn func(cc *Client) HandlerFunc) func(args map[string]interface{}) ([]byte, error) {
	return func(args map[string]interface{}) ([]byte, error) {
		cc, err := newClient()
		if err != nil {
			return nil, err
		}

		s, err := fn(cc)(args)
		if err != nil {
			return nil, err
		}

		return []byte(s), nil
	}
}
