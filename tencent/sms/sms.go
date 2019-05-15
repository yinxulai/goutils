package sms

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	httpd "github.com/yinxulai/grpc-services/toolbox/http"
	"github.com/yinxulai/grpc-services/toolbox/random"
)

const (
	sendURL = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s"
)

// SMS 短信
type SMS struct {
	AppKey string
}

// SendParams 发送参数
type SendParams struct {
	Ext    string        `json:"ext"`
	Sig    string        `json:"sig"`
	Sign   string        `json:"sign"`
	Time   uint64        `json:"time"`
	TplID  uint64        `json:"tplID"`
	Extend string        `json:"extend"`
	Params []interface{} `json:"params"`
	Tel    struct {
		Mobile     string `json:"mobile"`
		Nationcode string `json:"nationcode"`
	} `json:"tel"`
}

// SendResponse 发送短信效应
type SendResponse struct {
	Sid    string `json:"sid"`
	Fee    uint64 `json:"fee"`
	Ext    string `json:"ext"`
	Result uint64 `json:"result"`
	Errmsg string `json:"errmsg"`
}

// Send 发送短信
func (srv *SMS) Send(params *SendParams) (response *SendResponse, err error) {
	response = new(SendResponse)
	urlStr := fmt.Sprintf(sendURL, srv.AppKey, random.Number(100, 10000))

	result, err := httpd.PostJSON(urlStr, params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Sig 签名
func (srv *SMS) Sig(params *SendParams) string {
	h := sha256.New()
	mobile := params.Tel.Mobile
	randomNumber := random.Number(10000, 99999)
	h.Write([]byte(fmt.Sprintf("appkey=%s&random=%s&time=%s&mobile=%s", srv.AppKey, randomNumber, params.Time, mobile)))
	return string(h.Sum(nil))
}
