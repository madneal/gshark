package config

type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	OssType       string `mapstructure:"oss-type" json:"ossType" yaml:"oss-type"`
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	GitlabBase    string `json:"gitlabBase" yaml:"gitlab-base"`
	AiServer      string `json:"ai_server" yaml:"ai_server"`
	AiToken       string `json:"ai_token" yaml:"ai_token"`
}
