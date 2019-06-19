package bridge

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/yaotian/gowechat/mp/oauth"
	"github.com/yaotian/gowechat/mp/user"
	"github.com/yaotian/gowechat/wxcontext"
)

//PageOAuthHandler 微信网页授权
type PageOAuthHandler struct {
	// *wxcontext.Context
	*oauth.Oauth

	myURLOfPageOAuthCallback string
	myURLNeedPageOAuth       string

	openID                  string
	openIDExisting          bool
	checkOpenIDExistingFunc func(openID string) (existing bool, stopNow bool)

	user.Info
	afterGetUserInfoFunc func(user user.Info) (stopNow bool)
}

//NewPageOAuthHandler PageOAuthHandler初始化
func NewPageOAuthHandler(context *wxcontext.Context, myURLOfPageOAuthCallback string) *PageOAuthHandler {
	pa := new(PageOAuthHandler)
	pa.Oauth = oauth.NewOauth(context)
	pa.myURLOfPageOAuthCallback = myURLOfPageOAuthCallback
	return pa
}

func (c *PageOAuthHandler) getCallbackURL() (u string) {
	return fmt.Sprintf("%s?target=%s", c.myURLOfPageOAuthCallback, url.QueryEscape(c.myURLNeedPageOAuth))
}

/*SetFuncCheckOpenIDExisting 设置检查OpenID在您的系统中是否已经存在

handler:

	func(openID string) (existing bool){
		//用获得的openID，检查是否在你的系统中已经存在此用户
		//如果存在，调用你的Login方法，设置cookie, session等，然后return true

		//如果你的系统中不存在此openID用户, return false, handler会自动去获取用户信息
	}

*/
func (c *PageOAuthHandler) SetFuncCheckOpenIDExisting(handler func(string) (existing bool, stopNow bool)) {
	c.checkOpenIDExistingFunc = handler
}

/*SetFuncAfterGetUserInfo 设置获得用户信息后执行

handler:

	func(user oauth.UserInfo) (needStop bool) {
		//handler已经获得了用户信息，你可以用此信息，自动为用户完成一些动作，比如注册，头像等

		//默认needStop为false, 表示handler会自动redirect到你最开始需要授权的网页，此时你的系统已经完成了自动登陆等动作
		//如果你需要redirect到你需要的url，直接调用http.redirect； return true。 表示需要停止后面的动作
	}


*/
func (c *PageOAuthHandler) SetFuncAfterGetUserInfo(handler func(user.Info) bool) {
	c.afterGetUserInfoFunc = handler
}

//Handle handler
func (c *PageOAuthHandler) Handle() (err error) {
	code := c.Query("code")
	c.myURLNeedPageOAuth = c.Query("target")
	if code != "" {
		var acsTkn oauth.ResAccessToken
		acsTkn, err = c.GetUserAccessToken(code)
		if err != nil {
			return
		}
		openID := acsTkn.OpenID
		existing, stopNow := c.checkOpenIDExistingFunc(openID)
		if stopNow {
			return
		}
		if existing {
			http.Redirect(c.Writer, c.Request, c.myURLNeedPageOAuth, 302)
			return
		}
		//用 user模块的，没用oauth模板，可以获得更多信息
		u, err := user.NewUser(c.Oauth.Context).GetUserInfo(openID)
		if err != nil {
			return err
		}
		stopNow = c.afterGetUserInfoFunc(*u)
		if stopNow {
			return nil
		}
		http.Redirect(c.Writer, c.Request, c.myURLNeedPageOAuth, 302)
		return nil
	} else {
		//code为空时
		c.Redirect(c.getCallbackURL(), "snsapi_base", "base")
	}
	return
}
