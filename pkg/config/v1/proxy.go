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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/gofrp/tiny-frpc/pkg/util"
)

type ProxyBackend struct {
	// LocalIP specifies the IP address or host name of the backend.
	LocalIP string `json:"localIP,omitempty"`
	// LocalPort specifies the port of the backend.
	LocalPort int `json:"localPort,omitempty"`
}

type DomainConfig struct {
	CustomDomains []string `json:"customDomains,omitempty"`
	SubDomain     string   `json:"subdomain,omitempty"`
}

type ProxyBaseConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
	// metadata info for each proxy
	Metadatas map[string]string `json:"metadatas,omitempty"`
	ProxyBackend
}

func (c *ProxyBaseConfig) GetBaseConfig() *ProxyBaseConfig {
	return c
}

func (c *ProxyBaseConfig) Complete(namePrefix string) {
	c.Name = util.Ternary(namePrefix == "", "", namePrefix+".") + c.Name
	c.LocalIP = util.EmptyOr(c.LocalIP, "127.0.0.1")
}

func (c *ProxyBaseConfig) MarshalToMsg(m *NewProxy) {
	m.ProxyName = c.Name
	m.ProxyType = c.Type
	m.Metas = c.Metadatas
}

type TypedProxyConfig struct {
	Type string `json:"type"`
	ProxyConfigurer
}

func (c *TypedProxyConfig) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		return errors.New("type is required")
	}

	typeStruct := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(b, &typeStruct); err != nil {
		return err
	}

	c.Type = typeStruct.Type
	configurer := NewProxyConfigurerByType(ProxyType(typeStruct.Type))
	if configurer == nil {
		return fmt.Errorf("unknown proxy type: %s", typeStruct.Type)
	}
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	if DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(configurer); err != nil {
		return err
	}
	c.ProxyConfigurer = configurer
	return nil
}

type ProxyConfigurer interface {
	Complete(namePrefix string)
	GetBaseConfig() *ProxyBaseConfig
	// MarshalToMsg marshals this config into a NewProxy message. This
	// function will be called on the frpc side.
	MarshalToMsg(*NewProxy)
	GetLocalServerAddr() string
	GetProxyType() ProxyType
}

type ProxyType string

const (
	ProxyTypeTCP    ProxyType = "tcp"
	ProxyTypeTCPMUX ProxyType = "tcpmux"
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeHTTPS  ProxyType = "https"
	ProxyTypeSTCP   ProxyType = "stcp"
)

var proxyConfigTypeMap = map[ProxyType]reflect.Type{
	ProxyTypeTCP:    reflect.TypeOf(TCPProxyConfig{}),
	ProxyTypeHTTP:   reflect.TypeOf(HTTPProxyConfig{}),
	ProxyTypeHTTPS:  reflect.TypeOf(HTTPSProxyConfig{}),
	ProxyTypeTCPMUX: reflect.TypeOf(TCPMuxProxyConfig{}),
	ProxyTypeSTCP:   reflect.TypeOf(STCPProxyConfig{}),
}

func NewProxyConfigurerByType(proxyType ProxyType) ProxyConfigurer {
	v, ok := proxyConfigTypeMap[proxyType]
	if !ok {
		return nil
	}
	pc := reflect.New(v).Interface().(ProxyConfigurer)
	pc.GetBaseConfig().Type = string(proxyType)
	return pc
}

var _ ProxyConfigurer = &TCPProxyConfig{}

type TCPProxyConfig struct {
	ProxyBaseConfig

	RemotePort int `json:"remotePort,omitempty"`
}

func (c *TCPProxyConfig) MarshalToMsg(m *NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.RemotePort = c.RemotePort
}

func (c *TCPProxyConfig) GetProxyType() ProxyType {
	return ProxyTypeTCP
}

func (c *TCPProxyConfig) GetLocalServerAddr() string {
	return net.JoinHostPort(c.LocalIP, strconv.Itoa(c.LocalPort))
}

var _ ProxyConfigurer = &HTTPProxyConfig{}

type HTTPProxyConfig struct {
	ProxyBaseConfig
	DomainConfig

	Locations         []string         `json:"locations,omitempty"`
	HTTPUser          string           `json:"httpUser,omitempty"`
	HTTPPassword      string           `json:"httpPassword,omitempty"`
	HostHeaderRewrite string           `json:"hostHeaderRewrite,omitempty"`
	RequestHeaders    HeaderOperations `json:"requestHeaders,omitempty"`
	RouteByHTTPUser   string           `json:"routeByHTTPUser,omitempty"`
}

func (c *HTTPProxyConfig) MarshalToMsg(m *NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.CustomDomains = c.CustomDomains
	m.SubDomain = c.SubDomain
	m.Locations = c.Locations
	m.HostHeaderRewrite = c.HostHeaderRewrite
	m.HTTPUser = c.HTTPUser
	m.HTTPPwd = c.HTTPPassword
	m.Headers = c.RequestHeaders.Set
	m.RouteByHTTPUser = c.RouteByHTTPUser
}

func (c *HTTPProxyConfig) GetProxyType() ProxyType {
	return ProxyTypeHTTP
}

func (c *HTTPProxyConfig) GetLocalServerAddr() string {
	return net.JoinHostPort(c.LocalIP, strconv.Itoa(c.LocalPort))
}

var _ ProxyConfigurer = &HTTPSProxyConfig{}

type HTTPSProxyConfig struct {
	ProxyBaseConfig
	DomainConfig
}

func (c *HTTPSProxyConfig) MarshalToMsg(m *NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.CustomDomains = c.CustomDomains
	m.SubDomain = c.SubDomain
}

func (c *HTTPSProxyConfig) GetProxyType() ProxyType {
	return ProxyTypeHTTPS
}

func (c *HTTPSProxyConfig) GetLocalServerAddr() string {
	return net.JoinHostPort(c.LocalIP, strconv.Itoa(c.LocalPort))
}

type TCPMultiplexerType string

const (
	TCPMultiplexerHTTPConnect TCPMultiplexerType = "httpconnect"
)

var _ ProxyConfigurer = &TCPMuxProxyConfig{}

type TCPMuxProxyConfig struct {
	ProxyBaseConfig
	DomainConfig

	HTTPUser        string `json:"httpUser,omitempty"`
	HTTPPassword    string `json:"httpPassword,omitempty"`
	RouteByHTTPUser string `json:"routeByHTTPUser,omitempty"`
	Multiplexer     string `json:"multiplexer,omitempty"`
}

func (c *TCPMuxProxyConfig) MarshalToMsg(m *NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.CustomDomains = c.CustomDomains
	m.SubDomain = c.SubDomain
	m.Multiplexer = c.Multiplexer
	m.HTTPUser = c.HTTPUser
	m.HTTPPwd = c.HTTPPassword
	m.RouteByHTTPUser = c.RouteByHTTPUser
}

func (c *TCPMuxProxyConfig) GetProxyType() ProxyType {
	return ProxyTypeTCPMUX
}

func (c *TCPMuxProxyConfig) GetLocalServerAddr() string {
	return net.JoinHostPort(c.LocalIP, strconv.Itoa(c.LocalPort))
}

var _ ProxyConfigurer = &STCPProxyConfig{}

type STCPProxyConfig struct {
	ProxyBaseConfig

	Secretkey  string   `json:"secretKey,omitempty"`
	AllowUsers []string `json:"allowUsers,omitempty"`
}

func (c *STCPProxyConfig) MarshalToMsg(m *NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.Sk = c.Secretkey
	m.AllowUsers = c.AllowUsers
}

func (c *STCPProxyConfig) GetProxyType() ProxyType {
	return ProxyTypeSTCP
}

func (c *STCPProxyConfig) GetLocalServerAddr() string {
	return net.JoinHostPort(c.LocalIP, strconv.Itoa(c.LocalPort))
}

type NewProxy struct {
	ProxyName string `flag:"proxy-name"`
	ProxyType string

	Metas map[string]string `flag:"metas"`

	// tcp
	RemotePort int `flag:"remote-port"`

	// http and https only
	CustomDomains     []string `flag:"custom-domain"`
	SubDomain         string   `flag:"sd"`
	Locations         []string `flag:"locations"`
	HTTPUser          string   `flag:"http-user"`
	HTTPPwd           string   `flag:"http-pwd"`
	HostHeaderRewrite string   `flag:"host-header-rewrite"`

	// stcp
	Sk         string   `flag:"sk"`
	AllowUsers []string `flag:"allow-users"`

	// tcpmux
	Multiplexer string `flag:"mux"`

	// TODO deprecated ?
	Headers         map[string]string `flag:"headers"`
	RouteByHTTPUser string
}

func validateProxyBaseConfigForClient(c *ProxyBaseConfig) error {
	if c.Name == "" {
		return errors.New("name should not be empty")
	}

	return nil
}

func validateDomainConfigForClient(c *DomainConfig) error {
	if c.SubDomain == "" && len(c.CustomDomains) == 0 {
		return errors.New("subdomain and custom domains should not be both empty")
	}
	return nil
}

func ValidateProxyConfigurerForClient(c ProxyConfigurer) error {
	base := c.GetBaseConfig()
	if err := validateProxyBaseConfigForClient(base); err != nil {
		return err
	}

	switch v := c.(type) {
	case *TCPProxyConfig:
		return validateTCPProxyConfigForClient(v)
	case *TCPMuxProxyConfig:
		return validateTCPMuxProxyConfigForClient(v)
	case *HTTPProxyConfig:
		return validateHTTPProxyConfigForClient(v)
	case *HTTPSProxyConfig:
		return validateHTTPSProxyConfigForClient(v)
	case *STCPProxyConfig:
		return validateSTCPProxyConfigForClient(v)
	}
	return errors.New("unknown proxy config type")
}

func validateTCPProxyConfigForClient(c *TCPProxyConfig) error {
	return nil
}

func validateTCPMuxProxyConfigForClient(c *TCPMuxProxyConfig) error {
	if err := validateDomainConfigForClient(&c.DomainConfig); err != nil {
		return err
	}

	if !Contains([]string{string(TCPMultiplexerHTTPConnect)}, c.Multiplexer) {
		return fmt.Errorf("not support multiplexer: %s", c.Multiplexer)
	}
	return nil
}

func validateHTTPProxyConfigForClient(c *HTTPProxyConfig) error {
	return validateDomainConfigForClient(&c.DomainConfig)
}

func validateHTTPSProxyConfigForClient(c *HTTPSProxyConfig) error {
	return validateDomainConfigForClient(&c.DomainConfig)
}

func validateSTCPProxyConfigForClient(c *STCPProxyConfig) error {
	return nil
}
