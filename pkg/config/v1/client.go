// Copyright 2023 The frp Authors
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

package v1

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrp/tiny-frpc/pkg/util"
)

type ClientConfig struct {
	ClientCommonConfig

	Proxies  []TypedProxyConfig   `json:"proxies,omitempty"`
	Visitors []TypedVisitorConfig `json:"visitors,omitempty"`
}

type ClientCommonConfig struct {
	Auth AuthClientConfig `json:"auth,omitempty"`
	// User specifies a prefix for proxy names to distinguish them from other
	// clients. If this value is not "", proxy names will automatically be
	// changed to "{user}.{proxy_name}".
	User string `json:"user,omitempty"`

	// ServerAddr specifies the address of the server to connect to. By
	// default, this value is "0.0.0.0".
	ServerAddr string `json:"serverAddr,omitempty"`

	// SSHServerPort specifies the port to connect to the ssh server on. By default,
	// this value is 2200.
	SSHServerPort int `json:"sshServerPort,omitempty"`

	// Include other config files for proxies.
	IncludeConfigFiles []string `json:"includes,omitempty"`
}

func (c *ClientCommonConfig) Complete() {
	c.ServerAddr = util.EmptyOr(c.ServerAddr, "0.0.0.0")
	c.SSHServerPort = util.EmptyOr(c.SSHServerPort, 2200)
	c.Auth.Complete()
}

type AuthClientConfig struct {
	// Method specifies what authentication method to use to
	// authenticate frpc with frps. If "token" is specified - token will be
	// read into login message. If "oidc" is specified - OIDC (Open ID Connect)
	// token will be issued using OIDC settings. By default, this value is "token".
	Method AuthMethod `json:"method,omitempty"`

	// Token specifies the authorization token used to create keys to be sent
	// to the server. The server must have a matching token for authorization
	// to succeed.  By default, this value is "".
	Token string `json:"token,omitempty"`
}

func (c *AuthClientConfig) Complete() {
	c.Method = util.EmptyOr(c.Method, "token")
}

// Contains returns true if an element is present in a collection.
func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}

	return false
}

func ValidateClientCommonConfig(c *ClientCommonConfig) (Warning, error) {
	var (
		warnings Warning
		errs     error
	)

	if !Contains(SupportedAuthMethods, c.Auth.Method) {
		errs = AppendError(errs, fmt.Errorf("invalid auth method, optional values are %v", SupportedAuthMethods))
	}

	for _, f := range c.IncludeConfigFiles {
		absDir, err := filepath.Abs(filepath.Dir(f))
		if err != nil {
			errs = AppendError(errs, fmt.Errorf("include: parse directory of %s failed: %v", f, err))
			continue
		}
		if _, err := os.Stat(absDir); os.IsNotExist(err) {
			errs = AppendError(errs, fmt.Errorf("include: directory of %s not exist", f))
		}
	}
	return warnings, errs
}

func ValidateAllClientConfig(c *ClientCommonConfig, proxyCfgs []ProxyConfigurer, visitorCfgs []VisitorConfigurer) (Warning, error) {
	var warnings Warning
	if c != nil {
		warning, err := ValidateClientCommonConfig(c)
		warnings = AppendError(warnings, warning)
		if err != nil {
			return warnings, err
		}
	}

	for _, c := range proxyCfgs {
		if err := ValidateProxyConfigurerForClient(c); err != nil {
			return warnings, fmt.Errorf("proxy %s: %v", c.GetBaseConfig().Name, err)
		}
	}

	for _, c := range visitorCfgs {
		if err := ValidateVisitorConfigurer(c); err != nil {
			return warnings, fmt.Errorf("visitor %s: %v", c.GetBaseConfig().Name, err)
		}
	}
	return warnings, nil
}
