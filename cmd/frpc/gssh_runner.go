//go:build gssh
// +build gssh

package main

import (
	"sync"

	"github.com/gofrp/tiny-frpc/pkg/config"
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/gssh"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

type GoSSHRun struct {
	commonCfg *v1.ClientCommonConfig
	pxyCfg    []v1.ProxyConfigurer
	vCfg      []v1.VisitorConfigurer

	wg *sync.WaitGroup
	mu *sync.RWMutex

	tcs map[int]*gssh.TunnelClient
}

func (gr *GoSSHRun) New(commonCfg *v1.ClientCommonConfig, pxyCfg []v1.ProxyConfigurer, vCfg []v1.VisitorConfigurer) error {
	log.Infof("init go ssh runner")

	runner = &GoSSHRun{
		commonCfg: commonCfg,
		pxyCfg:    pxyCfg,
		vCfg:      vCfg,

		wg: new(sync.WaitGroup),
		mu: &sync.RWMutex{},

		tcs: make(map[int]*gssh.TunnelClient, 0),
	}
	return nil
}

func (gr *GoSSHRun) Run() error {
	params := config.ParseFRPCConfigToGoSSHParam(gr.commonCfg, gr.pxyCfg, gr.vCfg)

	log.Infof("proxy total len: %v", len(params))

	for i, cmd := range params {
		gr.wg.Add(1)

		go func(cmd config.GoSSHParam, idx int) {
			defer gr.wg.Done()

			log.Infof("start to run: %v", cmd)

			tc, err := gssh.NewTunnelClient(cmd.LocalAddr, cmd.ServerAddr, cmd.SSHExtraCmd)
			if err != nil {
				log.Errorf("new ssh tunnel client error: %v", err)
				return
			}

			gr.mu.Lock()
			gr.tcs[idx] = tc
			gr.mu.Unlock()

			err = tc.Start()
			if err != nil {
				log.Errorf("cmd: %v run error: %v", cmd, err)

				gr.mu.Lock()
				delete(gr.tcs, idx)
				gr.mu.Unlock()

				return
			}
		}(cmd, i)
	}

	gr.wg.Wait()

	log.Infof("stopping ssh tunnel to frps")
	return nil
}

func (gr *GoSSHRun) Close() error {
	gr.mu.Lock()
	defer gr.mu.Unlock()

	for _, tc := range gr.tcs {
		tc.Close()
	}
	return nil
}

func init() {
	runner = &GoSSHRun{}
}
