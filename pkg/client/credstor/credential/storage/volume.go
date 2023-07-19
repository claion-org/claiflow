package storage

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

var _ Storage = &StorageVolume{}

var defaultVolumeDirPath = ".claiflow_credentials"

type StorageVolume struct {
	dirPath      string
	encrypt      bool
	getSecretKey func() ([]byte, error)
}

func NewStorageVolume(path string, encrypt bool, getSecretKey func() ([]byte, error)) (Storage, error) {
	cs := &StorageVolume{encrypt: encrypt, getSecretKey: getSecretKey}
	if len(path) > 0 {
		cs.dirPath = path
	} else {
		cs.dirPath = defaultVolumeDirPath
	}

	// check dir
	if err := checkAndCreateDir(cs.dirPath); err != nil {
		return nil, err
	}

	return cs, nil
}

func (s *StorageVolume) Add(ctx context.Context, dataset map[string][]byte) (err error) {
	// check dir
	if err := checkAndCreateDir(s.dirPath); err != nil {
		return err
	}
	var writeDataKeys []string

	defer func() {
		if err != nil {
			for _, key := range writeDataKeys {
				os.Remove(key)
			}
		}
	}()

	for key, data := range dataset {
		var wBytes []byte
		if s.encrypt {
			secretKey, err := s.getSecretKey()
			if err != nil {
				return err
			}

			encryptedData, err := encryptAES256GCM(secretKey, data)
			if err != nil {
				return err
			}

			wBytes = encryptedData
		} else {
			wBytes = data
		}

		f, err := os.OpenFile(filepath.Join(s.dirPath, key), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return err
		}

		_, err = f.Write(wBytes)
		if err1 := f.Close(); err1 != nil && err == nil {
			err = err1
		}

		if err != nil {
			return err
		}
		writeDataKeys = append(writeDataKeys, key)
	}

	return nil
}

func (s *StorageVolume) Get(ctx context.Context, key string, getData bool) ([]byte, error) {
	// check dir
	if err := checkAndCreateDir(s.dirPath); err != nil {
		return nil, err
	}

	if !getData {
		if _, err := os.Stat(filepath.Join(s.dirPath, key)); err != nil {
			return nil, err
		}

		return nil, nil
	}

	b, err := os.ReadFile(filepath.Join(s.dirPath, key))
	if err != nil {
		return nil, err
	}

	var rBytes []byte
	if s.encrypt {
		secretKey, err := s.getSecretKey()
		if err != nil {
			return nil, err
		}

		// decrypt file
		decryptedData, err := decryptAES256GCM(secretKey, b)
		if err != nil {
			return nil, err
		}
		rBytes = decryptedData
	} else {
		rBytes = b
	}

	return rBytes, nil
}

func (s *StorageVolume) GetAll(ctx context.Context, getData bool) (map[string][]byte, error) {
	// check dir
	if err := checkAndCreateDir(s.dirPath); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.dirPath)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for _, e := range entries {
		if !e.IsDir() {
			key := e.Name()
			b, err := s.Get(ctx, key, getData)
			if err != nil {
				return nil, err
			}

			result[key] = b
		}
	}

	return result, nil
}

func (s *StorageVolume) Update(ctx context.Context, dataset map[string][]byte) (err error) {
	// check dir
	if err := checkAndCreateDir(s.dirPath); err != nil {
		return err
	}
	oldDataset := make(map[string][]byte)

	defer func() {
		if err != nil {
			for key, data := range oldDataset {
				os.WriteFile(filepath.Join(s.dirPath, key), data, 0644)
			}
		}
	}()

	for key, data := range dataset {
		var wBytes []byte
		if s.encrypt {
			secretKey, err := s.getSecretKey()
			if err != nil {
				return err
			}

			encryptedData, err := encryptAES256GCM(secretKey, data)
			if err != nil {
				return err
			}

			wBytes = encryptedData
		} else {
			wBytes = data
		}

		oldData, err := os.ReadFile(filepath.Join(s.dirPath, key))
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(s.dirPath, key), wBytes, 0644); err != nil {
			return err
		}

		oldDataset[key] = oldData
	}

	return nil
}

func (s *StorageVolume) Delete(ctx context.Context, keys []string) (err error) {
	// check dir
	if err := checkAndCreateDir(s.dirPath); err != nil {
		return err
	}
	oldDataset := make(map[string][]byte)

	defer func() {
		if err != nil {
			for key, data := range oldDataset {
				os.WriteFile(filepath.Join(s.dirPath, key), data, 0644)
			}
		}
	}()

	for _, key := range keys {
		oldData, err := os.ReadFile(filepath.Join(s.dirPath, key))
		if err != nil {
			return err
		}

		if err := os.Remove(filepath.Join(s.dirPath, key)); err != nil {
			return err
		}

		oldDataset[key] = oldData
	}

	return nil
}

func checkAndCreateDir(path string) error {
	// check base dir
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			// if not exist, create directory
			if err := os.Mkdir(path, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// var getSecretKey = func() ([]byte, error) {
// 	return fetcher.ClusterUuid
// }

// func getSecretKey() ([]byte, error)

func encryptAES256GCM(key, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key length must be 32 for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func decryptAES256GCM(key []byte, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key length must be 32 for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, ciphertextWithoutNonce := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertextWithoutNonce, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
