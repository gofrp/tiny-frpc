package main

import (
	"context"
	"sync"

	"github.com/blizard/frpc-slim/pkg/config"
	v1 "github.com/blizard/frpc-slim/pkg/config/v1"
	"github.com/blizard/frpc-slim/pkg/nssh"
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

	sshCmds := config.ParseFRPCConfigToSSHCmd(cfg, proxyCfgs, visitorCfgs)

	log.Infof("ssh cmds len_num: %v", len(sshCmds))

	closeCh := make(chan struct{})
	wg := new(sync.WaitGroup)

	for _, cmd := range sshCmds {
		wg.Add(1)

		go func(cmd string) {
			defer wg.Done()
			ctx := context.Background()

			log.Infof("start to run %v", cmd)

			task := nssh.NewCmdWrapper(ctx, cmd, closeCh)
			task.ExecuteCommand(ctx)
		}(cmd)
	}

	wg.Wait()
	close(closeCh)

	log.Infof("stopping process calling native ssh to frps, exit...")
}
