package config

import (
	"awesomeProjectRucenter/pkg/erx"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"sync"
)

type Logger struct {
	Path string
}

type Server struct {
	Domain string
	Port   int
}

type VmDb struct {
	PathDb string
}

type Config struct {
	Logger Logger
	Server Server
	VmDb   VmDb
}

var instance *Config
var once sync.Once

func GetConfigInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func GetConfig(config *Config) error {
	file, err := ioutil.ReadFile("./app.conf")
	if err != nil {
		return erx.New(err)
	}
	err = toml.Unmarshal(file, config)
	if err != nil {
		return erx.New(err)
	}
	return nil
}
