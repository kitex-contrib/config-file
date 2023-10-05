// Copyright 2023 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"sync/atomic"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	kitexserver "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-file/monitor"
	"github.com/kitex-contrib/config-file/parser"
)

// WithLimiter returns a server.Option that sets the limiter for the server.
func WithLimiter(watcher *monitor.ConfigMonitor) kitexserver.Option {
	return kitexserver.WithLimit(initLimitOptions(watcher))
}

// initLimitOptions init the limiter options
func initLimitOptions(watcher *monitor.ConfigMonitor) *limit.Option {
	var updater atomic.Value
	opt := &limit.Option{}

	opt.UpdateControl = func(u limit.Updater) {
		klog.Debugf("[local] server file limiter updater init, config %+v\n", *opt)
		u.UpdateLimit(opt)
		updater.Store(u)
	}

	onChangeCallback := func() {
		lc := watcher.Config().(*parser.ServerFileConfig).Limit

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
