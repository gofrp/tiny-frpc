// Copyright 2024 gofrp (https://github.com/gofrp)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"net"
	"reflect"
	"strings"

	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
)

type GoSSHParam struct {
	LocalAddr   string
	ServerAddr  string
	SSHExtraCmd string
}

func ParseFRPCConfigToGoSSHParam(
	cfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
) (res []GoSSHParam) {
	res = make([]GoSSHParam, 0)

	for _, pv := range proxyCfgs {
		switch pv.GetProxyType() {
		case v1.ProxyTypeTCP, v1.ProxyTypeHTTP, v1.ProxyTypeHTTPS, v1.ProxyTypeTCPMUX, v1.ProxyTypeSTCP:

			res = append(res, GoSSHParam{
				LocalAddr:  pv.GetLocalServerAddr(),
				ServerAddr: net.JoinHostPort(cfg.ServerAddr, fmt.Sprintf("%d", cfg.ServerPort)),

				SSHExtraCmd: genSSHExtraCmd(*cfg, pv),
			})
		default:
			panic("invalid proxy type: " + pv.GetProxyType())
		}
	}

	return
}

// ParseFRPCConfigToSSHCmd parse standard frpc config to standard ssh commands
func ParseFRPCConfigToSSHCmd(
	cfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
) []string {
	res := make([]string, 0)

	for _, pv := range proxyCfgs {
		switch pv.GetProxyType() {
		case v1.ProxyTypeTCP, v1.ProxyTypeHTTP, v1.ProxyTypeHTTPS, v1.ProxyTypeTCPMUX, v1.ProxyTypeSTCP:
			cmd := genFullSSHCmd(*cfg, pv)
			res = append(res, cmd)
		default:
			panic("invalid proxy type: " + pv.GetProxyType())
		}
	}

	// visitorCfgs now is useless but reserved now

	return res
}

// ssh raw cmd contains 5 parts
// part1: "ssh v0@%v -p %v"
// part2: "-R :80:%v"
// part3: "{proxy type}"
// part4: "{proxy related args}"
// part5: "{auth and users args}"
func genFullSSHCmd(c v1.ClientCommonConfig, pc v1.ProxyConfigurer) string {
	return strings.TrimSpace(fmt.Sprintf("%v %v %v %v %v", genDialedCmd(c), genReverseCmd(pc), genProxyTypeCmd(pc), genProxyCmd(pc), genAuthCmd(c)))
}

func genSSHExtraCmd(c v1.ClientCommonConfig, pc v1.ProxyConfigurer) string {
	return strings.TrimSpace(fmt.Sprintf("%v %v %v", genProxyTypeCmd(pc), genProxyCmd(pc), genAuthCmd(c)))
}

func genDialedCmd(c v1.ClientCommonConfig) string {
	return fmt.Sprintf("ssh v0@%v -p %v", c.ServerAddr, c.ServerPort)
}

func genReverseCmd(pc v1.ProxyConfigurer) string {
	return fmt.Sprintf("-R :80:%v", pc.GetLocalServerAddr())
}

func genProxyTypeCmd(pc v1.ProxyConfigurer) string {
	return string(pc.GetProxyType())
}

func genProxyCmd(pc v1.ProxyConfigurer) string {
	np := v1.NewProxy{}
	pc.MarshalToMsg(&np)

	var cmd string

	t := reflect.TypeOf(np)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := reflect.ValueOf(np).Field(i)

		if !value.IsZero() && field.Tag.Get("flag") != "" {
			cmd += fmt.Sprintf("--%s %v ", field.Tag.Get("flag"), value)
		}
	}

	return strings.TrimSpace(cmd)
}

func genAuthCmd(c v1.ClientCommonConfig) string {
	res := ""

	if c.User != "" {
		res += fmt.Sprintf("--user %v ", c.User)
	}
	if c.Auth.Method != "" && c.Auth.Token != "" {
		res += fmt.Sprintf("--token %v ", c.Auth.Token)
	}
	return strings.TrimSpace(res)
}
