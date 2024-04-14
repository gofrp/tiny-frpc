package main

import (
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/model"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

var runner model.Runner = defaultRunner{}

type defaultRunner struct{}

func (r defaultRunner) New(commonCfg *v1.ClientCommonConfig, pxyCfg []v1.ProxyConfigurer, vCfg []v1.VisitorConfigurer) (err error) {
	log.Infof("init default runner")
	return
}

func (r defaultRunner) Run() (err error) {
	return
}

func (r defaultRunner) Close() (err error) {
	return
}
