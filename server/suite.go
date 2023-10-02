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
)

type FileConfigServerSuite struct {
	service  string // service name
	filePath string // config filepath
}

// NewSuite service is the destination service.
func NewSuite(service, filePath string) *FileConfigServerSuite {
	return &FileConfigServerSuite{
		service:  service,
		filePath: filePath,
	}
}

// Options return a list client.Option
func (s *FileConfigServerSuite) Options() []server.Option {
	watcher := NewConfigWatcher(s.filePath, s.service)

	opts := make([]server.Option, 0, 1)
	opts = append(opts, WithLimiter(watcher))

	watcher.Start()
	server.RegisterShutdownHook(watcher.Stop)

	return opts
}
