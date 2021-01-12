package base

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yaotian/gowechat/util"
	"github.com/yaotian/gowechat/wxcontext"
)

//MpBase 微信公众平台,基本类
type MpBase struct {
	*wxcontext.Context
}

//HTTPGetWithAccessToken 微信公众平台中，自动加上access_token变量的GET调用，
//如果失败，会清空AccessToken cache, 再试一次
func (c *MpBase) HTTPGetWithAccessToken(url string) (resp []byte, err error) {
	retry := 1
Do:
	var accessToken string
	accessToken, err = c.GetAccessToken()
	if err != nil {
		return
	}

	var target = ""
	if strings.Contains(url, "?") {
		target = fmt.Sprintf("%s&access_token=%s", url, accessToken)
	} else {
		target = fmt.Sprintf("%s?access_token=%s", url, accessToken)
	}

	var response *http.Response
	response, err = http.Get(target)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = ioutil.ReadAll(response.Body)
	err = util.CheckCommonError(resp)
	if err == util.ErrUnmarshall {
		return
	}
	if err != nil {
		if retry > 0 {
			retry--
			c.CleanAccessTokenCache()
			goto Do
		}
		return
	}
	return
}

//HTTPPostJSONWithAccessToken post json 自动加上access token, 并retry
func (c *MpBase) HTTPPostJSONWithAccessToken(url string, obj interface{}) (resp []byte, err error) {
	retry := 1
Do:
	var accessToken string
	accessToken, err = c.GetAccessToken()
	if err != nil {
		return
	}

	var target = ""
	if strings.Contains(url, "?") {
		target = fmt.Sprintf("%s&access_token=%s", url, accessToken)
	} else {
		target = fmt.Sprintf("%s?access_token=%s", url, accessToken)
	}

	resp, err = util.PostJSON(target, obj)

	err = util.CheckCommonError(resp)
	if err == util.ErrUnmarshall {
		return
	}
	if err != nil {
		if retry > 0 {
			retry--
			c.CleanAccessTokenCache()
			goto Do
		}
		return
	}
	return
}
