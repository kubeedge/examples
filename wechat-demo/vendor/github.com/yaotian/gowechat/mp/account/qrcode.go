package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/yaotian/gowechat/mp/base"
	"github.com/yaotian/gowechat/wxcontext"
)

const (
	qrcodeURL      = "https://api.weixin.qq.com/cgi-bin/qrcode/create"
	ticketToImgURL = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s"
)

//Qrcode 带参数的二维码
type Qrcode struct {
	base.MpBase
}

//NewQrcode 实例化
func NewQrcode(context *wxcontext.Context) *Qrcode {
	qrcode := new(Qrcode)
	qrcode.Context = context
	return qrcode
}

const (
	//TemporaryQRCodeExpireSecondsLimit 临时二维码 expire_seconds 限制
	TemporaryQRCodeExpireSecondsLimit = 2592000
	//PermanentQRCodeSceneIDLimit 永久二维码 scene_id 限制
	PermanentQRCodeSceneIDLimit = 100000
)

//QrcodeResult Qrcode Result
type QrcodeResult struct {
	Ticket        string `json:"ticket"`                   // 获取的二维码ticket, 凭借此ticket可以在有效时间内换取二维码.
	URL           string `json:"url"`                      // 二维码图片解析后的地址, 开发者可根据该地址自行生成需要的二维码图片
	ExpireSeconds int    `json:"expire_seconds,omitempty"` // 二维码的有效时间, 以秒为单位. 最大不超过 604800.
}

//ImageURL ticket 换取二维码图片
func (c *QrcodeResult) ImageURL() (imgURL string) {
	return fmt.Sprintf(ticketToImgURL, url.QueryEscape(c.Ticket))
}

//CreateTemporaryQRCode  创建临时二维码
//  SceneId:       场景值ID, 为32位非0整型
//  ExpireSeconds: 二维码有效时间, 以秒为单位.  最大不超过 604800.
func (c *Qrcode) CreateTemporaryQRCode(SceneID uint32, ExpireSeconds int) (result *QrcodeResult, err error) {
	if SceneID == 0 {
		err = errors.New("SceneId should be greater than 0")
		return
	}
	if ExpireSeconds <= 0 {
		err = errors.New("ExpireSeconds should be greater than 0")
		return
	}
	var request struct {
		ExpireSeconds int    `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneID uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ExpireSeconds = ExpireSeconds
	request.ActionName = "QR_SCENE"
	request.ActionInfo.Scene.SceneID = SceneID

	result = new(QrcodeResult)

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, result)
	return
}

//CreateTemporaryQRCodeWithSceneString 创建临时二维码 scene_str
func (c *Qrcode) CreateTemporaryQRCodeWithSceneString(SceneString string, ExpireSeconds int) (result *QrcodeResult, err error) {
	if SceneString == "" {
		err = errors.New("SceneString should not be empty")
		return
	}
	if ExpireSeconds <= 0 {
		err = errors.New("ExpireSeconds should be greater than 0")
		return
	}
	var request struct {
		ExpireSeconds int    `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneString string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ExpireSeconds = ExpireSeconds
	request.ActionName = "QR_STR_SCENE"
	request.ActionInfo.Scene.SceneString = SceneString

	result = new(QrcodeResult)

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, result)
	return
}

//CreatePermanentQRCode 创建永久二维码
//  SceneId: 场景值ID, 目前参数只支持1--100000
func (c *Qrcode) CreatePermanentQRCode(sceneID uint32) (result *QrcodeResult, err error) {
	if sceneID == 0 {
		err = errors.New("SceneId should be greater than 0")
		return
	}
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneID uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneID = sceneID

	result = new(QrcodeResult)

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, result)
	return
}

//CreatePermanentQRCodeWithSceneString 创建永久二维码
//  SceneString: 场景值ID(字符串形式的ID), 字符串类型, 长度限制为1到64
func (c *Qrcode) CreatePermanentQRCodeWithSceneString(SceneString string) (result *QrcodeResult, err error) {
	if SceneString == "" {
		err = errors.New("SceneString should not be empty")
		return
	}
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneString string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_STR_SCENE"
	request.ActionInfo.Scene.SceneString = SceneString

	result = new(QrcodeResult)

	response, err := c.HTTPPostJSONWithAccessToken(qrcodeURL, &request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, result)
	return
}
