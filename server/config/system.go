package config

type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	OssType       string `mapstructure:"oss-type" json:"ossType" yaml:"oss-type"`
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	GitlabBase    string `mapstructure:"gitlab-base" yaml:"gitlab-base"`
	AiServer      string `mapstructure:"ai_server" yaml:"ai_server"`
	AiToken       string `mapstructure:"ai_token" yaml:"ai_token"`
	Model         string `mapstructure:"model" json:"model" yaml:"model"`
}
