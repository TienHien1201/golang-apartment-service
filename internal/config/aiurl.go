package config

type AiConfig struct {
	URL         string
	DownloadURL string
	StaffURL    string
	StaffToken  string
}

func (c *Config) InitAiURLConfig() *AiConfig {
	return &AiConfig{
		URL:         c.Ai.URL,
		DownloadURL: c.Ai.DownloadURL,
		StaffURL:    c.Ai.StaffURL,
		StaffToken:  c.Ai.StaffToken,
	}
}
