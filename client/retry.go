package client

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/kitex-contrib/config-file/utils"
)

// WithRetryPolicy returns a server.Option that sets the retry policies for the client
func WithRetryPolicy(watcher *ConfigWatcher) client.Option {
	return client.WithRetryContainer(initRetryContainer(watcher))
}

// initRetryOptions init the retry container
func initRetryContainer(watcher *ConfigWatcher) *retry.Container {
	retryContainer := retry.NewRetryContainer()

	ts := utils.ThreadSafeSet{}

	onChangeCallback := func() {
		// the key is method name, wildcard "*" can match anything.
		rcs := watcher.Config().Retry

		if rcs == nil {
			klog.Warnf("[local] %s file retry config: failed as the config not found, skip...", watcher.Key())
			return
		}

		set := utils.Set{}

		for method, policy := range rcs {
			set[method] = true

			if policy.BackupPolicy != nil && policy.FailurePolicy != nil {
				klog.Warnf("[local] %s client policy for method %s BackupPolicy and FailurePolicy must not be set at same time",
					watcher.Key(), method)
				continue
			}

			if policy.BackupPolicy == nil && policy.FailurePolicy == nil {
				klog.Warnf("[local] %s client policy for method %s BackupPolicy and FailurePolicy must not be empty at same time",
					watcher.Key(), method)
				continue
			}

			retryContainer.NotifyPolicyChange(method, *policy)
		}

		for _, method := range ts.DiffAndEmplace(set) {
			retryContainer.DeletePolicy(method)
		}
	}

	watcher.AddCallback(onChangeCallback)

	return retryContainer
}
