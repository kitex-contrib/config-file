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
