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
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-file/monitor"
	"github.com/kitex-contrib/config-file/parser"
)

type FileConfigServerSuite struct {
	watcher *monitor.ConfigMonitor
}

// NewSuite service is the destination service.
func NewSuite(watcher *monitor.ConfigMonitor) *FileConfigServerSuite {
	return &FileConfigServerSuite{
		watcher: watcher,
	}
}

// Options return a list client.Option
func (s *FileConfigServerSuite) Options() []server.Option {
	s.watcher.SetManager(&parser.ServerFileManager{})

	opts := make([]server.Option, 0, 1)
	opts = append(opts, WithLimiter(s.watcher))

	if err := s.watcher.Start(); err != nil {
		panic(err)
	}

	server.RegisterShutdownHook(s.watcher.Stop)

	return opts
}
