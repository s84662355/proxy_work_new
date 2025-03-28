package config

type mysql_config struct {
	Source             string `json:"source" mapstructure:"source" `
	MaxIdleConns       int    `json:"maxIdleConns" mapstructure:"maxIdleConns" `
	MaxOpenConns       int    `json:"maxOpenConns" mapstructure:"maxOpenConns" `
	SetConnMaxLifetime int    `json:"setConnMaxLifetime" mapstructure:"setConnMaxLifetime" ` // 秒
}
