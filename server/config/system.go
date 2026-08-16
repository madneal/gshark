package config

type System struct {
	Env                  string `mapstructure:"env" json:"env" yaml:"env"`
	Addr                 int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType               string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	UseMultipoint        bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	GitlabBase           string `mapstructure:"gitlab-base" json:"gitlabBase" yaml:"gitlab-base"`
	AiServer             string `mapstructure:"ai_server" json:"aiServer" yaml:"ai_server"`
	AiToken              string `mapstructure:"ai_token" json:"aiToken" yaml:"ai_token"`
	Model                string `mapstructure:"model" json:"model" yaml:"model"`
	AiAnalysisEnabled    bool   `mapstructure:"ai_analysis_enabled" json:"aiAnalysisEnabled" yaml:"ai_analysis_enabled"`
	AiAnalysisTimeout    int    `mapstructure:"ai_analysis_timeout" json:"aiAnalysisTimeout" yaml:"ai_analysis_timeout"`
	AiAnalysisMaxContent int    `mapstructure:"ai_analysis_max_content" json:"aiAnalysisMaxContent" yaml:"ai_analysis_max_content"`
}
