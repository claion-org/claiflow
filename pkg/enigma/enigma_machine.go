package enigma

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
)

type Cipher interface {
	EncodeDetail(src []byte, callback ...func(map[string]interface{})) ([]byte, error)
	Encode(src []byte) ([]byte, error)
	DecodeDetail(src []byte, callback ...func(map[string]interface{})) ([]byte, error)
	Decode(src []byte) ([]byte, error)
}

type Encoder func(src, salt []byte) (dst []byte, err error)
type Decoder func(src, salt []byte) (dst []byte, err error)

type Machine struct {
	method  func() EncryptionMethod
	mode    func() CipherMode
	key     func() []byte
	padding func() Padding
	strconv func() StrConv
	salt    func() *Salt
	block   func() cipher.Block
	Encoder
	Decoder
}

func (machine *Machine) EncodeDetail(src []byte, callback ...func(map[string]interface{})) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered %s %s encoder",
					machine.method().String(),
					machine.mode().String(),
				)
			default:
				err = fmt.Errorf("recovered %s %s encoder: %v",
					machine.method().String(),
					machine.mode().String(),
					r,
				)
			}
		}
	}()

	//salt
	// salt, hasSalt := machine.salt().GenSalt(), machine.salt().Has()
	err = machine.salt().Scope(func(ss *ScopeSalt) error {
		//padding
		src = machine.padding().Padder()(src, machine.block().BlockSize())
		//encode
		dst, err = machine.Encoder(src, ss.GenSalt())
		if err != nil {
			err := fmt.Errorf("%w: enigma encode src=%x salt=%x", err,
				src,
				base64.StdEncoding.EncodeToString(ss.GenSalt()))
			return err
		}

		for _, callback := range callback {
			callback(map[string]interface{}{
				"encript":     dst,
				"method":      machine.method().String(),
				"block_size":  machine.block().BlockSize(),
				"block_key":   machine.key(),
				"cipher_mode": machine.mode().String(),
				"cipher_salt": ss.GenSalt(),
				"padding":     machine.padding().String(),
			})
		}

		//salt encode rule
		dst = SaltEncodeRule(dst, ss.GenSalt(), ss.Has())
		//string converter encode
		dst = machine.strconv().Encoder()(dst)

		return nil
	})

	return
}

func (machine *Machine) Encode(src []byte) ([]byte, error) {
	return machine.EncodeDetail(src)
}

func (machine *Machine) DecodeDetail(src []byte, callback ...func(map[string]interface{})) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered %s %s decoder",
					machine.method().String(),
					machine.mode().String(),
				)
			default:
				err = fmt.Errorf("recovered %s %s decoder: %v",
					machine.method().String(),
					machine.mode().String(),
					r,
				)
			}
		}
	}()

	//salt
	err = machine.salt().Scope(func(ss *ScopeSalt) error {
		//string converter decode
		src, err = machine.strconv().Decoder()(src)
		//salt decode rule
		src, salt_ := SaltDecodeRule(src, ss.GenSalt(), ss.Has())
		//decode
		dst, err = machine.Decoder(src, salt_)
		if err != nil {
			err := fmt.Errorf("%w: enigma decode src=%x salt=%x", err,
				src,
				base64.StdEncoding.EncodeToString(ss.GenSalt()))
			return err
		}

		//unpadding
		dst = machine.padding().Unpadder()(dst)

		for _, callback := range callback {
			callback(map[string]interface{}{
				"encript":     dst,
				"method":      machine.method().String(),
				"block_size":  machine.block().BlockSize(),
				"block_key":   machine.key(),
				"cipher_mode": machine.mode().String(),
				"cipher_salt": salt_,
				"padding":     machine.padding().String(),
			})
		}

		return nil
	})

	return
}

func (machine *Machine) Decode(src []byte) ([]byte, error) {
	return machine.DecodeDetail(src)
}

type MachineOption struct {
	Block struct {
		Method string `json:"block-method"`
		Size   int    `json:"block-size"`
		Key    string `json:"block-key"`
	} `json:",inline"`
	Cipher struct {
		Mode string  `json:"cipher-mode"`
		Salt *string `json:"cipher-salt,omitempty"`
	} `json:",inline"`
	Padding string `json:"padding"`
	StrConv string `json:"strconv"`
}

func NewMachine(opt MachineOption) (m *Machine, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "enigma new machine")
			default:
				err = fmt.Errorf("enigma new machine: %v", r)
			}
		}
	}()

	method, err := ParseEncryptionMethod(opt.Block.Method)
	if err != nil {
		err := errors.Wrapf(err, "parse encryption method value=%q", opt.Block.Method)
		return nil, err
	}

	cipherMode, err := ParseCipherMode(opt.Cipher.Mode)
	if err != nil {
		err := errors.Wrapf(err, "parse cipher mode value=%q", opt.Cipher.Mode)
		return nil, err
	}

	padding, err := ParsePadding(opt.Padding)
	if err != nil {
		err := errors.Wrapf(err, "parse padding value=%q", opt.Padding)
		return nil, err
	}

	strconv, err := ParseStrConv(opt.StrConv)
	if err != nil {
		err := errors.Wrapf(err, "parse strconv value=%q", opt.StrConv)
		return nil, err
	}

	buf, err := base64.StdEncoding.DecodeString(opt.Block.Key)
	if err != nil {
		err := errors.Wrapf(err, "decode key value=%q", opt.Block.Key)
		return nil, err
	}
	blockKey := make([]byte, opt.Block.Size/8)
	copy(blockKey, buf)

	var salt Salt
	if opt.Cipher.Salt != nil {
		b, err := base64.StdEncoding.DecodeString(*opt.Cipher.Salt)
		if err != nil {
			err := errors.Wrapf(err, "decode salt value=%q", *opt.Cipher.Salt)
			return nil, err
		}
		salt.SetValue(b)
	}

	block, err := method.BlockFactory()(blockKey)
	if err != nil {
		err := errors.Wrapf(err, "block factory method=%q blockKey=%q", method, opt.Block.Key)
		return nil, err
	}

	encoder, decoder, err := cipherMode.CipherFactory(block, &salt)
	if err != nil {
		err := errors.Wrapf(err, "block factory method=%q blockKey=%q mode=%q",
			method,
			opt.Block.Key,
			opt.Cipher.Mode)
		return nil, err
	}

	m = &Machine{
		method:  func() EncryptionMethod { return method },
		mode:    func() CipherMode { return cipherMode },
		key:     func() []byte { return blockKey },
		padding: func() Padding { return padding },
		strconv: func() StrConv { return strconv },
		salt:    func() *Salt { return &salt },
		block:   func() cipher.Block { return block },
		Encoder: encoder,
		Decoder: decoder,
	}

	return
}
