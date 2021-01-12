package paytool

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/yaotian/gowechat/mch/base"
	"github.com/yaotian/gowechat/util"
)

//官方文档： https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_4&index=3
var (
	//ErrNoEnoughMoney 商户平台上的余额不足，给用户发不了红包
	ErrNoEnoughMoney = errors.New("No enough money")
)

const (
	//SceneIDPromotion 商品促销
	SceneIDPromotion = "PRODUCT_1"

	//SceneIDLuckyDraw 抽奖
	SceneIDLuckyDraw = "PRODUCT_2"

	//SceneIDPrize 虚拟物品兑奖
	SceneIDPrize = "PRODUCT_3"

	//SceneIDBenefit 企业内部福利
	SceneIDBenefit = "PRODUCT_4"

	//SceneIDAgentBonous 渠道分润
	SceneIDAgentBonous = "PRODUCT_5"

	//SceneIDInsurance 保险回馈
	SceneIDInsurance = "PRODUCT_6"

	//SceneIDLottery 彩票派奖
	SceneIDLottery = "PRODUCT_7"

	//SceneIDTax 税务刮奖
	SceneIDTax = "PRODUCT_8"
)

//RedPackInput 发红包的配置
type RedPackInput struct {
	ToOpenID string //接红包的OpenID
	MoneyFen int    //分为单位

	SendName string //商户名称，String(32) 谁发的红包，一般为发红包的单位
	Wishing  string //红包祝福语 String(128) “感谢您参加猜灯谜活动，祝您元宵节快乐！”
	ActName  string //活动名称 String(32) 猜灯谜抢红包活动
	Remark   string //备注 String(256)

	IP string

	//非必填，但大于200元，此必填, 有8个选项可供选择
	SceneID string
}

//Check check input
func (m *RedPackInput) Check() (isGood bool, err error) {
	if m.ToOpenID == "" || m.MoneyFen == 0 || m.SendName == "" || m.Wishing == "" || m.ActName == "" || m.Remark == "" || m.IP == "" {
		err = fmt.Errorf("%s", "Input有必填项没有值")
		return
	}

	if m.MoneyFen >= 200*100 && m.SceneID == "" {
		err = fmt.Errorf("%s", "大于200元的红包，必须设置SceneID")
		return
	}
	return true, nil
}

//SendRedPack 发红包
func (c *PayTool) SendRedPack(input RedPackInput) (isSuccess bool, err error) {
	if isGood, err := input.Check(); !isGood {
		return false, err
	}

	now := time.Now()
	dayStr := beego.Date(now, "Ymd")

	billno := c.MchID + dayStr + util.RandomStr(10)

	var signMap = make(map[string]string)
	signMap["nonce_str"] = util.RandomStr(5)
	signMap["mch_billno"] = billno //mch_id+yyyymmdd+10位一天内不能重复的数字
	signMap["mch_id"] = c.MchID
	signMap["wxappid"] = c.AppID
	signMap["send_name"] = input.SendName
	signMap["re_openid"] = input.ToOpenID
	signMap["total_amount"] = util.ToStr(input.MoneyFen)
	signMap["total_num"] = "1"
	signMap["wishing"] = input.Wishing
	signMap["client_ip"] = input.IP
	signMap["act_name"] = input.ActName
	signMap["remark"] = input.Remark
	signMap["sign"] = base.Sign(signMap, c.MchAPIKey, nil)

	respMap, err := c.SendRedPackRaw(signMap)
	if err != nil {
		return false, err
	}

	resultCode, ok := respMap["result_code"]
	if !ok {
		err = errors.New("no result_code")
		return false, err
	}

	if resultCode != "SUCCESS" {
		returnMsg, _ := respMap["return_msg"]
		errMsg, _ := respMap["err_code_des"]
		errCode, _ := respMap["err_code"]

		if errCode == "NOTENOUGH" {
			return false, ErrNoEnoughMoney
		}

		err = fmt.Errorf("Err:%s return_msg:%s err_code:%s err_code_des:%s", "result code is not success", returnMsg, errCode, errMsg)
		return false, err
	}

	mchBillNo, ok := respMap["mch_billno"]
	if !ok {
		err = errors.New("no mch_billno")
		return false, err
	}

	if billno != mchBillNo {
		err = errors.New("billno is not correct")
		return false, err
	}

	return true, nil
}
