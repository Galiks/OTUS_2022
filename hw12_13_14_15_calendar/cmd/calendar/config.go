package main

import (
	"github.com/pelletier/go-toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf
	PostgreSQL PostgresConf
	Server     ServerConf
}

type ServerConf struct {
	Host string
	Port string
}

type PostgresConf struct {
	ConnectString string
}

type LoggerConf struct {
	Level           string
	PrintStackTrace bool
	PathToFile      string
}

func NewConfig(path string) (*Config, error) {
	var conf = new(Config)

	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	data, err := tree.Marshal()
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
