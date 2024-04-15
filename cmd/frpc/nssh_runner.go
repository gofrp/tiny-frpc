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
}

func (nr *NativeSSHRun) New(commonCfg *v1.ClientCommonConfig, pxyCfg []v1.ProxyConfigurer, vCfg []v1.VisitorConfigurer) error {
	log.Infof("init native ssh runner")

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
			ctx := context.Background()

			log.Infof("start to run: %v", cmd)

			cmdWrapper := nssh.NewCmdWrapper(ctx, cmd)

			nr.mu.Lock()
			nr.cws[idx] = cmdWrapper
			nr.mu.Unlock()

			cmdWrapper.ExecuteCommand(ctx)
		}(cmd, i)
	}

	nr.wg.Wait()

	log.Infof("stopping native ssh tunnel to frps")

	return nil
}

func (nr *NativeSSHRun) Close() error {
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
