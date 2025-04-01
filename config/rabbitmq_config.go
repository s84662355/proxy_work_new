package config

type rabbitmq_config struct {
	Host              string `json:"host" mapstructure:"host" `
	Port              int    `json:"port" mapstructure:"port" `
	User              string `json:"user" mapstructure:"user" `
	Password          string `json:"password" mapstructure:"password" `
	VirtualHost       string `json:"virtual_host" mapstructure:"virtual_host" `
	BlacklistExchange string `json:"black_list_exchange" mapstructure:"black_list_exchange" `
}
