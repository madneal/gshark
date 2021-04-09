package config

type Search struct {
	GobusterFilePath  string `mapstructure:"gobuster-filepath" json:"gobuster-filepath" yaml:"gobuster-filepath"`
	SubdomainWordList string `mapstructure:"subdomain-wordlist" json:"subdomain-wordlist" yaml:"subdomain-wordlist"`
	SearchNum         int    `mapstructure:"searchnum" json:"searchnum" yaml:"searchnum"`
}
