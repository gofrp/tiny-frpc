// Copyright 2024 gofrp (https://github.com/gofrp)
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

//go:build nssh
// +build nssh

package main

import (
	"context"
	"sync"

	"github.com/gofrp/tiny-frpc/pkg/config"
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/nssh"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

type NativeSSHRun struct {
	commonCfg *v1.ClientCommonConfig
	pxyCfg    []v1.ProxyConfigurer
	vCfg      []v1.VisitorConfigurer

	wg *sync.WaitGroup
	mu *sync.RWMutex

	cws map[int]*nssh.CmdWrapper
	ctx context.Context
	cancelFunc context.CancelFunc
}

func (nr *NativeSSHRun) New(commonCfg *v1.ClientCommonConfig, pxyCfg []v1.ProxyConfigurer, vCfg []v1.VisitorConfigurer) error {
	log.Infof("init native ssh runner")

	nr.ctx, nr.cancelFunc = context.WithCancel(context.Background())

	runner = &NativeSSHRun{
		commonCfg: commonCfg,
		pxyCfg:    pxyCfg,
		vCfg:      vCfg,

		wg: new(sync.WaitGroup),
		mu: &sync.RWMutex{},

		cws: make(map[int]*nssh.CmdWrapper, 0),
	}
	return nil
}

func (nr *NativeSSHRun) Run() error {
	cmdParams := config.ParseFRPCConfigToSSHCmd(nr.commonCfg, nr.pxyCfg, nr.vCfg)

	log.Infof("proxy total len: %v", len(cmdParams))

	for i, cmd := range cmdParams {
		nr.wg.Add(1)

		go func(cmd string, idx int) {
			defer nr.wg.Done()

			log.Infof("start to run: %v", cmd)

			cmdWrapper := nssh.NewCmdWrapper(nr.ctx, cmd)

			nr.mu.Lock()
			nr.cws[idx] = cmdWrapper
			nr.mu.Unlock()

			cmdWrapper.ExecuteCommand(nr.ctx)
		}(cmd, i)
	}

	nr.wg.Wait()

	log.Infof("stopping native ssh tunnel to frps")

	return nil
}

func (nr *NativeSSHRun) Close() error {
	nr.cancelFunc() // Ensure all goroutines are signaled to stop

	nr.mu.Lock()
	defer nr.mu.Unlock()

	for _, c := range nr.cws {
		c.Close()
	}

	return nil
}

func init() {
	runner = &NativeSSHRun{}
}
