# config-local

[English](https://github.com/kitex-contrib/config-file/blob/main/README.md)

读取、加载并监听本地配置文件

## 使用说明

### 基本使用

#### 服务端

```go
package main

import (
	"context"
	"log"

	"github.com/cloudwego/kitex-examples/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/kitex_gen/api/echo"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexserver "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-file/filewatcher"
	fileserver "github.com/kitex-contrib/config-file/server"
)

var _ api.Echo = &EchoImpl{}

const (
	filepath    = "kitex_server.json"
	key         = "ServiceName"
	serviceName = "ServiceName"
)

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the Echo interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	klog.Info("echo called")
	return &api.Response{Message: req.Message}, nil
}

func main() {
	klog.SetLevel(klog.LevelDebug)

	// create a file watcher object
	fw, err := filewatcher.NewFileWatcher(filepath)
	if err != nil {
		panic(err)
	}
	// start watching file changes
	if err = fw.StartWatching(); err != nil {
		panic(err)
	}
	defer fw.StopWatching()

	svr := echo.NewServer(
		new(EchoImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
		server.WithSuite(fileserver.NewSuite(key, fw)), // add watcher
	)
	if err := svr.Run(); err != nil {
		log.Println("server stopped with error:", err)
	} else {
		log.Println("server stopped")
	}
}
```

#### 客户端

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/cloudwego/kitex-examples/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/kitex_gen/api/echo"
	kitexclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	fileclient "github.com/kitex-contrib/config-file/client"
	"github.com/kitex-contrib/config-file/filewatcher"
)

const (
	filepath    = "kitex_client.json"
	key         = "ClientName/ServiceName"
	serviceName = "ServiceName"
	clientName  = "ClientName"
)

func main() {
	klog.SetLevel(klog.LevelDebug)

	// create a file watcher object
	fw, err := filewatcher.NewFileWatcher(filepath)
	if err != nil {
		panic(err)
	}
	// start watching file changes
	if err = fw.StartWatching(); err != nil {
		panic(err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		fw.StopWatching()
		os.Exit(1)
	}()

	client, err := echo.NewClient(
		serviceName,
		kitexclient.WithHostPorts("0.0.0.0:8888"),
		kitexclient.WithSuite(fileclient.NewSuite(serviceName, key, fw)),
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		req := &api.Request{Message: "my request"}
		resp, err := client.Echo(context.Background(), req)
		if err != nil {
			klog.Errorf("take request error: %v", err)
		} else {
			klog.Infof("receive response %v", resp)
		}
		time.Sleep(time.Second * 10)
	}
}

```

#### 治理策略
> 服务名称为 ServiceName，客户端名称为 ClientName

##### 限流：Category=limit
> 限流目前只支持服务端，所以只需要设置服务端的 ServiceName。

[JSON Schema](https://github.com/cloudwego/kitex/blob/develop/pkg/limiter/item_limiter.go#L33)

|字段|说明|
|----|----|
|connection_limit|最大并发数量|
|qps_limit|每 100ms 内的最大请求数量|

样例:
```json
{
    "ServiceName": {
        "limit": {
            "connection_limit": 300,
            "qps_limit": 200
        }
    }
}
```

注：

- 限流配置的粒度是 Server 全局，不分 client、method
- 「未配置」或「取值为 0」表示不开启
- connection_limit 和 qps_limit 可以独立配置，例如 connection_limit = 100, qps_limit = 0

##### 重试：Category=retry
[JSON Schema](https://github.com/cloudwego/kitex/blob/develop/pkg/retry/policy.go#L63)

|参数|说明|
|----|----|
|type| 0: failure_policy 1: backup_policy|
|failure_policy.backoff_policy| 可以设置的策略： `fixed` `none` `random` |

样例：
```json
{
    "ClientName/ServiceName": {
        "retry": {
            "*": {
                "enable": true,
                "type": 0,
                "failure_policy": {
                    "stop_policy": {
                        "max_retry_times": 3,
                        "max_duration_ms": 2000,
                        "cb_policy": {
                            "error_rate": 0.2
                        }
                    }
                }
            },
            "Echo": {
                "enable": true,
                "type": 1,
                "backup_policy": {
                    "retry_delay_ms": 200,
                    "stop_policy": {
                        "max_retry_times": 2,
                        "max_duration_ms": 1000,
                        "cb_policy": {
                            "error_rate": 0.3
                        }
                    }
                }
            }
        }
    }
}
```
注：retry.Container 内置支持用 * 通配符指定默认配置（详见 [getRetryer](https://github.com/cloudwego/kitex/blob/v0.5.1/pkg/retry/retryer.go#L240) 方法）

##### 超时：Category=rpc_timeout

[JSON Schema](https://github.com/cloudwego/kitex/blob/develop/pkg/rpctimeout/item_rpc_timeout.go#L42)

样例：
```json
{
    "ClientName/ServiceName": {
        "timeout": {
            "*": {
                "conn_timeout_ms": 100,
                "rpc_timeout_ms": 2000
            },
            "Pay": {
                "conn_timeout_ms": 50,
                "rpc_timeout_ms": 1000
            }
        },
    }
}
```

##### 熔断: Category=circuit_break

[JSON Schema](https://github.com/cloudwego/kitex/blob/develop/pkg/circuitbreak/item_circuit_breaker.go#L30)

|参数|说明|
|----|----|
|min_sample|最小的统计样本数|

样例：
echo 方法使用下面的配置（0.3、100），其他方法使用全局默认配置（0.5、200）
```json
{
    "ClientName/ServiceName": {
        "circuitbreaker": {
            "Echo": {
                "enable": true,
                "err_rate": 0.3,
                "min_sample": 100
            }
        },
    }
}
```
注：kitex 的熔断实现目前不支持修改全局默认配置（详见 [initServiceCB](https://github.com/cloudwego/kitex/blob/v0.5.1/pkg/circuitbreak/cbsuite.go#L195)）
### 更多信息

更多示例请参考 [example](https://github.com/kitex-contrib/config-file/tree/main/example)


## 注意事项

### 自定义键

对于客户端配置，您应该将它们的所有配置写入同一对`$UserServiceName/$ServerServiceName`中，例如

```json
{
    "ClientName/ServiceName": {
        "timeout": {
            "*": {
                "conn_timeout_ms": 100,
                "rpc_timeout_ms": 2000
            },
            "Pay": {
                "conn_timeout_ms": 50,
                "rpc_timeout_ms": 1000
            }
        },
        "circuitbreaker": {
            "Echo": {
                "enable": true,
                "err_rate": 0.3,
                "min_sample": 100
            }
        },
        "retry": {
            "*": {
                "enable": true,
                "type": 0,
                "failure_policy": {
                    "stop_policy": {
                        "max_retry_times": 3,
                        "max_duration_ms": 2000,
                        "cb_policy": {
                            "error_rate": 0.2
                        }
                    }
                }
            },
            "Echo": {
                "enable": true,
                "type": 1,
                "backup_policy": {
                    "retry_delay_ms": 200,
                    "stop_policy": {
                        "max_retry_times": 2,
                        "max_duration_ms": 1000,
                        "cb_policy": {
                            "error_rate": 0.3
                        }
                    }
                }
            }
        }
    }
}
```