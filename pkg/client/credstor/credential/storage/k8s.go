package storage

import (
	"context"
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/claion-org/claiflow/pkg/client/k8s"
)

var _ Storage = &StorageKubernetes{}

var (
	storageKubernetesNamespace  = "default"
	storageKubernetesSecretName = "claiflow-credential"
)

func init() {
	// get namespace
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		storageKubernetesNamespace = string(namespace)
	}
}

type StorageKubernetes struct {
	k8sClient *k8s.Client
	// secret    *corev1.Secret
}

func NewStorageKubernetes() (Storage, error) {
	k8sClient, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}

	return &StorageKubernetes{k8sClient: k8sClient}, nil
}

func (s *StorageKubernetes) Add(ctx context.Context, dataset map[string][]byte) (err error) {
	found := true
	secret, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Get(ctx, storageKubernetesSecretName, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		found = false
		secret = newSecretForCredStor()
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	for key, data := range dataset {
		if errs := validation.IsConfigMapKey(key); len(errs) != 0 {
			return fmt.Errorf("credential %q's key is invalid: %s", key, strings.Join(errs, ";"))
		}
		if _, ok := secret.Data[key]; ok {
			return fmt.Errorf("credential %q already exists", key)
		}

		secret.Data[key] = data
	}

	if !found {
		// k8s create secret
		_, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		// k8s update secret
		_, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StorageKubernetes) Get(ctx context.Context, key string, getData bool) ([]byte, error) {
	secret, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Get(ctx, storageKubernetesSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []byte
	if secret != nil && secret.Data != nil {
		b, ok := secret.Data[key]
		if !ok {
			return nil, fmt.Errorf("%q key's data isn't exist", key)
		}

		if getData {
			result = b
		}
	}

	return result, nil
}

func (s *StorageKubernetes) GetAll(ctx context.Context, getData bool) (map[string][]byte, error) {
	secret, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Get(ctx, storageKubernetesSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	results := make(map[string][]byte)
	if secret != nil && secret.Data != nil {
		for key, data := range secret.Data {
			if getData {
				results[key] = data
			} else {
				results[key] = nil
			}

		}
	}

	return results, nil
}

func (s *StorageKubernetes) Update(ctx context.Context, dataset map[string][]byte) (err error) {
	secret, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Get(ctx, storageKubernetesSecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if secret == nil || secret.Data == nil || len(secret.Data) <= 0 {
		return fmt.Errorf("no credentials exist")
	}

	for key, data := range dataset {
		if _, ok := secret.Data[key]; !ok {
			return fmt.Errorf("credential %q does not exist", key)
		}

		secret.Data[key] = data
	}

	// k8s update secret
	if _, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (s *StorageKubernetes) Delete(ctx context.Context, keys []string) (err error) {
	secret, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Get(ctx, storageKubernetesSecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if secret == nil || secret.Data == nil || len(secret.Data) <= 0 {
		return fmt.Errorf("no credentials exist")
	}

	for _, key := range keys {
		if _, ok := secret.Data[key]; !ok {
			return fmt.Errorf("credential %q does not exist", key)
		}
		delete(secret.Data, key)
	}

	// k8s update secret
	if _, err := s.k8sClient.GetK8sClientset().CoreV1().Secrets(storageKubernetesNamespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func newSecretForCredStor() *corev1.Secret {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      storageKubernetesSecretName,
			Namespace: storageKubernetesNamespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{},
	}
	return secret
}
