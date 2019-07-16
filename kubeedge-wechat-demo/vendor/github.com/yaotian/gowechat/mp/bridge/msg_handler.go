package bridge

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/yaotian/gowechat/mp/message"
	"github.com/yaotian/gowechat/util"
	"github.com/yaotian/gowechat/wxcontext"
)

//MsgHandler struct
type MsgHandler struct {
	*wxcontext.Context

	handleMessageFunc func(message.MixMessage) *message.Reply

	requestRawXMLMsg  []byte
	requestMsg        message.MixMessage
	responseRawXMLMsg []byte
	responseMsg       interface{}

	isSafeMode bool
	random     []byte
	nonce      string
	timestamp  int64
}

//NewMsgHandler init
func NewMsgHandler(context *wxcontext.Context) *MsgHandler {
	srv := new(MsgHandler)
	fmt.Println("NewMsgHandler:", srv)
	srv.Context = context
	return srv
}

//Handle 处理微信的请求消息
func (srv *MsgHandler) Handle() error {
	//Request is GET
	//微信公众平台，设置服务器后保存，会调用此方法来做验证
	if strings.ToLower(srv.Context.Request.Method) == "get" {
		if !srv.Validate() {
			return fmt.Errorf("请求校验失败")
		}

		echostr, exists := srv.GetQuery("echostr")
		if exists {
			srv.String(echostr) //微信公众平台需要将此值发送回去，来完成验证
		}
		return nil
	}

	//Request is POST
	//微信公众平台将消息post到服务器上
	if strings.ToLower(srv.Context.Request.Method) == "post" {
		replyMsg, err := srv.handleRequest()
		if err != nil {
			return err
		}
		//debug
		// fmt.Println("request msg = ", string(srv.requestRawXMLMsg))
		err = srv.buildResponse(replyMsg)
		if err == nil {
			srv.Send()
		} else {
			return err
		}
	}
	return nil
}

//Validate 校验请求是否合法
func (srv *MsgHandler) Validate() bool {
	timestamp := srv.Query("timestamp")
	nonce := srv.Query("nonce")
	signature := srv.Query("signature")
	return signature == util.Signature(srv.Token, timestamp, nonce)
}

//HandleRequest 处理微信的请求
func (srv *MsgHandler) handleRequest() (reply *message.Reply, err error) {
	//set isSafeMode
	srv.isSafeMode = false
	encryptType := srv.Query("encrypt_type")
	if encryptType == "aes" {
		srv.isSafeMode = true
	}

	var msg interface{}
	msg, err = srv.getMessage()
	if err != nil {
		return
	}
	mixMessage, success := msg.(message.MixMessage)
	if !success {
		err = errors.New("消息类型转换失败")
	}
	srv.requestMsg = mixMessage
	reply = srv.handleMessageFunc(mixMessage)
	return
}

//getMessage 解析微信返回的消息
func (srv *MsgHandler) getMessage() (interface{}, error) {
	var rawXMLMsgBytes []byte
	var err error
	if srv.isSafeMode {
		var encryptedXMLMsg message.EncryptedXMLMsg
		if err := xml.NewDecoder(srv.Request.Body).Decode(&encryptedXMLMsg); err != nil {
			return nil, fmt.Errorf("从body中解析xml失败,err=%v", err)
		}

		//验证消息签名
		timestamp := srv.Query("timestamp")
		srv.timestamp, err = strconv.ParseInt(timestamp, 10, 32)
		if err != nil {
			return nil, err
		}
		nonce := srv.Query("nonce")
		srv.nonce = nonce
		msgSignature := srv.Query("msg_signature")
		msgSignatureGen := util.Signature(srv.Token, timestamp, nonce, encryptedXMLMsg.EncryptedMsg)
		if msgSignature != msgSignatureGen {
			return nil, fmt.Errorf("消息不合法，验证签名失败")
		}

		//解密
		srv.random, rawXMLMsgBytes, err = util.DecryptMsg(srv.AppID, encryptedXMLMsg.EncryptedMsg, srv.EncodingAESKey)
		if err != nil {
			return nil, fmt.Errorf("消息解密失败, err=%v", err)
		}
	} else {
		rawXMLMsgBytes, err = ioutil.ReadAll(srv.Request.Body)
		if err != nil {
			return nil, fmt.Errorf("从body中解析xml失败, err=%v", err)
		}
	}

	srv.requestRawXMLMsg = rawXMLMsgBytes

	return srv.parseRequestMessage(rawXMLMsgBytes)
}

func (srv *MsgHandler) parseRequestMessage(rawXMLMsgBytes []byte) (msg message.MixMessage, err error) {
	msg = message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, &msg)
	return
}

//SetHandleMessageFunc 设置用户自定义的回调方法
func (srv *MsgHandler) SetHandleMessageFunc(handler func(message.MixMessage) *message.Reply) {
	srv.handleMessageFunc = handler
}

func (srv *MsgHandler) buildResponse(reply *message.Reply) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: %v\n%s", e, debug.Stack())
		}
	}()
	if reply == nil {
		//do nothing
		return nil
	}
	msgType := reply.MsgType
	switch msgType {
	case message.MsgTypeText:
	case message.MsgTypeImage:
	case message.MsgTypeVoice:
	case message.MsgTypeVideo:
	case message.MsgTypeMusic:
	case message.MsgTypeNews:
	case message.MsgTypeTransfer:
	default:
		err = message.ErrUnsupportReply
		return
	}

	msgData := reply.MsgData
	value := reflect.ValueOf(msgData)
	//msgData must be a ptr
	kind := value.Kind().String()
	if "ptr" != kind {
		return message.ErrUnsupportReply
	}

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(srv.requestMsg.FromUserName)
	value.MethodByName("SetToUserName").Call(params)

	params[0] = reflect.ValueOf(srv.requestMsg.ToUserName)
	value.MethodByName("SetFromUserName").Call(params)

	params[0] = reflect.ValueOf(msgType)
	value.MethodByName("SetMsgType").Call(params)

	params[0] = reflect.ValueOf(util.GetCurrTs())
	value.MethodByName("SetCreateTime").Call(params)

	srv.responseMsg = msgData
	srv.responseRawXMLMsg, err = xml.Marshal(msgData)
	return
}

//Send 将自定义的消息发送
func (srv *MsgHandler) Send() (err error) {
	replyMsg := srv.responseMsg
	if srv.isSafeMode {
		//安全模式下对消息进行加密
		var encryptedMsg []byte
		encryptedMsg, err = util.EncryptMsg(srv.random, srv.responseRawXMLMsg, srv.AppID, srv.EncodingAESKey)
		if err != nil {
			return
		}
		//TODO 如果获取不到timestamp nonce 则自己生成
		timestamp := srv.timestamp
		timestampStr := strconv.FormatInt(timestamp, 10)
		msgSignature := util.Signature(srv.Token, timestampStr, srv.nonce, string(encryptedMsg))
		replyMsg = message.ResponseEncryptedXMLMsg{
			EncryptedMsg: string(encryptedMsg),
			MsgSignature: msgSignature,
			Timestamp:    timestamp,
			Nonce:        srv.nonce,
		}
	}
	if replyMsg != nil {
		srv.XML(replyMsg)
	}
	return
}
