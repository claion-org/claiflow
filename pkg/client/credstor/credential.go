package credstor

import (
	"context"
	"fmt"

	"github.com/claion-org/claiflow/pkg/client/credstor/credential"
)

func addCredential(cc *credential.Client, params map[string]interface{}) ([]*credential.CredentialInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if params == nil && len(params) <= 0 {
		return nil, fmt.Errorf("params is empty")
	}

	credentialsInf, ok := params["credentials"]
	if !ok {
		return nil, fmt.Errorf("credentials argument is empty")
	}

	credentials, ok := credentialsInf.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type assertion for credentials argument")
	}

	if len(credentials) <= 0 {
		return nil, fmt.Errorf("credentials is empty")
	}

	var result []*credential.CredentialInfo
	var creds []*credential.CredentialInfo
	for key, data := range credentials {
		dataM, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed type assertion for credential %q's data: want be map[string]interface{}, not %T", key, data)
		}

		var storageType string
		if storageTypeInf, ok := dataM["storage"]; ok {
			storageType, ok = storageTypeInf.(string)
			if !ok {
				return nil, fmt.Errorf("failed type assertion for credential.storage's data: want be string, not %T", storageTypeInf)
			}
		}

		cred := &credential.CredentialInfo{Key: key, Storage: storageType}
		if typeInf, ok := dataM["type"]; !ok {
			return nil, fmt.Errorf("credential type is empty")
		} else {
			if typeStr, ok := typeInf.(string); !ok {
				return nil, fmt.Errorf("credential type must be string, not %T", typeInf)
			} else {
				cred.Type = typeStr
			}
		}
		if dataInf, ok := dataM["data"]; !ok {
			return nil, fmt.Errorf("credential data is empty")
		} else {
			cred.Data = dataInf
		}

		creds = append(creds, cred)
		result = append(result, &credential.CredentialInfo{Key: cred.Key, Storage: cred.Storage})
	}

	if err := cc.Add(ctx, creds); err != nil {
		return nil, err
	}

	return result, nil
}

func getCredential(cc *credential.Client, params map[string]interface{}) ([]*credential.CredentialInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return cc.GetAll(ctx, false)
}

func updateCredential(cc *credential.Client, params map[string]interface{}) ([]*credential.CredentialInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if params == nil && len(params) <= 0 {
		return nil, fmt.Errorf("params is empty")
	}

	credentialsInf, ok := params["credentials"]
	if !ok {
		return nil, fmt.Errorf("credentials argument is empty")
	}

	credentials, ok := credentialsInf.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type assertion for credentials argument")
	}

	if len(credentials) <= 0 {
		return nil, fmt.Errorf("credentials is empty")
	}

	var result []*credential.CredentialInfo
	var creds []*credential.CredentialInfo
	for key, data := range credentials {
		dataM, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed type assertion for credential %q's data: want be map[string]interface{}, not %T", key, data)
		}

		var storageType string
		if storageTypeInf, ok := dataM["storage"]; ok {
			storageType, ok = storageTypeInf.(string)
			if !ok {
				return nil, fmt.Errorf("failed type assertion for credential.storage's data: want be string, not %T", storageTypeInf)
			}
		}

		cred := &credential.CredentialInfo{Key: key, Storage: storageType}
		if typeInf, ok := dataM["type"]; !ok {
			return nil, fmt.Errorf("credential type is empty")
		} else {
			if typeStr, ok := typeInf.(string); !ok {
				return nil, fmt.Errorf("credential type must be string, not %T", typeInf)
			} else {
				cred.Type = typeStr
			}
		}
		if dataInf, ok := dataM["data"]; !ok {
			return nil, fmt.Errorf("credential data is empty")
		} else {
			cred.Data = dataInf
		}

		creds = append(creds, cred)
		result = append(result, &credential.CredentialInfo{Key: cred.Key, Storage: cred.Storage})
	}

	if err := cc.Update(ctx, creds); err != nil {
		return nil, err
	}

	return result, nil
}

func removeCredential(cc *credential.Client, params map[string]interface{}) ([]*credential.CredentialInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if params == nil && len(params) <= 0 {
		return nil, fmt.Errorf("params is empty")
	}

	credentialsInf, ok := params["credential_keys"]
	if !ok {
		return nil, fmt.Errorf("credential_keys argument is empty")
	}

	credentials, ok := credentialsInf.([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type assertion for credential_keys argument")
	}

	if len(credentials) <= 0 {
		return nil, fmt.Errorf("credential_keys is empty")
	}

	var creds []*credential.CredentialInfo
	for _, data := range credentials {
		dataM, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed type assertion for credential's data: want be map[string]interface{}, not %T", data)
		}

		var key string
		if keyInf, ok := dataM["key"]; ok {
			key, ok = keyInf.(string)
			if !ok {
				return nil, fmt.Errorf("failed type assertion for credential.key's data: want be string, not %T", keyInf)
			}
		}

		var storageType string
		if storageTypeInf, ok := dataM["storage"]; ok {
			storageType, ok = storageTypeInf.(string)
			if !ok {
				return nil, fmt.Errorf("failed type assertion for credential.storage's data: want be string, not %T", storageTypeInf)
			}
		}

		cred := &credential.CredentialInfo{Key: key, Storage: storageType}

		creds = append(creds, cred)
	}

	if err := cc.Remove(ctx, creds); err != nil {
		return nil, err
	}

	return creds, nil
}
