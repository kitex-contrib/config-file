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
)

type FileConfigClientSuite struct {
	from     string // client service name
	to       string // server service name
	filePath string // config filepath
}

// NewSuite service is the destination service.
func NewSuite(from, to, filePath string) *FileConfigClientSuite {
	return &FileConfigClientSuite{
		from:     from,
		to:       to,
		filePath: filePath,
	}
}

// Options return a list client.Option
func (s *FileConfigClientSuite) Options() []client.Option {
	watcher := NewConfigWatcher(s.filePath, s.from, s.to)

	opts := make([]client.Option, 0, 5)
	opts = append(opts, WithRetryPolicy(watcher))
	opts = append(opts, WithCircuitBreaker(watcher)...)
	opts = append(opts, WithRPCTimeout(watcher))
	opts = append(opts, client.WithCloseCallbacks(watcher.Stop))

	watcher.Start()

	return opts
}
