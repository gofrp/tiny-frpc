// Copyright 2024 gofrp (https://github.com/gofrp)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	cfgFilePath, showVersion := parseArgs()

	if showVersion {
		fmt.Println(version.Full())
		return
	}

	cfg, proxyCfgs, visitorCfgs, err := loadConfig(cfgFilePath)
	if err != nil {
		log.Errorf("load frpc config error: %v", err)
		return
	}

	runner, err := NewRunner(cfg, proxyCfgs, visitorCfgs)
	if err != nil {
		log.Errorf("new runner error: %v", err)
		return
	}

	setupSignalHandler(runner)

	err = runner.Run()
	if err != nil {
		log.Errorf("run error: %v", err)
		return
	}

	time.Sleep(time.Millisecond * 10)
	log.Infof("process exit...")
}

func parseArgs() (cfgFilePath string, showVersion bool) {
	flag.StringVar(&cfgFilePath, "c", "frpc.toml", "path to the configuration file")
	flag.BoolVar(&showVersion, "v", false, "version of tiny-frpc")
	flag.Parse()
	return
}

func loadConfig(cfgFilePath string) (*v1.ClientCommonConfig, []v1.ProxyConfigurer, []v1.VisitorConfigurer, error) {
	cfg, proxyCfgs, visitorCfgs, _, err := config.LoadClientConfig(cfgFilePath, true)
	if err != nil {
		return nil, nil, nil, err
	}

	_, err = v1.ValidateAllClientConfig(cfg, proxyCfgs, visitorCfgs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("validate frpc config error: %v", err)
	}

	log.Infof("common cfg: %v, proxy cfg: %v, visitor cfg: %v", util.JSONEncode(cfg), util.JSONEncode(proxyCfgs), util.JSONEncode(visitorCfgs))
	return cfg, proxyCfgs, visitorCfgs, nil
}

func setupSignalHandler(run model.Runner) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		v := <-ch
		log.Infof("get signal term: %v, gracefully shutdown", v)
		run.Close()
	}()
}
