package config

import (
	"github.com/claion-org/claiflow/pkg/echov4/middleware"
	"github.com/claion-org/claiflow/pkg/enigma"
	"github.com/claion-org/claiflow/pkg/logger"
	"github.com/jinzhu/configor"
)

type Config struct {
	APPName     string            `yaml:"appName" default:"flow-api"`
	HttpService Http              `yaml:"http"`
	GrpcService Grpc              `yaml:"grpc"`
	Database    Database          `yaml:"database"`
	Migrate     Migrate           `yaml:"migrate"`
	Enigma      enigma.Config     `yaml:"enigma"`
	Logger      logger.ZaprConfig `yaml:"logger"`
}

func New(c *Config, configPath string) (*Config, error) {
	if c == nil {
		c = &Config{}
	}
	if err := configor.Load(c, configPath); err != nil {
		return nil, err
	}
	return c, nil
}

type Http struct {
	Port      int32                 `yaml:"port"      default:"8099"`
	URLPrefix string                `yaml:"urlPrefix" default:""`
	Tls       Tls                   `yaml:"tls"`
	CORS      middleware.CORSConfig `yaml:"cors"`
}

type Grpc struct {
	Port           int32 `yaml:"port"           default:"18099"`
	Tls            Tls   `yaml:"tls"`
	MaxRecvMsgSize int   `yaml:"maxRecvMsgSize" default:"1073741824"` // 1GiB
}

type Tls struct {
	Enable   bool   `yaml:"enable"   default:"false"`
	CertFile string `yaml:"certFile" default:"server.crt"`
	KeyFile  string `yaml:"keyFile"  default:"server.key"`
}

type Database struct {
	Type            string `yaml:"type"            default:"mysql"`
	Protocol        string `yaml:"protocol"        default:"tcp"`
	Host            string `yaml:"host"            default:"localhost"`
	Port            string `yaml:"port"            default:"3306"`
	DBName          string `yaml:"dbname"          default:"flow"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	MaxOpenConns    int    `yaml:"maxOpenConns"    default:"15"`
	MaxIdleConns    int    `yaml:"maxIdleConns"    default:"5"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime" default:"1"`
	// ShowSQL         bool   `yaml:"showSql"         default:"false"`
	// LogLevel        string `yaml:"logLevel"        default:"warn"`
}

type Migrate struct {
	Source string `yaml:"source" default:"migrations/mysql"`
}
