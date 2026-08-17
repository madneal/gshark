package config

type AIProvider struct {
	Name   string `mapstructure:"name" json:"name" yaml:"name"`
	Server string `mapstructure:"server" json:"server" yaml:"server"`
	Token  string `mapstructure:"token" json:"token" yaml:"token"`
	Model  string `mapstructure:"model" json:"model" yaml:"model"`
}

type System struct {
	Env               string       `mapstructure:"env" json:"env" yaml:"env"`
	Addr              int          `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType            string       `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	UseMultipoint     bool         `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	GitlabBase        string       `mapstructure:"gitlab-base" json:"gitlabBase" yaml:"gitlab-base"`
	AiServer          string       `mapstructure:"ai_server" json:"aiServer" yaml:"ai_server"`
	AiToken           string       `mapstructure:"ai_token" json:"aiToken" yaml:"ai_token"`
	Model             string       `mapstructure:"model" json:"model" yaml:"model"`
	AiAnalysisEnabled bool         `mapstructure:"ai_analysis_enabled" json:"aiAnalysisEnabled" yaml:"ai_analysis_enabled"`
	AiProviders       []AIProvider `mapstructure:"ai_providers" json:"aiProviders" yaml:"ai_providers"`
}
