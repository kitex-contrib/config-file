// Copyright 2024 CloudWeGo Authors
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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/cloudwego/kitex-examples/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/kitex_gen/api/echo"
	kitexclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpctimeout"
	fileclient "github.com/kitex-contrib/config-file/client"
	"github.com/kitex-contrib/config-file/filewatcher"
	"github.com/kitex-contrib/config-file/parser"
	"github.com/kitex-contrib/config-file/utils"
	"gopkg.in/ini.v1"
)

const (
	filepath                      = "kitex_client.ini"
	key                           = "ClientName/ServiceName"
	serviceName                   = "ServiceName"
	clientName                    = "echo"
	INI         parser.ConfigType = "ini"
)

// Customed by user
type MyParser struct{}

// one example for custom parser
// if the type of client config is json or yaml,just using default parser
func (p *MyParser) Decode(kind parser.ConfigType, data []byte, config interface{}) error {
	cfg, err := ini.Load(data)
	if err != nil {
		return fmt.Errorf("load config error: %v", err)
	}

	cfm := make(parser.ClientFileManager, 0)
	cfc := &parser.ClientFileConfig{
		Timeout:        make(map[string]*rpctimeout.RPCTimeout, 0),
		Retry:          make(map[string]*retry.Policy, 0),
		Circuitbreaker: make(map[string]*circuitbreak.CBConfig, 0),
	}

	timeout := &rpctimeout.RPCTimeout{}
	circ := &circuitbreak.CBConfig{}
	ret := &retry.Policy{}
	stop := &retry.StopPolicy{}
	cb := &retry.CBPolicy{}

	cfg.Section(key).MapTo(timeout)
	cfg.Section("ClientName/ServiceName.Circuitbreaker.Echo").MapTo(circ)
	cfg.Section("ClientName/ServiceName.Retry.*").MapTo(ret)
	cfg.Section("ClientName/ServiceName.Retry.*.FailurePolicy.StopPolicy").MapTo(stop)
	cfg.Section("ClientName/ServiceName.Retry.*.CBPolicy").MapTo(cb)
	stop.CBPolicy = *cb
	ret.FailurePolicy = &retry.FailurePolicy{
		StopPolicy: *stop,
	}

	cfc.Timeout[clientName] = timeout
	cfc.Circuitbreaker[clientName] = circ
	cfc.Retry[clientName] = ret

	cfm[key] = cfc

	v := config.(*parser.ClientFileManager)
	*v = cfm
	return err
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

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		fw.StopWatching()
		os.Exit(1)
	}()

	// customed by user
	params := &parser.ConfigParam{
		Type: INI,
	}
	opts := &utils.Options{
		CustomParser: &MyParser{},
		CustomParams: params,
	}

	client, err := echo.NewClient(
		serviceName,
		kitexclient.WithHostPorts("0.0.0.0:8888"),
		kitexclient.WithSuite(fileclient.NewSuite(serviceName, key, fw, opts)),
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
