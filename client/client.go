// Copyright 2024 CloudWeGo Authors
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
