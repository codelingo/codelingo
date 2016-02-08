package config

import "github.com/juju/errors"

type ServicesConfig struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

// serviceConfig contains configuration to connect to each service.
type ServiceConfig map[string]interface{}

func Services() (*ServicesConfig, error) {

	services := &ServicesConfig{}
	if err := Load(ServicesCfgFile, services); err != nil {
		return nil, err
	}

	return services, nil
}

func Service(serviceName string) (ServiceConfig, error) {
	if serviceName == "" {
		return nil, errors.New("service name cannot be empty")
	}

	cfg, err := Services()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if cfgBlock, ok := cfg.Services[serviceName]; ok {
		return cfgBlock, nil
	}

	return nil, errors.Errorf("configuration not found for %s in services config does not contain", serviceName)
}

var ServicesTmpl = `
services:
    github:

        # A Github API authentication token to allow Lingo to post reviews on
        # your behalf.
        # token: your-token

        # Domain of the service.
        # domain: http://github.com

    reviewboard:
        # Domain of the service.
        # domain: http://reviews.vapour.ws

        # A Reviewboard API authentication token to allow Lingo to post
        # reviews on your behalf.
        # token: your-token

        # whether or not to publish the review
        # publish: "true"
`[1:]
