package main

import (
	"sync"

	"github.com/blizard/frpc-slim/pkg/config"
	v1 "github.com/blizard/frpc-slim/pkg/config/v1"
	"github.com/blizard/frpc-slim/pkg/gssh"
	"github.com/blizard/frpc-slim/pkg/util"
	"github.com/blizard/frpc-slim/pkg/util/log"
)

func main() {
	cfgFilePath := "./frpc.toml"

	cfg, proxyCfgs, visitorCfgs, _, err := config.LoadClientConfig(cfgFilePath, true)
	if err != nil {
		panic(err)
	}

	_, err = v1.ValidateAllClientConfig(cfg, proxyCfgs, visitorCfgs)
	if err != nil {
		panic(err)
	}

	log.Infof("common cfg: %v, proxy cfg: %v, visitor cfg: %v", util.JSONEncode(cfg), util.JSONEncode(proxyCfgs), util.JSONEncode(visitorCfgs))

	goSSHParams := config.ParseFRPCConfigToGoSSHParam(cfg, proxyCfgs, visitorCfgs)

	log.Infof("ssh cmds len_num: %v", len(goSSHParams))

	wg := new(sync.WaitGroup)

	for _, cmd := range goSSHParams {
		wg.Add(1)

		go func(cmd config.GoSSHParam) {
			defer wg.Done()

			log.Infof("start to run %v", cmd)

			tc := gssh.NewTunnelClient(cmd.LocalAddr, cmd.ServerAddr, cmd.SSHExtraCmd)

			err := tc.Start()
			if err != nil {
				log.Errorf("cmd: %v run error: %v", cmd, err)
				return
			}
		}(cmd)
	}

	wg.Wait()

	log.Infof("stopping process calling native ssh to frps, exit...")
}
