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
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-file/parser"
	"github.com/kitex-contrib/config-file/utils"
)

type ConfigWatcher struct {
	config      *parser.ClientFileConfig
	callbacks   []func()
	key         string
	from        string // client name
	to          string // server name
	filewatcher *utils.FileWatcher
}

// NewConfigWatcher init a watcher for the config file
func NewConfigWatcher(filepath, from, to string) *ConfigWatcher {
	fw, err := utils.NewFileWatcher(filepath)
	if err != nil {
		panic(err)
	}

	return &ConfigWatcher{
		filewatcher: fw,
		key:         fmt.Sprintf("%s%s%s", from, "/", to),
		from:        from,
		to:          to,
	}
}

func (c *ConfigWatcher) ToService() string { return c.to }

func (c *ConfigWatcher) Key() string { return c.key }

func (c *ConfigWatcher) Config() *parser.ClientFileConfig { return c.config }

func (c *ConfigWatcher) Start() {
	c.parseHandler(c.filewatcher.FilePath())
	c.filewatcher.AddCallback(c.parseHandler)
	c.filewatcher.StartWatching()
}

func (c *ConfigWatcher) Stop() error {
	c.filewatcher.StopWatching()
	klog.Infof("[local] stop watching file: %s", c.filewatcher.FilePath())
	return nil
}

func (c *ConfigWatcher) AddCallback(callback func()) {
	c.callbacks = append(c.callbacks, callback)
}

// parseHandler parse and invoke each function in the callbacks array
func (c *ConfigWatcher) parseHandler(filepath string) {
	data, err := utils.ReadFileAll(filepath)
	if err != nil {
		klog.Errorf("[local] read config file failed: %v\n", err)
		return
	}

	resp := &parser.ClientFileManager{}
	err = parser.Decode(data, resp)
	if err != nil {
		klog.Errorf("[local] failed to parse the config file: %v\n", err)
		return
	}

	c.config = resp.GetConfig(c.key)
	if c.config == nil {
		klog.Warnf("[local] not matching key found, skip\n")
		return
	}

	if len(c.callbacks) > 0 {
		for _, callback := range c.callbacks {
			callback()
		}
	}
	klog.Infof("[local] server config parse and update complete \n")
}
