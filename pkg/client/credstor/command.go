package credstor

import (
	"encoding/json"

	"github.com/claion-org/claiflow/pkg/client/credstor/credential"
)

type Handler = func(params map[string]interface{}) ([]byte, error)

func InitCommandFuncs() (map[string]func(map[string]interface{}) ([]byte, error), error) {
	credClient, err := credential.NewClient()
	if err != nil {
		return nil, err
	}

	return map[string]Handler{
		"credstor.credential.add":    WrapFunc(WrapCredential(credClient, addCredential)),
		"credstor.credential.get":    WrapFunc(WrapCredential(credClient, getCredential)),
		"credstor.credential.update": WrapFunc(WrapCredential(credClient, updateCredential)),
		"credstor.credential.remove": WrapFunc(WrapCredential(credClient, removeCredential)),
	}, nil
}

func WrapFunc[A any](fn func(args map[string]interface{}) (A, error)) func(args map[string]interface{}) ([]byte, error) {
	return func(args map[string]interface{}) ([]byte, error) {
		a, err := fn(args)
		if err != nil {
			return nil, err
		}

		return json.Marshal(a)
	}
}

func WrapCredential[A any](cc *credential.Client, fn func(cc *credential.Client, args map[string]interface{}) (A, error)) func(args map[string]interface{}) (A, error) {
	return func(args map[string]interface{}) (A, error) {
		return fn(cc, args)
	}
}
