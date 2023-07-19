package cryptography

import (
	"github.com/claion-org/claiflow/pkg/enigma"
	"github.com/pkg/errors"
)

var (
	Cipher enigma.Cipher
)

func CipherSet() (enigma.Cipher, bool) {
	return Cipher, Cipher != nil
}

func EnigmaEncode(bytes []byte) ([]byte, error) {
	cipher, ok := CipherSet()
	if !ok {
		return nil, errors.Errorf("no cipher")
	}

	out, err := cipher.Encode(bytes)
	return out, errors.Wrapf(err, "enigma encode")
}

func EnigmaDecode(bytes []byte) ([]byte, error) {
	cipher, ok := CipherSet()
	if !ok {
		return nil, errors.Errorf("no cipher")
	}

	out, err := cipher.Decode(bytes)
	return out, errors.Wrapf(err, "enigma decode")
}
