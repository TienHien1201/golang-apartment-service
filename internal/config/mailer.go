package config

// MailerConfig holds all configuration for the outgoing mail service.
// Values are loaded from the YAML config file — never hardcoded in code.
type MailerConfig struct {
	SMTP     MailerSMTPConfig
	FromName string
}

// MailerSMTPConfig holds SMTP connection and authentication details.
type MailerSMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
}
