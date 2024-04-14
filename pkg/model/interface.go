package model

import (
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
)

type Runner interface {
	New(commonCfg *v1.ClientCommonConfig, pxyCfg []v1.ProxyConfigurer, vCfg []v1.VisitorConfigurer) error
	Run() error
	Close() error
}
