package storage

import "context"

type StorageType string

const (
	StorageType_KubernetesSecret StorageType = "k8s/secret"
	StorageType_Volume           StorageType = "volume"
)

type Storage interface {
	Add(context.Context, map[string][]byte) error
	Get(context.Context, string, bool) ([]byte, error)
	GetAll(context.Context, bool) (map[string][]byte, error)
	Update(context.Context, map[string][]byte) error
	Delete(context.Context, []string) error
}
