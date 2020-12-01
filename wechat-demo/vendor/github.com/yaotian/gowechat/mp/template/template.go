package template

import (
	"encoding/json"

	"github.com/yaotian/gowechat/mp/base"
	"github.com/yaotian/gowechat/util"
	"github.com/yaotian/gowechat/wxcontext"
)

const (
	templateSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/send"

	templateAddURL         = "https://api.weixin.qq.com/cgi-bin/template/api_add_template"
	templateAllURL         = "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template"
	templateSetIndustryURL = "https://api.weixin.qq.com/cgi-bin/template/api_set_industry"
	templateGetIndustryURL = "https://api.weixin.qq.com/cgi-bin/template/get_industry"
)

//Template 模板消息
type Template struct {
	base.MpBase
}

//NewTemplate 实例化
func NewTemplate(context *wxcontext.Context) *Template {
	tpl := new(Template)
	tpl.Context = context
	return tpl
}

//Message 发送的模板消息内容
type Message struct {
	ToUser     string               `json:"touser"`          // 必须, 接受者OpenID
	TemplateID string               `json:"template_id"`     // 必须, 模版ID
	URL        string               `json:"url,omitempty"`   // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	Color      string               `json:"color,omitempty"` // 可选, 整个消息的颜色, 可以不设置
	Data       map[string]*DataItem `json:"data"`            // 必须, 模板数据

	MiniProgram struct {
		AppID    string `json:"appid"`    //所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
		PagePath string `json:"pagepath"` //所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar）
	} `json:"miniprogram"` //可选,跳转至小程序地址
}

//DataItem 模版内某个 .DATA 的值
type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type resTemplateSend struct {
	util.CommonError

	MsgID int64 `json:"msgid"`
}

//Send 发送模板消息
func (tpl *Template) Send(msg *Message) (msgID int64, err error) {
	response, err := tpl.HTTPPostJSONWithAccessToken(templateSendURL, msg)
	if err != nil {
		return 0, err
	}

	var result resTemplateSend
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	msgID = result.MsgID
	return
}

//IndustryList 行业列表
type IndustryList struct {
	PrimaryIndustry *Industry `json:"primary_industry"`
	SecondIndustry  *Industry `json:"secondary_industry"`
}

//Industry 行业
type Industry struct {
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
}

//Tmpl 模板
type Tmpl struct {
	TemplateId      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
}

//TmplList 模板列表
type TmplList struct {
	Templates []*Tmpl `json:"template_list"`
}

//AddTemplate 增加一个模板
func (tpl *Template) AddTemplate(templateIDShort string) (templateID string, err error) {
	type reqAddTmpl struct {
		TemplateIDShort string `json:"template_id_short"`
	}
	var response []byte
	response, err = tpl.HTTPPostJSONWithAccessToken(templateAddURL, reqAddTmpl{TemplateIDShort: templateIDShort})
	if err != nil {
		return "", err
	}

	var result Tmpl
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	templateID = result.TemplateId
	return
}

//GetTemplateList 查询模板列表
func (tpl *Template) GetTemplateList(templateIDShort string) (list TmplList, err error) {
	var response []byte
	response, err = tpl.HTTPGetWithAccessToken(templateAllURL)
	err = json.Unmarshal(response, &list)
	return
}

//GetTemplateIndustry 获得模板行业
func (tpl *Template) GetTemplateIndustry() (industryList IndustryList, err error) {
	var response []byte
	response, err = tpl.HTTPGetWithAccessToken(templateGetIndustryURL)
	err = json.Unmarshal(response, &industryList)
	return
}

//SetTemplateIndustry 设置模板行业
func (tpl *Template) SetTemplateIndustry(industry1, industry2 int) (err error) {
	type reqSetIndustry struct {
		Industry1 int `json:"industry_id1"`
		Industry2 int `json:"industry_id2"`
	}
	var req reqSetIndustry
	if industry1 > 0 {
		req.Industry1 = industry1
	}
	if industry2 > 0 {
		req.Industry2 = industry2
	}
	_, err = tpl.HTTPPostJSONWithAccessToken(templateSetIndustryURL, req)
	return
}
