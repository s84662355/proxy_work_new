package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

var ConfData *confData = &confData{}

func Load(path string) {
	viper.SetConfigFile(path)
	// 会查找和读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("加载配置文件失败 config file: %s \n", err))
	}
	if err := viper.Unmarshal(ConfData); err != nil {
		panic(fmt.Errorf("解释配置文件失败 config file: %s \n", err))
	}

	serviceConfData := viper.GetString(ServiceName)
	if serviceConfData != "" {
		if err := json.Unmarshal([]byte(serviceConfData), ConfData.ServiceConf); err != nil {
			panic(fmt.Errorf("解释配置文件serviceConfData失败: %s \n", err))
		}
	}
}
