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

package utils

import (
	"errors"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/fsnotify/fsnotify"
)

// FileWatcher is used for file monitoring
type FileWatcher struct {
	filePath string                // The path to the file to be monitored.
	callback func(filepath string) // Custom function to be executed when the file changes.
	watcher  *fsnotify.Watcher     // fsnotify file change watcher.
	done     chan struct{}         // A channel for signaling the watcher to stop.
}

// NewFileWatcher creates a new FileWatcher instance.
func NewFileWatcher(filePath string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	exist, err := PathExists(filePath)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("file [" + filePath + "] not exist")
	}

	fw := &FileWatcher{
		filePath: filePath,
		watcher:  watcher,
		done:     make(chan struct{}),
	}

	return fw, nil
}

func (fw *FileWatcher) FilePath() string { return fw.filePath }

// SetCallback sets the callback function.
func (fw *FileWatcher) AddCallback(callback func(filepath string)) {
	fw.callback = callback
}

// Start starts monitoring file changes.
func (fw *FileWatcher) StartWatching() error {
	err := fw.watcher.Add(fw.filePath)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				klog.Errorf("file watcher panic: %v\n", r)
			}
		}()
		fw.start()
	}()

	return nil
}

// Stop stops monitoring file changes.
func (fw *FileWatcher) StopWatching() {
	close(fw.done)
}

// StartWatching starts monitoring file changes.
func (fw *FileWatcher) start() {
	defer fw.watcher.Close()
	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				fw.callback(fw.filePath)
			}
			if event.Has(fsnotify.Remove) {
				klog.Warnf("file %s is removed, stop watching", fw.filePath)
				fw.StopWatching()
			}
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			klog.Errorf("file watcher meet error: %v\n", err)
		case <-fw.done:
			return
		}
	}
}
