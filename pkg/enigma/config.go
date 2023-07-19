package enigma

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// blockMethod: aes    # NONE, AES, DES
// blockSize: 128      # NONE: default(1), AES: 128|192|256, DES: 64
// blockKey: secret    # (base64 string)
// cipherMode: gcm     # NONE: NONE|AES|DES , GCM: AES, CBC: NONE|AES|DES
// cipherSalt: null    # NULL, (base64 string)
// padding: PKCS       # NONE: AES+GCM, PKCS: AES+NONE|AES+CBC|DES+NONE|DES+CBC
// strconv: base64     # plain|base64|hex
type Config struct {
	ConfigBlock   `yaml:",inline"`
	ConfigCipher  `yaml:",inline"`
	ConfigPadding `yaml:",inline"`
	ConfigStrConv `yaml:",inline"`
}

func (cfg Config) ToOption() MachineOption {
	return configToOption(cfg)
}

func configToOption(cfg Config) (opt MachineOption) {
	opt.Block.Method = cfg.ConfigBlock.EncryptionMethod
	opt.Block.Size = cfg.ConfigBlock.BlockSize
	opt.Block.Key = cfg.ConfigBlock.BlockKey
	opt.Cipher.Mode = cfg.ConfigCipher.CipherMode
	opt.Cipher.Salt = cfg.ConfigCipher.CipherSalt
	opt.Padding = cfg.ConfigPadding.Padding
	opt.StrConv = cfg.ConfigStrConv.StrConv

	return
}

type ConfigBlock struct {
	EncryptionMethod string `yaml:"blockMethod"` // NONE|AES|DES
	BlockSize        int    `yaml:"blockSize"`   // NONE: default(1), AES: [128|192|256], DES: [64]
	BlockKey         string `yaml:"blockKey"`    // (base64 string)
}

type ConfigCipher struct {
	CipherMode string  `yaml:"cipherMode"` // NONE|CBC|GCM
	CipherSalt *string `yaml:"cipherSalt"` // nil: auto-generate (base64 string)
}

type ConfigPadding struct {
	Padding string `yaml:"padding"` // none|PKCS
}

type ConfigStrConv struct {
	StrConv string `yaml:"strconv"` // plain|base64|hex
}

func PrintConfig(w io.Writer, cfgset map[string]Config, insecure bool) {
	nullstring := func(p *string) (s string) {
		s = fmt.Sprintf("%v", p)
		if p != nil {
			s = *p
		}
		return
	}

	fmt.Fprintln(w, "enigma configuration:")

	tabwrite := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	col := []string{}
	col = append(col, "")
	col = append(col, "name")
	col = append(col, "encryption-method")
	col = append(col, "block-size")
	if insecure {
		col = append(col, "block-key")
	}
	col = append(col, "cipher-mode")
	if insecure {
		col = append(col, "cipher-salt")
	}
	col = append(col, "padding")
	col = append(col, "strconv")

	tabwrite.Write([]byte(strings.Join(col, "\t") + "\n"))

	for name, cfg := range cfgset {

		row := []string{}
		row = append(row, "-")
		row = append(row, name)
		row = append(row, cfg.EncryptionMethod)
		row = append(row, fmt.Sprintf("%v", cfg.BlockSize))
		if insecure {
			row = append(row, cfg.BlockKey)
		}
		row = append(row, cfg.CipherMode)
		if insecure {
			row = append(row, nullstring(cfg.CipherSalt))
		}
		row = append(row, cfg.Padding)
		row = append(row, cfg.StrConv)

		tabwrite.Write([]byte(strings.Join(row, "\t") + "\n"))
	}

	tabwrite.Flush()

	fmt.Fprintln(w, strings.Repeat("_", 40))
}
