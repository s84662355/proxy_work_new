package service

import (
	"encoding/json"
	"fmt"

	"mproxy/utils/rabbitMQ"
)

// /推送黑名单
func PushDomainBackList(
	mqProduct *rabbitMQ.RabbitMQPool,
	mqExchange string,
	bl []string,
) error {
	body, err := json.Marshal(bl)
	if err != nil {
		return fmt.Errorf("推送黑名单json.Marshal error:%+v", err)
	}
	data := rabbitMQ.GetRabbitMqDataFormat(
		mqExchange,
		rabbitMQ.EXCHANGE_TYPE_FANOUT,
		"", "",
		body,
	)

	if err := mqProduct.Push(data); err != nil {
		return fmt.Errorf("推送黑名单mqProduct.Push error:%+v", err)
	}
	return nil
}
