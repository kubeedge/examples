package user

import (
	"encoding/json"
	"fmt"

	"github.com/yaotian/gowechat/mp/base"
	"github.com/yaotian/gowechat/util"
	"github.com/yaotian/gowechat/wxcontext"
)

const (
	userInfoURL = "https://api.weixin.qq.com/cgi-bin/user/info"
)

//User 用户管理
type User struct {
	base.MpBase
}

//NewUser 实例化
func NewUser(context *wxcontext.Context) *User {
	user := new(User)
	user.Context = context
	return user
}

//Info 用户基本信息
type Info struct {
	util.CommonError

	Subscribe     int32    `json:"subscribe"`
	OpenID        string   `json:"openid"`
	Nickname      string   `json:"nickname"`
	Sex           int      `json:"sex"`
	City          string   `json:"city"`
	Country       string   `json:"country"`
	Province      string   `json:"province"`
	Language      string   `json:"language"`
	Headimgurl    string   `json:"headimgurl"`
	SubscribeTime int64    `json:"subscribe_time"`
	UnionID       string   `json:"unionid"`
	Remark        string   `json:"remark"`
	GroupID       int32    `json:"groupid"`
	TagidList     []string `json:"tagid_list"`
}

//GetUserInfo 获取用户基本信息
func (user *User) GetUserInfo(openID string) (userInfo *Info, err error) {
	url := fmt.Sprintf("%s?openid=%s&lang=zh_CN", userInfoURL, openID)
	var response []byte
	response, err = user.HTTPGetWithAccessToken(url)
	if err != nil {
		return
	}
	userInfo = new(Info)
	err = json.Unmarshal(response, userInfo)
	return
}

//IsSubscribed 是否已经关注公众号
func (user *User) IsSubscribed(openID string) (subscribed bool, err error) {
	var userInfo *Info
	userInfo, err = user.GetUserInfo(openID)
	if err != nil {
		return
	}
	subscribed = userInfo.Subscribe == 1
	return
}
