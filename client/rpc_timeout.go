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
