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
