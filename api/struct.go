package api

import (
	_ "encoding/json"
)

type JSONData[T any] struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data T      `json:"data"`
}
