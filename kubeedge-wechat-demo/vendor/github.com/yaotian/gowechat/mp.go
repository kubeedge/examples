// @description  微信公共平台的接口

package gowechat

import (
	"net/http"

	"github.com/yaotian/gowechat/mp/account"
	"github.com/yaotian/gowechat/mp/bridge"
	"github.com/yaotian/gowechat/mp/jssdk"
	"github.com/yaotian/gowechat/mp/material"
	"github.com/yaotian/gowechat/mp/menu"
	"github.com/yaotian/gowechat/mp/oauth"
	"github.com/yaotian/gowechat/mp/template"
	"github.com/yaotian/gowechat/mp/user"
)

//MpMgr mp mgr
type MpMgr struct {
	*Wechat
}

//GetAccessToken 获取access_token
func (c *MpMgr) GetAccessToken() (string, error) {
	return c.Context.GetAccessToken()
}

// GetOauth oauth2网页授权
func (c *MpMgr) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(c.Context)
}

// GetMaterial 素材管理
func (c *MpMgr) GetMaterial() *material.Material {
	return material.NewMaterial(c.Context)
}

// GetJs js-sdk配置
func (c *MpMgr) GetJs() *jssdk.Js {
	return jssdk.NewJs(c.Context)
}

// GetMenu 菜单管理接口
func (c *MpMgr) GetMenu() *menu.Menu {
	return menu.NewMenu(c.Context)
}

// GetUser 用户管理接口
func (c *MpMgr) GetUser() *user.User {
	return user.NewUser(c.Context)
}

// GetTemplate 模板消息接口
func (c *MpMgr) GetTemplate() *template.Template {
	return template.NewTemplate(c.Context)
}

// GetMsgHandler 消息管理
func (c *MpMgr) GetMsgHandler(req *http.Request, writer http.ResponseWriter) *bridge.MsgHandler {
	c.Context.Request = req
	c.Context.Writer = writer
	return bridge.NewMsgHandler(c.Context)
}

//GetPageOAuthHandler 网页授权
func (c *MpMgr) GetPageOAuthHandler(req *http.Request, writer http.ResponseWriter, myURLOfPageOAuthCallback string) *bridge.PageOAuthHandler {
	c.Context.Request = req
	c.Context.Writer = writer
	handler := bridge.NewPageOAuthHandler(c.Context, myURLOfPageOAuthCallback)
	return handler
}

// GetQrcode 带参数的二维码
func (c *MpMgr) GetQrcode() *account.Qrcode {
	return account.NewQrcode(c.Context)
}
