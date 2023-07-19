package enigma_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/claion-org/claiflow/pkg/enigma"
)

func TestEnigma_10(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = ""
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 CBC PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 CBC PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_101(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 CBC PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 CBC PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_111(t *testing.T) {
	//AES 128 GCM PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "GCM"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_121(t *testing.T) {
	//AES 128 NONE PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_131(t *testing.T) {
	//AES 128 GCM NONE SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "NONE"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM NONE SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM NONE SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_141(t *testing.T) {
	//AES 128 GCM PKCS NULL
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	// crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM NONE NULL
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM NONE NULL
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_151(t *testing.T) {
	//AES 128 GCM NONE NONE
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "NONE"
	// crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM NONE SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM NONE SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_102(t *testing.T) {
	//DES 64 CBC PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_112(t *testing.T) {
	//DES 64 CBC PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_103(t *testing.T) {
	//DES 64 NONE PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

}

func TestEnigma_1(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "cbc"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_12(t *testing.T) {
	//NONE 128 NONE PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "NONE"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//NONE 128 CBC PKCS NONE
	crypto_alg.CipherMode = "CBC"
	EnigmaMachine(t, crypto_alg)

	//NONE 128 GCM PKCS NONE
	crypto_alg.CipherMode = "GCM"
	if false {
		EnigmaMachine(t, crypto_alg)
	}
}

func TestEnigma_13(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_14(t *testing.T) {
	//AES 128 NONE PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	crypto_alg.CipherMode = "CBC"
	EnigmaMachine(t, crypto_alg)

	crypto_alg.CipherMode = "GCM"
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_15(t *testing.T) {
	//DES 128 NONE PKCS SALTY
	var crypto_alg enigma.Config
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	crypto_alg.CipherMode = "CBC"
	EnigmaMachine(t, crypto_alg)
}

func EnigmaMachine(t *testing.T, alg enigma.Config) {

	crypto, err := enigma.NewMachine(alg.ToOption())
	if err != nil {
		t.Fatal(err)
	}

	s := "세종어제 훈민정음\n" +
		"나랏말이\n" +
		"중국과 달라\n" +
		"문자와 서로 통하지 아니하므로\n" +
		"이런 까닭으로 어리석은 백성이 이르고자 하는 바가 있어도\n" +
		"마침내 제 뜻을 능히 펴지 못하는 사람이 많다.\n" +
		"내가 이를 위해 불쌍히 여겨\n" +
		"새로 스물여덟 글자를 만드니\n" +
		"사람마다 하여금 쉬이 익혀 날마다 씀에 편안케 하고자 할 따름이다.\n"

	var encripttext, plaintext []byte

	var salt_a, salt_b []byte

	encripttext, err = crypto.Encode([]byte(s))
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err = crypto.Decode(encripttext)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(salt_a, salt_b) {
		t.Fatal("diff salt")
	}

	if s != string(plaintext) {
		t.Fatal("diff text", string(plaintext))
	}
}

func TestMain(t *testing.T) {
	NewString := func(s string) *string { return &s }

	config := enigma.Config{}
	config.EncryptionMethod = "AES"
	config.BlockSize = 128
	config.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	config.CipherMode = "GCM"
	config.Padding = "PKCS"
	config.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	config.StrConv = "base64"

	const example = "세종어제 훈민정음\n" +
		"나랏말이\n" +
		"중국과 달라\n" +
		"문자와 서로 통하지 아니하므로\n" +
		"이런 까닭으로 어리석은 백성이 이르고자 하는 바가 있어도\n" +
		"마침내 제 뜻을 능히 펴지 못하는 사람이 많다.\n" +
		"내가 이를 위해 불쌍히 여겨\n" +
		"새로 스물여덟 글자를 만드니\n" +
		"사람마다 하여금 쉬이 익혀 날마다 씀에 편안케 하고자 할 따름이다.\n"

	machine, err := enigma.NewMachine(config.ToOption())
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := machine.Encode([]byte(example))
	if err != nil {
		t.Fatal(err)
	}

	plain, err := machine.Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.EqualFold(example, string(plain)) {
		t.Fatal("diff")
	}
}
