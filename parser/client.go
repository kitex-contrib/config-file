package parser

import (
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpctimeout"
)

type ClientFileConfig struct {
	Timeout        map[string]*rpctimeout.RPCTimeout `mapstructure:"timeout"`
	Retry          map[string]*retry.Policy          `mapstructure:"retry"`
	Circuitbreaker map[string]*circuitbreak.CBConfig `mapstructure:"circuitbreaker"`
}

type ClientFileManager map[string]*ClientFileConfig

func (s *ClientFileManager) GetConfig(key string) *ClientFileConfig {
	config, exist := (*s)[key]

	if !exist {
		return nil
	}

	return config
}
