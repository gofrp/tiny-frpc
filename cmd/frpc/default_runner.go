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
	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/model"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

func NewRunner(commonCfg *v1.ClientCommonConfig, pxyCfg []v1.ProxyConfigurer, vCfg []v1.VisitorConfigurer) (model.Runner, error) {
	log.Infof("init default runner")
	return defaultRunner{}, nil
}

type defaultRunner struct{}

func (r defaultRunner) Run() (err error) {
	return
}

func (r defaultRunner) Close() (err error) {
	return
}
