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
	"fmt"
	"strings"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

type KeyProvider interface {
	GetKey(rpcinfo.RPCInfo) string
	ParseKey(key string) (*Key, error)
}

type Key struct {
	ClientName string
	ServerName string
	Tags       map[string]string
}

type DefaultKeyProvider struct{}

func NewDefaultKeyProvider() *DefaultKeyProvider {
	return &DefaultKeyProvider{}
}

func (d *DefaultKeyProvider) GetKey(info rpcinfo.RPCInfo) string {
	return fmt.Sprintf("%s%s%s", info.From().ServiceName(), "/", info.To().ServiceName())
}

func (d *DefaultKeyProvider) ParseKey(key string) (*Key, error) {
	parts := strings.Split(key, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid key: %s", key)
	}

	return &Key{
		ClientName: parts[0],
		ServerName: parts[1],
	}, nil
}
