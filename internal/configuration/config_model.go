package configuration

var RuntimeConf = RuntimeConfig{}

type RuntimeConfig struct {
	Scenario 		[]Scene      `yaml:"Scenario" mapstructure:"Scenario"`
}

type Scene struct {
	Name string `yaml:"Name" mapstructure:"name"`
	Url string `yaml:"Url" mapstructure:"url"`
	Type string `yaml:"Type" mapstructure:"type"`
}
