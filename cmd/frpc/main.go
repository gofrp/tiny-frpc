package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofrp/tiny-frpc/pkg/config"
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/model"
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
	flag.BoolVar(&showVersion, "v", false, "version of tiny-frpc")
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

	err = runner.New(cfg, proxyCfgs, visitorCfgs)
	if err != nil {
		log.Errorf("new runner error: %v", err)
		return
	}

	go handleTermSignal(runner)

	err = runner.Run()
	if err != nil {
		log.Errorf("run error: %v", err)
		return
	}

	time.Sleep(time.Millisecond * 10)
	log.Infof("process exit...")
}

func handleTermSignal(run model.Runner) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	v := <-ch
	log.Infof("get signal term: %v, gracefully shutdown", v)
	run.Close()
}
