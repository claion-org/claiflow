package cryptography

import (
	"github.com/claion-org/claiflow/pkg/enigma"
	"github.com/pkg/errors"
)

// Config
//
//	 enigma:
//		blockMethod: aes    # NONE, AES, DES
//		blockSize: 128      # NONE: default(1), AES: 128|192|256, DES: 64
//		blockKey: secret    # (base64 string)
//		cipherMode: gcm     # NONE: NONE|AES|DES , GCM: AES, CBC: NONE|AES|DES
//		cipherSalt: null    # NULL, (base64 string)
//		padding: PKCS       # NONE: AES+GCM, PKCS: AES+NONE|AES+CBC|DES+NONE|DES+CBC
//		strconv: base64     # plain|base64|hex
type Config struct {
	EnigmaConfig enigma.Config `yaml:"enigma"`
}

func LoadConfig(config Config) error {
	machine, err := enigma.NewMachine(config.EnigmaConfig.ToOption())
	if err != nil {
		return errors.Wrapf(err, "new machine")
	}
	Cipher = machine

	return nil
}
