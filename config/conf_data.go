package config

type confData struct {
	LogDir string `json:"log_dir" mapstructure:"log_dir" `
	Redis  struct {
		Addr     string `json:"addr"  mapstructure:"addr"  `
		Password string `json:"password"  mapstructure:"password" `
		DB       int    `json:"db" mapstructure:"db" `
	} `json:"redis"`
	ServiceConf serviceConf `json:"-"`
}
