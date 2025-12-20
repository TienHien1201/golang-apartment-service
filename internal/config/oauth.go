package config

type AuthConfig struct {
	Google GoogleConfig `mapstructure:"google"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL  string `mapstructure:"callback_url"`
}
