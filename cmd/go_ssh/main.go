package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/gofrp/tiny-frpc/pkg/config"
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/gssh"
	"github.com/gofrp/tiny-frpc/pkg/util"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
	"github.com/gofrp/tiny-frpc/pkg/util/version"
)

func main() {
	var (
		cfgFilePath string
		showVersion bool
	)

	flag.StringVar(&cfgFilePath, "c", "frpc.toml", "path to the configuration file")
	flag.BoolVar(&showVersion, "v", false, "version of frpc-gssh")
	flag.Parse()

	if showVersion {
		fmt.Println(version.Full())
		return
	}

	cfg, proxyCfgs, visitorCfgs, _, err := config.LoadClientConfig(cfgFilePath, true)
	if err != nil {
		log.Errorf("load frpc config error: %v", err)
		return
	}

	_, err = v1.ValidateAllClientConfig(cfg, proxyCfgs, visitorCfgs)
	if err != nil {
		log.Errorf("validate frpc config error: %v", err)
		return
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
