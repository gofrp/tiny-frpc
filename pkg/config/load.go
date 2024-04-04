package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"

	v1 "github.com/gofrp/tiny-frpc/pkg/config/v1"
	"github.com/gofrp/tiny-frpc/pkg/util"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

var glbEnvs map[string]string

func init() {
	glbEnvs = make(map[string]string)
	envs := os.Environ()
	for _, env := range envs {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		glbEnvs[pair[0]] = pair[1]
	}
}

type Values struct {
	Envs map[string]string // environment vars
}

func GetValues() *Values {
	return &Values{
		Envs: glbEnvs,
	}
}

func RenderWithTemplate(in []byte, values *Values) ([]byte, error) {
	tmpl, err := template.New("frp").Parse(string(in))
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBufferString("")
	if err := tmpl.Execute(buffer, values); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func LoadFileContentWithTemplate(path string, values *Values) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return RenderWithTemplate(b, values)
}

func LoadConfigureFromFile(path string, c any, strict bool) error {
	content, err := LoadFileContentWithTemplate(path, GetValues())
	if err != nil {
		return err
	}
	return LoadConfigure(content, c, strict)
}

func LoadConfigure(b []byte, c any, strict bool) error {
	v1.DisallowUnknownFieldsMu.Lock()
	defer v1.DisallowUnknownFieldsMu.Unlock()
	v1.DisallowUnknownFields = strict

	log.Infof("file content: %v", string(b))

	var tomlObj interface{}
	// Try to unmarshal as TOML first; swallow errors from that (assume it's not valid TOML).
	if err := toml.Unmarshal(b, &tomlObj); err == nil {
		log.Infof("use toml unmarshal file to struct: toml obj: %v", util.JSONEncode(tomlObj))

		b, err = json.Marshal(&tomlObj)
		if err != nil {
			return err
		}
	} else {
		panic(err)
	}

	// If the buffer smells like JSON (first non-whitespace character is '{'), unmarshal as JSON directly.
	if util.IsJSONBuffer(b) {
		decoder := json.NewDecoder(bytes.NewBuffer(b))
		if strict {
			decoder.DisallowUnknownFields()
		}

		err := decoder.Decode(c)

		log.Infof("json decode config: %v", util.JSONEncode(c))

		return err
	}
	return nil
}

// LoadConfigure loads configuration from bytes and unmarshal into c.
// Now it supports json, yaml and toml format.
func LoadClientConfig(path string, strict bool) (
	*v1.ClientCommonConfig,
	[]v1.ProxyConfigurer,
	[]v1.VisitorConfigurer,
	bool, error,
) {
	var (
		cliCfg         *v1.ClientCommonConfig
		proxyCfgs      = make([]v1.ProxyConfigurer, 0)
		visitorCfgs    = make([]v1.VisitorConfigurer, 0)
		isLegacyFormat bool
	)

	allCfg := v1.ClientConfig{}
	if err := LoadConfigureFromFile(path, &allCfg, strict); err != nil {
		return nil, nil, nil, false, err
	}

	log.Infof("get all cfg: %v", util.JSONEncode(allCfg))

	cliCfg = &allCfg.ClientCommonConfig
	for _, c := range allCfg.Proxies {
		proxyCfgs = append(proxyCfgs, c.ProxyConfigurer)
	}
	for _, c := range allCfg.Visitors {
		visitorCfgs = append(visitorCfgs, c.VisitorConfigurer)
	}

	// Load additional config from includes.
	// legacy ini format already handle this in ParseClientConfig.
	if len(cliCfg.IncludeConfigFiles) > 0 && !isLegacyFormat {
		extProxyCfgs, extVisitorCfgs, err := LoadAdditionalClientConfigs(cliCfg.IncludeConfigFiles, isLegacyFormat, strict)
		if err != nil {
			return nil, nil, nil, isLegacyFormat, err
		}
		proxyCfgs = append(proxyCfgs, extProxyCfgs...)
		visitorCfgs = append(visitorCfgs, extVisitorCfgs...)
	}

	if cliCfg != nil {
		cliCfg.Complete()
	}
	for _, c := range proxyCfgs {
		c.Complete(cliCfg.User)
	}
	for _, c := range visitorCfgs {
		c.Complete(cliCfg)
	}
	return cliCfg, proxyCfgs, visitorCfgs, isLegacyFormat, nil
}

func LoadAdditionalClientConfigs(paths []string, isLegacyFormat bool, strict bool) ([]v1.ProxyConfigurer, []v1.VisitorConfigurer, error) {
	proxyCfgs := make([]v1.ProxyConfigurer, 0)
	visitorCfgs := make([]v1.VisitorConfigurer, 0)
	for _, path := range paths {
		absDir, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			return nil, nil, err
		}
		if _, err := os.Stat(absDir); os.IsNotExist(err) {
			return nil, nil, err
		}
		files, err := os.ReadDir(absDir)
		if err != nil {
			return nil, nil, err
		}
		for _, fi := range files {
			if fi.IsDir() {
				continue
			}
			absFile := filepath.Join(absDir, fi.Name())
			if matched, _ := filepath.Match(filepath.Join(absDir, filepath.Base(path)), absFile); matched {
				// support yaml/json/toml
				cfg := v1.ClientConfig{}
				if err := LoadConfigureFromFile(absFile, &cfg, strict); err != nil {
					return nil, nil, fmt.Errorf("load additional config from %s error: %v", absFile, err)
				}
				for _, c := range cfg.Proxies {
					proxyCfgs = append(proxyCfgs, c.ProxyConfigurer)
				}
				for _, c := range cfg.Visitors {
					visitorCfgs = append(visitorCfgs, c.VisitorConfigurer)
				}
			}
		}
	}
	return proxyCfgs, visitorCfgs, nil
}
