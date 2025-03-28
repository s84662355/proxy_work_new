package config

type confData struct {
	LogDir      string        `json:"log_dir" mapstructure:"log_dir" `
	Redis       *redis_config `json:"redis" mapstructure:"redis"`
	Mysql       *mysql_config `json:"mysql" mapstructure:"mysql"`
	ServiceConf serviceConf   `json:"-"`
}
