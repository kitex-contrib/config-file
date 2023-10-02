package parser

import "github.com/cloudwego/kitex/pkg/limiter"

type ServerFileConfig struct {
	Limit limiter.LimiterConfig `mapstructure:"limit"`
}

type ServerFileManager map[string]*ServerFileConfig

func (s *ServerFileManager) GetConfig(key string) *ServerFileConfig {
	config, exist := (*s)[key]

	if !exist {
		return nil
	}

	return config
}
