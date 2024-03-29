// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mock

import "github.com/kitex-contrib/config-file/filewatcher"

type fwmock struct{}

// NewMockFileWatcher will return a mock filewatcher
func NewMockFileWatcher() filewatcher.FileWatcher { return &fwmock{} }

func (fw *fwmock) FilePath() string { return "test" }

func (fw *fwmock) CallbackSize() int { return 1 }

func (fw *fwmock) RegisterCallback(callback func(data []byte)) int64 { return 0 }

func (fw *fwmock) DeregisterCallback(uniqueID int64) {}

func (fw *fwmock) StartWatching() error { return nil }

func (fw *fwmock) StopWatching() {}

func (fw *fwmock) CallOnceAll() error { return nil }

func (fw *fwmock) CallOnceSpecific(uniqueID int64) error { return nil }
