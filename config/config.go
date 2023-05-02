package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Aws      AwsConfig
}

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     string
	Sslmode  string
	Timezone string
}

type RedisConfig struct {
	Addr     string
	Password string
	Db       int
}

type AwsConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func NewConfig(path string) (*Config, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, err
	}
	var c Config
	if err := k.Unmarshal("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
