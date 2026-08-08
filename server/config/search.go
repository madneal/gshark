package config

type Search struct {
	GobusterFilePath    string `mapstructure:"gobuster-filepath" json:"gobuster-filepath" yaml:"gobuster-filepath"`
	SubdomainWordList   string `mapstructure:"subdomain-wordlist" json:"subdomain-wordlist" yaml:"subdomain-wordlist"`
	SearchNum           int    `mapstructure:"searchnum" json:"searchnum" yaml:"searchnum"`
	GitlabDiscoverPages int    `mapstructure:"gitlab-discover-pages" json:"gitlab-discover-pages" yaml:"gitlab-discover-pages"`
	GitlabBatchSize     int    `mapstructure:"gitlab-batch-size" json:"gitlab-batch-size" yaml:"gitlab-batch-size"`
	SourcegraphURL      string `mapstructure:"sourcegraph-url" json:"sourcegraph-url" yaml:"sourcegraph-url"`
	SourcegraphToken    string `mapstructure:"sourcegraph-token" json:"sourcegraph-token" yaml:"sourcegraph-token"`
}
