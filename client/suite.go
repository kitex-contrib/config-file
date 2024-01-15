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
	kitexclient "github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/config-file/filewatcher"
	"github.com/kitex-contrib/config-file/monitor"
	"github.com/kitex-contrib/config-file/parser"
)

type FileConfigClientSuite struct {
	watcher monitor.ConfigMonitor
	service string
	fns     []parser.CustomFunction
}

// NewSuite service is the destination service.
func NewSuite(service, key string, watcher filewatcher.FileWatcher, cp parser.ConfigParser, cfs ...parser.CustomFunction) *FileConfigClientSuite {
	cm, err := monitor.NewConfigMonitor(key, watcher)
	if err != nil {
		panic(err)
	}

	if cp == nil {
		cm.SetParser(parser.DefaultConfigParse())
	} else {
		// use customized parser
		cm.SetParser(cp)
	}

	return &FileConfigClientSuite{
		watcher: cm,
		service: service,
		fns:     cfs,
	}
}

// Options return a list client.Option
func (s *FileConfigClientSuite) Options() []kitexclient.Option {
	s.watcher.SetManager(&parser.ClientFileManager{})

	opts := make([]kitexclient.Option, 0, 7)
	opts = append(opts, WithRetryPolicy(s.watcher)...)
	opts = append(opts, WithCircuitBreaker(s.service, s.watcher)...)
	opts = append(opts, WithRPCTimeout(s.watcher)...)
	opts = append(opts, kitexclient.WithCloseCallbacks(func() error {
		s.watcher.Stop()
		return nil
	}))

	if err := s.watcher.Start(); err != nil {
		panic(err)
	}

	return opts
}
