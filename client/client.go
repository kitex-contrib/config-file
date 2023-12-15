package client

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-file/monitor"
	"github.com/kitex-contrib/config-file/parser"
)

// getFileConfig returns the config from the watcher.
// if the config type is not *parser.ClientFileConfig, it will log an error and return nil.
func getFileConfig(watcher monitor.ConfigMonitor) *parser.ClientFileConfig {
	config, ok := watcher.Config().(*parser.ClientFileConfig)
	if !ok {
		// This should never happen.
		// But if it does, we should log it and do nothing.
		// Otherwise, the program will panic.
		klog.Errorf("[local] Invalid config type: %T", watcher.Config())
		return nil
	}
	return config
}
