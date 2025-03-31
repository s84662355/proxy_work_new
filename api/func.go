package api

import (
	"fmt"
	"net/http"

	"mproxy/config"
)

// /https://lekecloud.yuque.com/dk8avh/pb2sdo/xmoyhvco0o5glytg
func AutoBuyDayNoCurrentLimiting(
	userId interface{},
) (*JSONData[interface{}], error) {
	url := config.GetConf().ApiBaseUrl + "/transit/dynamic/userAutoBuyDay"
	data := make(map[string]interface{})
	data["userId"] = userId

	var err error
	var httpRequest *http.Request
	httpRequest, err = postRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("AutoBuyDayNoCurrentLimiting postRequest err:%+v", err)
	}

	var jsonData *JSONData[interface{}]

	jsonData, err = getResponse[interface{}](httpRequest, jsonData)
	if err != nil {
		return nil, fmt.Errorf("AutoBuyDayNoCurrentLimiting getResponse err:%+v", err)
	}

	return jsonData, nil
}

func AutomaticRechargeTraffic(userId interface{}, normal, datacenter bool) (*JSONData[interface{}], error) {
	url := config.GetConf().ApiBaseUrl + "/transit/dynamic/automaticRechargeTraffic"
	data := make(map[string]interface{})
	data["userId"] = userId
	data["normal"] = normal
	data["datacenter"] = datacenter

	var err error
	var httpRequest *http.Request
	httpRequest, err = postRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("AutomaticRechargeTraffic postRequest err:%+v", err)
	}

	var jsonData *JSONData[interface{}]

	jsonData, err = getResponse[interface{}](httpRequest, jsonData)
	if err != nil {
		return nil, fmt.Errorf("AutomaticRechargeTraffic getResponse err:%+v", err)
	}

	return jsonData, nil
}

// /transitDynamicAccount/syncNewAccount
func SyncNewAccount(username, password interface{}) (*JSONData[interface{}], error) {
	url := config.GetConf().ApiBaseUrl + "/transitDynamicAccount/syncSonAccount"
	data := make(map[string]interface{})
	data["username"] = username
	data["password"] = password

	var err error
	var httpRequest *http.Request
	httpRequest, err = postRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("SyncNewAccount postRequest err:%+v", err)
	}

	var jsonData *JSONData[interface{}]

	jsonData, err = getResponse[interface{}](httpRequest, jsonData)
	if err != nil {
		return nil, fmt.Errorf("SyncNewAccount postgetResponseRequest err:%+v", err)
	}
	return jsonData, nil
}
