package server

import (
	"sync/atomic"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/server"
)

// WithLimiter returns a server.Option that sets the limiter for the server.
func WithLimiter(watcher *ConfigWatcher) server.Option {
	return server.WithLimit(initLimitOptions(watcher))
}

// initLimitOptions init the limiter options
func initLimitOptions(watcher *ConfigWatcher) *limit.Option {
	var updater atomic.Value
	opt := &limit.Option{}

	opt.UpdateControl = func(u limit.Updater) {
		klog.Debugf("[local] server file limiter updater init, config %+v\n", *opt)
		u.UpdateLimit(opt)
		updater.Store(u)
	}

	onChangeCallback := func() {
		lc := watcher.Config().Limit
		klog.Infof("current maxConnections: %v, new: %v\n", opt.MaxConnections, lc.ConnectionLimit)

		opt.MaxConnections = int(lc.ConnectionLimit)
		opt.MaxQPS = int(lc.QPSLimit)

		u := updater.Load()
		if u == nil {
			klog.Warnf("[local] %s server limiter config: failed as the updater is empty", watcher.Key())
			return
		}
		if !u.(limit.Updater).UpdateLimit(opt) {
			klog.Warnf("[local] %s server limiter config: update may do not take affect", watcher.Key())
		}
	}

	watcher.AddCallback(onChangeCallback)

	return opt
}
