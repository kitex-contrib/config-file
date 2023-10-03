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

package monitor

import (
	"errors"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-file/parser"
	"github.com/kitex-contrib/config-file/utils"
)

type Options struct {
	FilePath string      // config file path
	Key      string      // config key
	Provider KeyProvider // custom key provider
}

type ConfigMonitor struct {
	manager     parser.ConfigManager // Manager for the config file
	config      interface{}          // config details
	fileWatcher *utils.FileWatcher   // local config file watcher
	callbacks   []func()             // callbacks when config file changed
	key         string               // key
	provider    KeyProvider          // provider
}

// NewConfigMonitor init a monitor for the config file
func NewConfigMonitor(opts Options) (*ConfigMonitor, error) {
	if opts.FilePath == "" {
		return nil, errors.New("empty config file path")
	}
	if opts.Key == "" {
		return nil, errors.New("empty config key")
	}
	if opts.Provider == nil {
		opts.Provider = NewDefaultKeyProvider()
	}

	fw, err := utils.NewFileWatcher(opts.FilePath)
	if err != nil {
		return nil, err
	}

	return &ConfigMonitor{
		fileWatcher: fw,
		key:         opts.Key,
		provider:    opts.Provider,
	}, nil
}

// Key return the key of the config file
func (c *ConfigMonitor) Key() string { return c.key }

// Config return the config details
func (c *ConfigMonitor) Config() interface{} { return c.config }

// Start starts the file watch progress
func (c *ConfigMonitor) Start() error {
	if c.manager == nil {
		return errors.New("not set manager for config file")
	}
	c.parseHandler(c.fileWatcher.FilePath())
	c.fileWatcher.AddCallback(c.parseHandler)
	return c.fileWatcher.StartWatching()
}

// Stop stops the file watch progress
func (c *ConfigMonitor) Stop() {
	c.fileWatcher.StopWatching()
	klog.Infof("[local] stop watching file: %s", c.fileWatcher.FilePath())
}

// SetManager set the manager for the config file
func (c *ConfigMonitor) SetManager(manager parser.ConfigManager) { c.manager = manager }

// AddCallback add callback function, it will be called when file changed
func (c *ConfigMonitor) AddCallback(callback func()) {
	c.callbacks = append(c.callbacks, callback)
}

// parseHandler parse and invoke each function in the callbacks array
func (c *ConfigMonitor) parseHandler(filepath string) {
	data, err := utils.ReadFileAll(filepath)
	if err != nil {
		klog.Errorf("[local] read config file failed: %v\n", err)
		return
	}

	resp := c.manager
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
