package service

import (
	"io/ioutil"
	"path/filepath"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/util"
	"gopkg.in/yaml.v2"
)

// This file contains the configuration structs for each service. They are
// populated by $LINGO_HOME/services.yaml

// servicesConfig maps to $LINGO_HOME/services.yaml
type servicesConfig struct {
	Services map[string]serviceConfig
}

// serviceConfig contains configuration to connect to the service
type serviceConfig map[string]interface{}

func Config(serviceName string) (serviceConfig, error) {

	if serviceName == "" {
		return nil, errors.New("service name cannot be empty")
	}

	lHome, err := util.LingoHome()
	if err != nil {
		return nil, errors.Trace(err)
	}

	cfgPath := filepath.Join(lHome, common.ConfigFile)
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, errors.Errorf("problem reading service config: %v", err)
	}

	var services servicesConfig
	if err := yaml.Unmarshal(data, &services); err != nil {
		return nil, errors.Errorf("problem reading service config: %v", err)
	}

	if cfgBlock, ok := services.Services[serviceName]; ok {
		return cfgBlock, nil
	}

	return nil, errors.Errorf("configuration not found for %s in services config does not contain", serviceName)
}
