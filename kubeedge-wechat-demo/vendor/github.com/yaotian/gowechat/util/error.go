package util

import (
	"encoding/json"
	"errors"
	"fmt"
)

//ErrUnmarshall err when unmarshall
var ErrUnmarshall = errors.New("Json Unmarshal Error")

//CommonError 微信返回的错误信息
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

//NewCommonError new CommonError
func NewCommonError(code int64, msg string) *CommonError {
	return &CommonError{ErrCode: code, ErrMsg: msg}
}

func (e *CommonError) Error() string {
	return e.ErrMsg
}

//CheckCommonError check CommonError
func CheckCommonError(jsonData []byte) error {
	var errmsg CommonError
	if err := json.Unmarshal(jsonData, &errmsg); err != nil {
		return ErrUnmarshall
	}

	if errmsg.ErrCode != 0 {
		return fmt.Errorf("Error , errcode=%d , errmsg=%s", errmsg.ErrCode, errmsg.ErrMsg)
	}

	return nil
}
