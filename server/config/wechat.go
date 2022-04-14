package config

type Wechat struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
}
