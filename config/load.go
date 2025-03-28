package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

var ConfData *confData = &confData{}

func Load(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("json")
	// 会查找和读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("加载配置文件失败 config file: %s \n", err))
	}
	if err := viper.Unmarshal(ConfData); err != nil {
		panic(fmt.Errorf("解释配置文件失败 config file: %s \n", err))
	}
	fmt.Println(viper.GetString ( "log_dir"))
	serviceConfData := viper.GetStringMap(ServiceName)
	if len(serviceConfData) >0 {
		if b, err := json.Marshal(serviceConfData) ;err!=nil{
 		   panic(fmt.Errorf("解释配置文件失败   json Marshal config file: %s \n", err))
		}else if err := json.Unmarshal(b, &ConfData.ServiceConf); err != nil {
			panic(fmt.Errorf("解释配置文件serviceConfData失败: %s \n", err))
		}
		fmt.Println(ConfData.ServiceConf)
	}
}
