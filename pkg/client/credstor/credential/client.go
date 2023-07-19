package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/claion-org/claiflow/pkg/client/credstor/credential/storage"
	"github.com/claion-org/claiflow/pkg/client/internal/config"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	storages map[storage.StorageType]storage.Storage
}

func getSecretKey() ([]byte, error) {
	if config.ClusterUuid == "" {
		return []byte("claiflow-credential"), nil
	}
	return []byte(config.ClusterUuid), nil
}

func NewClient() (*Client, error) {
	storages := make(map[storage.StorageType]storage.Storage)
	if storKube, err := storage.NewStorageKubernetes(); err == nil {
		storages[storage.StorageType_KubernetesSecret] = storKube
	}
	if storVol, err := storage.NewStorageVolume("", true, getSecretKey); err == nil {
		storages[storage.StorageType_Volume] = storVol
	}

	return &Client{storages: storages}, nil
}

func (c *Client) checkStorage(st storage.StorageType) error {
	if _, ok := c.storages[st]; ok {
		return nil
	}

	switch st {
	case storage.StorageType_KubernetesSecret:
		if storKube, err := storage.NewStorageKubernetes(); err != nil {
			return err
		} else {
			c.storages[storage.StorageType_KubernetesSecret] = storKube
		}
	case storage.StorageType_Volume:
		if storVol, err := storage.NewStorageVolume("", true, getSecretKey); err != nil {
			return err
		} else {
			c.storages[storage.StorageType_Volume] = storVol
		}
	default:
		return fmt.Errorf("not supported storage type: %v", st)
	}

	return nil
}

func (c *Client) Add(ctx context.Context, creds []*CredentialInfo) error {
	storCredM := make(map[storage.StorageType][]*CredentialInfo)

	for _, cred := range creds {
		if cred.Storage == "" {
			storCredM[storage.StorageType_KubernetesSecret] = append(storCredM[storage.StorageType_KubernetesSecret], cred)
		} else {
			storCredM[storage.StorageType(cred.Storage)] = append(storCredM[storage.StorageType(cred.Storage)], cred)
		}
	}

	for storType, creds := range storCredM {
		if err := c.checkStorage(storType); err != nil {
			return err
		}

		dataset := make(map[string][]byte)
		for _, cred := range creds {
			b, err := cred.GetInfoAsBytes()
			if err != nil {
				return err
			}
			dataset[cred.Key] = b
		}

		if err := c.storages[storType].Add(ctx, dataset); err != nil {
			return err
		}
	}

	return nil
}

// func (c *Client) Get(ctx context.Context, cred *CredentialInfo, getData bool) (*CredentialInfo, error) {
func (c *Client) Get(ctx context.Context, st storage.StorageType, key string, getData bool) (*CredentialInfo, error) {
	if err := c.checkStorage(st); err != nil {
		return nil, err
	}

	storage := c.storages[storage.StorageType(st)]

	b, err := storage.Get(ctx, key, getData)
	if err != nil {
		return nil, err
	}

	cred := &CredentialInfo{}
	if getData {
		if len(b) <= 0 {
			return nil, fmt.Errorf("failed to get credential data")
		}

		if err := json.Unmarshal(b, cred); err != nil {
			return nil, err
		}
	} else {
		cred.Storage = string(st)
		cred.Key = key
	}

	return cred, nil
}

func (c *Client) GetAll(ctx context.Context, getData bool) ([]*CredentialInfo, error) {
	var creds []*CredentialInfo

	for storType, stor := range c.storages {
		m, err := stor.GetAll(ctx, getData)
		if err != nil {
			return nil, err
		}

		for key, b := range m {
			cred := &CredentialInfo{}
			if getData {
				if len(b) <= 0 {
					return nil, fmt.Errorf("failed to get credential data")
				}

				if err := json.Unmarshal(b, cred); err != nil {
					return nil, err
				}
			} else {
				cred.Key = key
				cred.Storage = string(storType)
			}
			creds = append(creds, cred)
		}
	}

	return creds, nil
}

func (c *Client) Update(ctx context.Context, creds []*CredentialInfo) error {
	storCredM := make(map[storage.StorageType][]*CredentialInfo)

	for _, cred := range creds {
		if cred.Storage == "" {
			storCredM[storage.StorageType_KubernetesSecret] = append(storCredM[storage.StorageType_KubernetesSecret], cred)
		} else {
			storCredM[storage.StorageType(cred.Storage)] = append(storCredM[storage.StorageType(cred.Storage)], cred)
		}
	}

	for storType, creds := range storCredM {
		if err := c.checkStorage(storType); err != nil {
			return err
		}

		dataset := make(map[string][]byte)
		for _, cred := range creds {
			b, err := cred.GetInfoAsBytes()
			if err != nil {
				return err
			}
			dataset[cred.Key] = b
		}

		if err := c.storages[storType].Update(ctx, dataset); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Remove(ctx context.Context, creds []*CredentialInfo) error {
	storCredM := make(map[storage.StorageType][]*CredentialInfo)

	for _, cred := range creds {
		if cred.Storage == "" {
			storCredM[storage.StorageType_KubernetesSecret] = append(storCredM[storage.StorageType_KubernetesSecret], cred)
		} else {
			storCredM[storage.StorageType(cred.Storage)] = append(storCredM[storage.StorageType(cred.Storage)], cred)
		}
	}

	for storType, creds := range storCredM {
		if err := c.checkStorage(storType); err != nil {
			return err
		}

		var keys []string
		for _, cred := range creds {
			keys = append(keys, cred.Key)
		}

		if err := c.storages[storType].Delete(ctx, keys); err != nil {
			return err
		}
	}

	return nil
}

type CredentialInfo struct {
	Key     string      `json:"key,omitempty"`
	Storage string      `json:"storage,omitempty"`
	Type    string      `json:"type,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (ci *CredentialInfo) GetInfoAsBytes() ([]byte, error) {
	return json.Marshal(ci)
}
