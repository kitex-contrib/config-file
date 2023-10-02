package client

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/rpctimeout"
)

// WithRPCTimeout returns a server.Option that sets the timeout provider for the client.
func WithRPCTimeout(watcher *ConfigWatcher) client.Option {
	return client.WithTimeoutProvider(initRPCTimeout(watcher))
}

// initRPCTimeout init the rpc timeout provider
func initRPCTimeout(watcher *ConfigWatcher) rpcinfo.TimeoutProvider {
	rpcTimeoutContainer := rpctimeout.NewContainer()

	onChangeCallback := func() {
		// the key is method name, wildcard "*" can match anything.
		configs := watcher.Config().Timeout
		rpcTimeoutContainer.NotifyPolicyChange(configs)
	}

	watcher.AddCallback(onChangeCallback)
	return rpcTimeoutContainer
}
