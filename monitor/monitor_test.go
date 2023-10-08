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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kitex-contrib/config-file/filewatcher"
	"github.com/kitex-contrib/config-file/mock"
	"github.com/kitex-contrib/config-file/parser"
)

func TestNewConfigMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	if _, err := NewConfigMonitor("test", m); err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
}

func TestNewConfigMonitorFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	if _, err := NewConfigMonitor("", m); err == nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	if _, err := NewConfigMonitor("test", nil); err == nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
}

func TestKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	if cm.Key() != "test" {
		t.Errorf("Key() error")
	}
}

func TestSetManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	cm.SetManager(&parser.ServerFileManager{})
}

func TestRegisterCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	cm.RegisterCallback(nil, "")
}

func TestDeregisterCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	cm.DeregisterCallback("")
}

func TestStartFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	if err := cm.Start(); err == nil {
		t.Errorf("filewatcher not sert manager, Start() should error, but not")
	}
}

func TestStartSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	m.EXPECT().FilePath().Return("test").AnyTimes()
	m.EXPECT().RegisterCallback(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().CallOnceSpecific(gomock.Any()).Return(nil).AnyTimes()

	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	cm.SetManager(&parser.ServerFileManager{})
	if err := cm.Start(); err != nil {
		t.Errorf("Start() error = %v", err)
	}
}

func TestStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockFileWatcher(ctrl)
	m.EXPECT().DeregisterCallback(gomock.Any()).AnyTimes()
	cm, err := NewConfigMonitor("test", m)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	cm.Stop()
}

func TestEntireProgress(t *testing.T) {
	var callbackKey string = "test"
	var filepath string = "test.json"
	mapKey := []string{"Test1", "Test2"}
	// create a file watcher object
	fw, err := filewatcher.NewFileWatcher(filepath)
	if err != nil {
		t.Errorf("NewFileWatcher() error = %v", err)
	}
	// start watching file changes
	if err = fw.StartWatching(); err != nil {
		t.Errorf("StartWatching() error = %v", err)
	}

	// create a config monitor object
	cm, err := NewConfigMonitor(mapKey[0], fw)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	// set manager
	cm.SetManager(&parser.ServerFileManager{})
	cm.RegisterCallback(testCallback, callbackKey)
	// start monitoring
	if err = cm.Start(); err != nil {
		t.Errorf("Start() error = %v", err)
	}
	t.Log("call specific ConfigManager0")
	fw.CallOnceSpecific(mapKey[0])

	// create a config monitor object
	cm1, err := NewConfigMonitor(mapKey[1], fw)
	if err != nil {
		t.Errorf("NewConfigMonitor() error = %v", err)
	}
	// set manager
	cm1.SetManager(&parser.ServerFileManager{})
	cm1.RegisterCallback(testCallback, callbackKey)
	// start monitoring
	if err = cm1.Start(); err != nil {
		t.Errorf("Start() error = %v", err)
	}
	t.Log("call specific ConfigManager1")
	fw.CallOnceSpecific(mapKey[1])

	t.Log("call all ConfigManager0 and ConfigManager1")
	fw.CallOnceAll()

	cm.Stop()
	cm1.Stop()
	fw.StopWatching()
}

func testCallback() {
	fmt.Println("testCallback callback once")
}
