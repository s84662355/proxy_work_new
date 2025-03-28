package config

type redis_config struct {
	Addr         string `json:"addr"  mapstructure:"addr"  `
	Password     string `json:"password"  mapstructure:"password"  `
	DB           int    `json:"db"  mapstructure:"db"  `
	MinIdleConns int    `json:"minIdleConns"  mapstructure:"minIdleConns"  ` // 最小空闲连接数
	MaxIdleConns int    `json:"maxIdleConns"  mapstructure:"maxIdleConns"  ` ///最大空闲连接数
	PoolSize     int    `json:"poolSize"  mapstructure:"poolSize"  `         // 连接池大小
}
