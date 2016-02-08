package config

type defaultsYamlCfg struct {
	Defaults *defaultsConfig `yaml:"defaults"`
}

type defaultsConfig struct {
	Tenet tenetDefaults `yaml:"tenet"`
}

type tenetDefaults struct {
	Driver string `yaml:"driver"`
}

func Defaults() (*defaultsConfig, error) {
	cfg := &defaultsYamlCfg{}
	if err := Load(DefaultsCfgFile, cfg); err != nil {
		return nil, err
	}

	return cfg.Defaults, nil
}

var DefaultsTmpl = `
defaults:
    tenet:

        # when adding a new tenet, the tenet driver will be set to this default. Valid values are: "binary" or "docker".
        driver: binary
`[1:]
