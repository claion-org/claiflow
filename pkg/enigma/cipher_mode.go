//go:generate go run github.com/abice/go-enum --file=cipher_mode.go --names --nocase
package enigma

import (
	"crypto/cipher"
	"fmt"
	"reflect"
	"runtime"
)

/*
ENUM(
NONE
CBC
GCM
)
*/
type CipherMode int

func (mode CipherMode) CipherFactory(block cipher.Block, salt *Salt) (Encoder, Decoder, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		switch r := r.(type) {
	// 		case error:
	// 			err = errors.Wrapf(r, "recovered")
	// 		default:
	// 			err = fmt.Errorf("recovered %v", r)
	// 		}
	// 	}
	// }()

	switch mode {
	case CipherModeNONE:

		switch block.BlockSize() {
		case 1:
			salt.SetLen(0) //set salt len

			encoder := func(src, salt []byte) (dst []byte, err error) {
				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Encrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				return
			}
			decoder := func(src, salt []byte) (dst []byte, err error) {
				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Decrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				return
			}

			return encoder, decoder, nil
		default:
			salt.SetLen(block.BlockSize()) //set salt len

			encoder := func(src, salt []byte) (dst []byte, err error) {
				src = PKCS7Padding(src, block.BlockSize())

				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Encrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				return
			}
			decoder := func(src, salt []byte) (dst []byte, err error) {
				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Decrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				dst = PKCS7Unpadding(dst)

				return
			}

			return encoder, decoder, nil
		}

	case CipherModeCBC:
		salt.SetLen(block.BlockSize()) //set salt len

		encoder := func(src, iv []byte) (dst []byte, err error) {
			dst = make([]byte, len(src))
			cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, src)
			return
		}
		decoder := func(src, iv []byte) (dst []byte, err error) {
			dst = make([]byte, len(src))
			cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, src)
			return
		}

		return encoder, decoder, nil

	case CipherModeGCM:
		c, err := cipher.NewGCM(block)
		if err != nil {
			err := fmt.Errorf("%w: new gcm cipher type=%q", err, fmt.Sprintf("%T", block))
			return nil, nil, err
		}

		salt.SetLen(c.NonceSize()) //set salt len

		encoder := func(src, nonce []byte) (dst []byte, err error) {
			dst = c.Seal(nonce, nonce, src, nil)
			dst = dst[len(nonce):] //remove nonce
			return
		}
		decoder := func(src, nonce []byte) (dst []byte, err error) {
			dst, err = c.Open(nil, nonce, src, nil)
			return
		}

		return encoder, decoder, nil

	default:
		err := fmt.Errorf("invalid cipher mode=%q", mode.String())
		return nil, nil, err
	}

}

func typeName(i interface{}) string {
	t := reflect.ValueOf(i).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	}
	return t.String()
}
