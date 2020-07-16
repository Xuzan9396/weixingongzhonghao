package weixin

import (
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"weixingongzhonghao/config"
	"weixingongzhonghao/redix"
)

type WeixinConfig struct {
	OffCount *officialaccount.OfficialAccount
}

var g_wc *WeixinConfig

func GetOa() *WeixinConfig {

	if g_wc == nil {
		g_wc = &WeixinConfig{}
		wc := wechat.NewWechat()

		memory := redix.GetReids().Redis

		cfg := &offConfig.Config{
			//AppID:     "wxc510627a460bc5af", // 夜空
			//AppSecret: "9c64db2d4daba73978b466a807d887d7",
			//AppID:     "wxf47f8e54a15168c1", // 浪舞
			//AppSecret: "4163d30ca9e0c91181eb2b09cfdc52ae",
			//AppID:     "wxb59c2e1b39a0ad33", // 测试
			//AppSecret: "86d5c2e8bb20380e1fc584a8707eb0fe",
			//Token:     "sdmcnjf67234hzaslkiashd892",

			AppID:     config.G_JsonConfig.Wx.AppID, // 测试
			AppSecret: config.G_JsonConfig.Wx.AppSecret,
			Token:     config.G_JsonConfig.Wx.Token,
			//EncodingAESKey: "xxxx",  // zCyLq87F7Q6r67Ubs3oYGPVMEBN0ItAFkiPBiGuusSx
			Cache: memory,
		}
		g_wc.OffCount = wc.GetOfficialAccount(cfg)
	}
	ak, _ := g_wc.OffCount.GetAccessToken()
	fmt.Println("accesstoken", ak)

	return g_wc
}
