package config

const (
	ConfigEnv  = "GSHARK_CONFIG"
	ConfigFile = "config.yaml"
)

type Server struct {
	Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Local   Local   `mapstructure:"local" json:"local" yaml:"local"`
	Search  Search  `mapstructure:"search" json:"search" yaml:"search"`
	Wechat  Wechat  `mapstructure:"wechat" json:"wechat" yaml:"wechat"`
}

type Local struct {
	Path string `mapstructure:"path" json:"path" yaml:"path" `
}

type Wechat struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
}
