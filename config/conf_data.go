package config

type confData struct {
	LogDir      string           `json:"log_dir" mapstructure:"log_dir" `
	Redis       *redis_config    `json:"redis" mapstructure:"redis"`
	Mysql       *mysql_config    `json:"mysql" mapstructure:"mysql"`
	ApiBaseUrl  string           `json:"api_base_url" mapstructure:"api_base_url"`
	FlowIncRate float64          `json:"flowIncRate" mapstructure:"flowIncRate"`
	Rabbitmq    *rabbitmq_config `json:"rabbitmq" mapstructure:"rabbitmq"`
	ServiceConf serviceConf      `json:"-"`
}
