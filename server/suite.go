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
	kitexserver "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-file/filewatcher"
	"github.com/kitex-contrib/config-file/monitor"
	"github.com/kitex-contrib/config-file/parser"

)

type FileConfigServerSuite struct {
	watcher monitor.ConfigMonitor
}

// NewSuite service is the destination service.
func NewSuite(key string, watcher filewatcher.FileWatcher) *FileConfigServerSuite {
	cm, err := monitor.NewConfigMonitor(key, watcher)
	if err != nil {
		panic(err)
	}

	return &FileConfigServerSuite{
		watcher: cm,
	}
}

// Options return a list client.Option
func (s *FileConfigServerSuite) Options() []kitexserver.Option {
	s.watcher.SetManager(&parser.ServerFileManager{})
	s.watcher.SetParser(&parser.Parser{})

	opts := make([]kitexserver.Option, 0, 1)
	opts = append(opts, WithLimiter(s.watcher))

	if err := s.watcher.Start(); err != nil {
		panic(err)
	}

	kitexserver.RegisterShutdownHook(s.watcher.Stop)

	return opts
}
