package config

type Server struct {
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	AutoCode Autocode `mapstructure:"autoCode" json:"autoCode" yaml:"autoCode"`
	// gorm
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	// oss
	Local  Local  `mapstructure:"local" json:"local" yaml:"local"`
	Excel  Excel  `mapstructure:"excel" json:"excel" yaml:"excel"`
	Search Search `mapstructure:"search" json:"search" yaml:"search"`
	Wechat Wechat `mapstructure:"wechat" json:"wechat" yaml:"wechat"`
}

type Local struct {
	Path string `mapstructure:"path" json:"path" yaml:"path" `
}

type Excel struct {
	Dir string `mapstructure:"dir" json:"dir" yaml:"dir"`
}

type Wechat struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
}
