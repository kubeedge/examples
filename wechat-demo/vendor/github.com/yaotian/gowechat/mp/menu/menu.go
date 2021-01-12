package menu

import (
	"encoding/json"

	"github.com/yaotian/gowechat/mp/base"
	"github.com/yaotian/gowechat/util"
	"github.com/yaotian/gowechat/wxcontext"
)

const (
	menuCreateURL            = "https://api.weixin.qq.com/cgi-bin/menu/create"
	menuGetURL               = "https://api.weixin.qq.com/cgi-bin/menu/get"
	menuDeleteURL            = "https://api.weixin.qq.com/cgi-bin/menu/delete"
	menuAddConditionalURL    = "https://api.weixin.qq.com/cgi-bin/menu/addconditional"
	menuDeleteConditionalURL = "https://api.weixin.qq.com/cgi-bin/menu/delconditional"
	menuTryMatchURL          = "https://api.weixin.qq.com/cgi-bin/menu/trymatch"
	menuSelfMenuInfoURL      = "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info"
)

//Menu struct
type Menu struct {
	base.MpBase
}

//reqMenu 设置菜单请求数据
type reqMenu struct {
	Button    []*Button  `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
}

//reqDeleteConditional 删除个性化菜单请求数据
type reqDeleteConditional struct {
	MenuID int64 `json:"menuid"`
}

//reqMenuTryMatch 菜单匹配请求
type reqMenuTryMatch struct {
	UserID string `json:"user_id"`
}

//resConditionalMenu 个性化菜单返回结果
type resConditionalMenu struct {
	Button    []Button  `json:"button"`
	MatchRule MatchRule `json:"matchrule"`
	MenuID    int64     `json:"menuid"`
}

//resMenuTryMatch 菜单匹配请求结果
type resMenuTryMatch struct {
	util.CommonError

	Button []Button `json:"button"`
}

//ResMenu 查询菜单的返回数据
type ResMenu struct {
	util.CommonError

	Menu struct {
		Button []Button `json:"button"`
		MenuID int64    `json:"menuid"`
	} `json:"menu"`
	Conditionalmenu []resConditionalMenu `json:"conditionalmenu"`
}

//ResSelfMenuInfo 自定义菜单配置返回结果
type ResSelfMenuInfo struct {
	util.CommonError

	IsMenuOpen   int32 `json:"is_menu_open"`
	SelfMenuInfo struct {
		Button []SelfMenuButton `json:"button"`
	} `json:"selfmenu_info"`
}

//SelfMenuButton 自定义菜单配置详情
type SelfMenuButton struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	URL       string `json:"url,omitempty"`
	Value     string `json:"value,omitempty"`
	SubButton struct {
		List []SelfMenuButton `json:"list"`
	} `json:"sub_button,omitempty"`
	NewsInfo struct {
		List []ButtonNew `json:"list"`
	} `json:"news_info,omitempty"`
}

//ButtonNew 图文消息菜单
type ButtonNew struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Digest     string `json:"digest"`
	ShowCover  int32  `json:"show_cover"`
	CoverURL   string `json:"cover_url"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url"`
}

//MatchRule 个性化菜单规则
type MatchRule struct {
	GroupID            int32  `json:"group_id,omitempty"`
	Sex                int32  `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType int32  `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

//NewMenu 实例
func NewMenu(context *wxcontext.Context) *Menu {
	menu := new(Menu)
	menu.Context = context
	return menu
}

//SetMenu 设置按钮
func (menu *Menu) SetMenu(buttons []*Button) error {
	reqMenu := &reqMenu{
		Button: buttons,
	}
	_, err := menu.HTTPPostJSONWithAccessToken(menuCreateURL, reqMenu)
	if err != nil {
		return err
	}
	return nil
}

//GetMenu 获取菜单配置
func (menu *Menu) GetMenu() (resMenu ResMenu, err error) {
	var response []byte
	response, err = menu.HTTPGetWithAccessToken(menuGetURL)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &resMenu)
	return
}

//DeleteMenu 删除菜单
func (menu *Menu) DeleteMenu() (err error) {
	_, err = menu.HTTPGetWithAccessToken(menuDeleteURL)
	return
}

//AddConditional 添加个性化菜单
func (menu *Menu) AddConditional(buttons []*Button, matchRule *MatchRule) error {
	reqMenu := &reqMenu{
		Button:    buttons,
		MatchRule: matchRule,
	}
	_, err := menu.HTTPPostJSONWithAccessToken(menuAddConditionalURL, reqMenu)
	if err != nil {
		return err
	}
	return nil
}

//DeleteConditional 删除个性化菜单
func (menu *Menu) DeleteConditional(menuID int64) error {
	reqDeleteConditional := &reqDeleteConditional{
		MenuID: menuID,
	}
	_, err := menu.HTTPPostJSONWithAccessToken(menuDeleteConditionalURL, reqDeleteConditional)
	if err != nil {
		return err
	}
	return nil
}

//MenuTryMatch 菜单匹配
func (menu *Menu) MenuTryMatch(userID string) (buttons []Button, err error) {
	reqMenuTryMatch := &reqMenuTryMatch{userID}
	var response []byte
	response, err = menu.HTTPPostJSONWithAccessToken(menuTryMatchURL, reqMenuTryMatch)
	if err != nil {
		return
	}
	var resMenuTryMatch resMenuTryMatch
	err = json.Unmarshal(response, &resMenuTryMatch)
	if err != nil {
		return
	}
	buttons = resMenuTryMatch.Button
	return
}

//GetCurrentSelfMenuInfo 获取自定义菜单配置接口
func (menu *Menu) GetCurrentSelfMenuInfo() (resSelfMenuInfo ResSelfMenuInfo, err error) {
	var response []byte
	response, err = menu.HTTPGetWithAccessToken(menuSelfMenuInfoURL)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &resSelfMenuInfo)
	return
}
