package config

type confData struct {
	LogDir string `json:"log_dir"`
	Redis  struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	ServiceConf *serviceConf `json:"-"`
}
