package main

type ConfigProvider struct {
	Path    string
	Targets []Configurable
	Config  *Config
}

func NewConfigProvider(path string) (*ConfigProvider, error) {
	conf, err := ReadConfig(path)
	if err != nil {
		return nil, err
	}

	return &ConfigProvider{
		Path:    path,
		Targets: []Configurable{},
		Config:  conf,
	}, nil
}

func (p *ConfigProvider) AddTarget(target Configurable) {
	p.Targets = append(p.Targets, target)
}

func (p *ConfigProvider) Reload() error {
	conf, err := ReadConfig(p.Path)
	if err != nil {
		return err
	}

	p.Config = conf
	return p.Reconfigure(conf)
}

func (p *ConfigProvider) Reconfigure(conf *Config) error {
	for _, target := range p.Targets {
		if err := target.Configure(conf); err != nil {
			return err
		}
	}

	return nil
}

type Configurable interface {
	Configure(config *Config) error
}
