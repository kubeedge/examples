package main

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	beegoContext "github.com/astaxie/beego/context"
	"github.com/yaotian/gowechat"
	"github.com/yaotian/gowechat/mp/message"
	"github.com/yaotian/gowechat/mp/user"
	"github.com/yaotian/gowechat/wxcontext"
	"k8s.io/client-go/rest"

	"github.com/kubeedge/examples/wechat-demo/kubeedge-wechat-app/utils"
	"github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha2"
)

// DeviceStatus is used to patch device status
type DeviceStatus struct {
	Status v1alpha2.DeviceStatus `json:"status"`
}

// The device id of the speaker
var deviceID = "speaker-01"

// The default namespace in which the speaker device instance resides
var namespace = "default"

// A regular expression which evaluates to true if the text contains `play`.
var rp = regexp.MustCompile(`\w*play\w*`)

// A regular expression which evaluates to true if the text contains `stop`.
var rs = regexp.MustCompile(`\w*stop\w*`)

// The CRD client used to patch the device instance.
var crdClient *rest.RESTClient

var appURL = "http://localhost:80"

// WeChat config
var config wxcontext.Config

func main() {
	// Create a client to talk to the K8S API server to patch the device CRDs
	kubeConfig, err := utils.KubeConfig()
	if err != nil {
		log.Fatalf("Failed to create KubeConfig, error : %v", err)
	}
	log.Println("Get kubeConfig successfully")

	crdClient, err = utils.NewCRDClient(kubeConfig)
	if err != nil {
		log.Fatalf("Failed to create device crd client , error : %v", err)
	}
	log.Println("Get crdClient successfully")

	// Set WeChat config from the secret mounted in the container.
	err = getWeChatConfigFromSecret()
	if err != nil {
		log.Fatalf("Failed to create device crd client , error : %v", err)
	}
	log.Println("Get WeChat Config from Secret successfully")

	// Listen
	beego.Any("/", textHandler)
	beego.Any("/oauth", wxOAuth) //需要网页授权的页面url  /oauth?target=url
	beego.Run(":80")
}

// getWeChatConfigFromSecret reads the credentials to authorize the application with WeChat.
func getWeChatConfigFromSecret() error {
	var err error
	appID, err := utils.ReadSecretKey(utils.AppID)
	if err != nil {
		log.Fatalf("Error reading %s : %v", utils.AppID, err)
		return err
	}
	appSecret, err := utils.ReadSecretKey(utils.AppSecret)
	if err != nil {
		log.Fatalf("Error reading %s : %v", utils.AppSecret, err)
		return err
	}
	token, err := utils.ReadSecretKey(utils.Token)
	if err != nil {
		log.Fatalf("Error reading %s : %v", utils.Token, err)
		return err
	}
	encodingAESKey, err := utils.ReadSecretKey(utils.EncodingAESKey)
	if err != nil {
		log.Fatalf("Error reading %s : %v", utils.EncodingAESKey, err)
		return err
	}

	log.Printf("AppID: %s\n", appID)
	log.Printf("AppSecret: %s\n", appSecret)
	log.Printf("Token: %s\n", token)
	log.Printf("EncodingAESKey: %s\n", encodingAESKey)

	// Contruct config
	config = wxcontext.Config{
		AppID:          appID,
		AppSecret:      appSecret,
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}

	return nil
}

//wxOAuth 微信公众平台，网页授权
func wxOAuth(ctx *beegoContext.Context) {
	var wechat = gowechat.NewWechat(config)
	mp, err := wechat.MpMgr()
	if err != nil {
		return
	}

	oauthHandler := mp.GetPageOAuthHandler(ctx.Request, ctx.ResponseWriter, appURL+"/oauth")

	oauthHandler.SetFuncCheckOpenIDExisting(func(openID string) (existing bool, stopNow bool) {
		//看自己的系统中是否已经存在此openID的用户
		//如果已经存在， 调用自己的Login 方法，设置cookie等，return true
		//如果还不存在，return false, handler会自动去取用户信息
		return false, true
	})

	oauthHandler.SetFuncAfterGetUserInfo(func(user user.Info) bool {
		//已获得用户信息，这里用信息做注册使用
		//调用自己的Login方法，设置cookie等
		return false
	})

	oauthHandler.Handle()
}

// textHandler
func textHandler(ctx *beegoContext.Context) {
	// 微信平台mp
	var wechat = gowechat.NewWechat(config)
	mp, err := wechat.MpMgr()
	log.Printf("Get MpMgr: %v", mp)
	if err != nil {
		log.Printf("Failed to get MpMgr: %v", err)
		return
	}

	// 传入request和responseWriter
	msgHandler := mp.GetMsgHandler(ctx.Request, ctx.ResponseWriter)
	log.Printf("Get msgHandler: %v", msgHandler)

	// 设置接收消息的处理方法
	msgHandler.SetHandleMessageFunc(func(msg message.MixMessage) *message.Reply {
		log.Printf("Got text %s\n", msg.Content)
		requestText := strings.ToLower(strings.TrimSpace(msg.Content))

		replyText := ""
		if requestText == "kubeedge" {
			replyText += "今日音乐推荐：\n"
			replyText += "1. 至少还有你\n"
			replyText += "2. 剩下的盛夏\n"
			replyText += "3. 普通朋友\n"
			replyText += "4. 起风啦\n"
			replyText += "5. 多余的解释\n"
			replyText += "回复play+数字开始播放音乐\n"
			replyText += "例如：play 1\n"
			replyText += "回复stop停止播放音乐\n"
		} else if rp.MatchString(requestText) {
			items := rp.Split(requestText, -1)
			songTrack := strings.TrimSpace(items[1])
			log.Printf("Music Track: %s\n", songTrack)
			// Update device twin
			UpdateDeviceTwinWithDesiredTrack(songTrack)
			replyText += "正在开始播放音乐，请稍等\n"
		} else if rs.MatchString(requestText) {
			songTrack := "stop"
			// Update device twin
			UpdateDeviceTwinWithDesiredTrack(songTrack)
			replyText += "正在停止播放音乐，请稍等\n"
		} else {
			log.Println("Could not parse the message")
			replyText += "不能识别您的输入\n"
		}
		return &message.Reply{message.MsgTypeText, message.NewText(replyText)}
	})

	// 处理消息接收以及回复
	err = msgHandler.Handle()
	if err != nil {
		log.Printf("Failed to handle message: %v", err)
	}
}

// UpdateDeviceTwinWithDesiredTrack patches the desired state of
// the device twin with the track to play.
func UpdateDeviceTwinWithDesiredTrack(track string) bool {
	status := buildStatusWithDesiredTrack(track)
	deviceStatus := &DeviceStatus{Status: status}
	body, err := json.Marshal(deviceStatus)
	if err != nil {
		log.Printf("Failed to marshal device status %v", deviceStatus)
		return false
	}
	result := crdClient.Patch(utils.MergePatchType).Namespace(namespace).Resource(utils.ResourceTypeDevices).Name(deviceID).Body(body).Do(context.TODO())
	if result.Error() != nil {
		log.Printf("Failed to patch device status %v of device %v in namespace %v \n error:%+v", deviceStatus, deviceID, namespace, result.Error())
		return false
	} else {
		log.Printf("Track [ %s ] will be played on speaker %s", track, deviceID)
	}
	return true
}

func buildStatusWithDesiredTrack(song string) v1alpha2.DeviceStatus {
	metadata := map[string]string{"timestamp": strconv.FormatInt(time.Now().Unix()/1e6, 10),
		"type": "string",
	}
	twins := []v1alpha2.Twin{{PropertyName: "track", Desired: v1alpha2.TwinProperty{Value: song, Metadata: metadata}, Reported: v1alpha2.TwinProperty{Value: "unknown", Metadata: metadata}}}
	devicestatus := v1alpha2.DeviceStatus{Twins: twins}
	return devicestatus
}
