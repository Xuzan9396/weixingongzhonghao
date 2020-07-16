package handle

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"weixingongzhonghao/redix2"
	"weixingongzhonghao/weixin"
	"weixingongzhonghao/word"
)

type ErrCodes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Menuid  int64  `json:"menuid"`
}

type RespImgInfo struct {
	Media_id string `json:"media_id"`
	Url      string `json:"url"`
}

type ReqImg struct {
	ImgPath string `json:"img_path"`
}

// 默认菜单列表
func GetMenu(rw http.ResponseWriter, req *http.Request) {

	wxModel := weixin.GetOa()
	officialAccount := wxModel.OffCount

	m := officialAccount.GetMenu()
	resp, _ := m.GetMenu()
	//resp,_ := m.GetCurrentSelfMenuInfo()
	bytes, _ := json.Marshal(resp)
	rw.Write(bytes)

}

//个性菜单列表
func GeCustomtMenu(rw http.ResponseWriter, req *http.Request) {

	wxModel := weixin.GetOa()
	officialAccount := wxModel.OffCount

	m := officialAccount.GetMenu()
	//resp,_ := m.GetMenu()
	resp, _ := m.GetCurrentSelfMenuInfo()
	bytes, _ := json.Marshal(resp)
	rw.Write(bytes)

}

// 发送http
func HttpPostGet(url, method, params string, headerArr ...map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	request, err := http.NewRequest(method, url, strings.NewReader(params))

	// 设置header
	if method == "POST" {
		//headerArr[0]["Content-Type"] = "application/json;charset=utf-8"
		request.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	if len(headerArr) > 0 {
		for key, value := range headerArr[0] {
			//req.Header.Add("Authorization", "3eex8dY04BdU1amui6bf20ECgtyc9s")
			request.Header.Add(key, value)
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	return body, err
}

func RedisPush() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	fmt.Println("push 启动")
	for {
		acts, err := redis.Values(redix2.RedisX.CommonCmdAct("BLPOP", "WEIXINGONGZHONGHAO:EVENT", 0))
		if err != nil {
			continue
		}
		if acts == nil || len(acts) <= 0 {
			return
		}
		smsStr := string(acts[1].([]byte))
		switch smsStr {
		case "menu":
			word.GetSensiPlat().SetMenuString() // 设置关键字
			word.GetSensiPlat().SetReplyString("all")
			setMenu("all")
		case "ios_menu":
			word.GetSensiPlat().SetIosMenuString() // 设置关键字
			word.GetSensiPlat().SetReplyString("ios")
			setMenu("ios")
		case "android_menu":

			setMenu("android")
		case "subscribe":

			word.GetSensiPlat().SetSubString()

		case "word":
			word.GetSensiPlat().SetString()
		}

	}

	switch {

	}

	//memory :=  redix.GetReids().Redis
	//redix.GetReids().Redis.

}

// 设置默认菜单列表
func setMenu(act string) {

	//if req.Method == "POST" {
	//	req.ParseForm()
	//	fmt.Println(req.PostFormValue("act"))
	//	act := req.PostFormValue("act")
	var url, jsonStr string
	switch act {
	case "ios":
		url = "https://api.weixin.qq.com/cgi-bin/menu/addconditional"
		jsonStr = word.GetSensiPlat().GetActString("ios")
		//			jsonStr =`{
		//	"button": [{
		//			"name": "夜空直播",
		//			"sub_button": [{
		//					"type": "click",
		//                    "name": "下载地址",
		//					"key": "download"
		//				}
		//			]
		//		},
		//		{
		//			"name": "联系我们",
		//			"sub_button": [{
		//					"type": "click",
		//					"name": "运营客服",
		//					"key": "kefu"
		//				}
		//			]
		//		},
		//        {
		//			"name": "夜空服务",
		//			"sub_button": [
		//                {
		//					"type": "view",
		//					"name": "优币服务",
		//					"url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxc510627a460bc5af&redirect_uri=https://h5cdn.wangran.live/wxpublickpay?&response_type=code&scope=snsapi_base&state=wft&connect_redirect=1#wechat_redirect"
		//				},
		//                {
		//					"type": "view",
		//					"name": "贵族服务",
		//					"url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxc510627a460bc5af&redirect_uri=https://h5cdn.wangran.live/mall?&response_type=code&scope=snsapi_base&state=wft&connect_redirect=1#wechat_redirect"
		//				},
		//                {
		//					"type": "view",
		//					"name": "靓号服务",
		//					"url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxc510627a460bc5af&redirect_uri=https://h5cdn.wangran.live/nicenum/home?&response_type=code&scope=snsapi_base&state=wft&connect_redirect=1#wechat_redirect"
		//				}
		//			]
		//		}
		//	],
		//	"matchrule": {
		//		"client_platform_type": "1",
		//		"language": "zh_CN"
		//	}
		//}`
	case "android":
		url = "https://api.weixin.qq.com/cgi-bin/menu/addconditional"
		jsonStr = `{
	"button": [{
			"name": "夜空直播安卓",
			"sub_button": [{
					"type": "click",
                    "name": "下载地址",
					"key": "download"
				}
			]
		},
		{
			"name": "联系我们",
			"sub_button": [{
					"type": "click",
					"name": "运营客服",
					"key": "kefu"
				}
			]
		},
        {
			"name": "夜空服务",
			"sub_button": [
                {
					"type": "view",
					"name": "优币服务",
					"url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxc510627a460bc5af&redirect_uri=https://h5cdn.wangran.live/wxpublickpay?&response_type=code&scope=snsapi_base&state=wft&connect_redirect=1#wechat_redirect"
				},
                {
					"type": "view",
					"name": "贵族服务",
					"url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxc510627a460bc5af&redirect_uri=https://h5cdn.wangran.live/mall?&response_type=code&scope=snsapi_base&state=wft&connect_redirect=1#wechat_redirect"
				},
                {
					"type": "view",
					"name": "靓号服务",
					"url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxc510627a460bc5af&redirect_uri=https://h5cdn.wangran.live/nicenum/home?&response_type=code&scope=snsapi_base&state=wft&connect_redirect=1#wechat_redirect"
				}
			]
		}
	],
	"matchrule": {
		"client_platform_type": "1",
		"language": "zh_CN"
	}
}`
	case "all":
		url = "https://api.weixin.qq.com/cgi-bin/menu/create"
		jsonStr = word.GetSensiPlat().GetActString("all")

		//			jsonStr =`{
		//    "button":[
		//        {
		//            "name":"直播1",
		//            "sub_button":[
		//                {
		//                    "type":"click",
		//                    "name":"下载地址",
		//                    "key":"download"
		//                }
		//            ]
		//        },
		//        {
		//            "name":"联系我们",
		//            "sub_button":[
		//                {
		//                    "type":"click",
		//                    "name":"运营客服",
		//                    "key":"kefu"
		//                }
		//            ]
		//        }
		//    ]
		//}`
	default:
		fmt.Println("方式不对")
		//rw.Write([]byte("fail"))
		return

	}

	respByte, err := setMenuPub(url, jsonStr)
	if err != nil {
		fmt.Println("fail")
		//rw.Write([]byte("fail"))
		return
	}
	datas := ErrCodes{}
	if err := json.Unmarshal(respByte, &datas); err != nil {
		fmt.Println(err)
		fmt.Println("fail")

		//rw.Write([]byte("fail"))
		return
	}
	if datas.ErrMsg == "ok" || datas.Menuid > 0 {
		//rw.Write([]byte("ok"))
		fmt.Println("ok")

	} else {
		fmt.Println("fail")

		//rw.Write([]byte("fail"))
	}
	//}

}

func setMenuPub(url, jsonStr string) (respByte []byte, err error) {
	wxModel := weixin.GetOa()
	officialAccount := wxModel.OffCount
	m := officialAccount.GetMenu()
	accessToken, err := m.GetAccessToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	//
	uri := fmt.Sprintf("%s?access_token=%s", url, accessToken)
	respByte, err = HttpPostGet(uri, "POST", jsonStr)
	fmt.Println(string(respByte))
	return

}

// 回调接收的消息
func ServeWechat(rw http.ResponseWriter, req *http.Request) {

	////memory := cache.NewMemory()

	wxModel := weixin.GetOa()

	officialAccount := wxModel.OffCount

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)

	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		//TODO
		fmt.Println(msg.Event, msg.EventKey, msg.Content, msg.MsgType, msg)
		switch msg.MsgType { // 事件类型
		case "event":
			switch msg.Event { // 事件行为
			case "subscribe":

				/*亲亲~欢迎来到夜空直播~

				你想看性感御姐的舞姿吗？

				你想跟可爱的萝莉窃窃私语吗？

				你想跟小清新美女低声密谈吗？

				<a href="https://h5cdn.wangran.live/download">~来这里，</a>你想要都有！（ios公众号充值已下线，请联系客服）*/
				text := message.NewText(word.GetSensiPlat().GetActString("sub"))
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			case "CLICK": // 点击按钮
				//msg.
				iosReply := word.GetSensiPlat().GetActReply("reply")
				iosReplys, ok := iosReply[msg.EventKey]
				if ok {
					if iosReplys.IsText == "1" {
						text := message.NewText(iosReplys.Value)
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					} else if iosReplys.IsText == "0" {
						image := message.NewImage(iosReplys.MediaId)
						return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
					}
				} else {
					reply := word.GetSensiPlat().GetActReply("ios_reply")
					replys, ok := reply[msg.EventKey]
					if ok {
						if replys.IsText == "1" {
							text := message.NewText(replys.Value)
							return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
						} else if replys.IsText == "0" {
							image := message.NewImage(replys.MediaId)
							return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
						}
					}
				}

				//switch msg.EventKey {
				//case "download":
				//	image := message.NewImage("qIKXAZDnG4gaHeVts7sWd6Kpzasz0hIt_c8kGrCj0Do")
				//	return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
				//case "kefu":
				//	image := message.NewImage("qIKXAZDnG4gaHeVts7sWd0GGGgiiDwdomZQtKMGFDUI")
				//	return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
				//}

			}
		case "text":
			model := word.GetSensiPlat()
			bools, img_bool, str := model.CheckString(msg.Content)
			if bools {
				if img_bool {
					image := message.NewImage(str)
					return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
				} else {
					text := message.NewText(str)
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}

			}

			//switch msg.Content {
			//case "test":
			//	article1 := message.NewArticle("夜空直播", "夜空app下载地址", "http://mmbiz.qpic.cn/mmbiz_png/R3kNehdENMsibSJrCib0zibya76w5sDotfEsPOyvWaeSsc7t60c7v3jrZMZnVd1Afrd0PSp1gxaibiabpCpXIaU95rw/0?wx_fmt=png", "https://h5cdn.wangran.live/download")
			//
			//
			//	articles := []*message.Article{article1}
			//	news := message.NewNews(articles)
			//	return &message.Reply{MsgType: message.MsgTypeNews, MsgData: news}
			//
			//	//text := message.NewText(msg.Content)
			//	//return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			//case "send_all":
			//	bd:=officialAccount.GetBroadcast()
			//	bd.SendText(nil,"textssss")
			//}

		}

		return nil
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

// 获取素材列表
func GetMaterial(rw http.ResponseWriter, req *http.Request) {
	wxModel := weixin.GetOa()
	officialAccount := wxModel.OffCount
	m := officialAccount.GetMaterial()
	list, _ := m.BatchGetMaterial(material.PermanentMaterialTypeImage, 0, 20)

	//if list.TotalCount > 0 {
	//	sendList := make([]*ImgInfo,list.ItemCount)
	//	for key, value := range list.Item {
	//		sendList[key] = &ImgInfo{
	//			Media_id:value.MediaID,
	//			Name:value.Name,
	//			Url:value.URL,
	//		}
	//
	//	}
	//	sendListStr,_ := json.Marshal(sendList)
	//	fmt.Println(string(sendListStr));
	//
	//	redix2.RedisX.Set("WEIXINGONGZHONGHAO:IMG_LIST",string(sendListStr))
	//}
	//if err != nil{
	//	fmt.Println(err,1111111);
	//	return
	//}
	//
	fmt.Println(list)
	bytes, _ := json.Marshal(list)

	rw.Write(bytes)
}

// 上传图片
func SetMaterial(rw http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("read body err, %v\n", err)
			return
		}
		println("json:", string(body))

		dataImg := ReqImg{}
		if err := json.Unmarshal(body, &dataImg); err != nil {
			fmt.Println(err, "解析错误")
			return
		}

		wxModel := weixin.GetOa()
		officialAccount := wxModel.OffCount
		m := officialAccount.GetMaterial()
		mediaID, url, err := m.AddMaterial(material.MediaTypeImage, dataImg.ImgPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := RespImgInfo{
			Media_id: mediaID,
			Url:      url,
		}

		bytes, _ := json.Marshal(data)

		rw.Write(bytes)
	}

}

func Test(rw http.ResponseWriter, req *http.Request) {
	// tesolt
	//ip := ClientIP(req)
	//fmt.Println("ip:", ip);
	//fmt.Println(req.Header);
	////fmt.Println(req.URL.Query());
	//for key, value := range req.URL.Query() {
	//	if key == "PHPSESSID" {
	//		fmt.Println(key,value);
	//	}
	//}

	//rw.Header().Set("Access-Control-Allow-Origin","*")
	//rw.Write([]byte("test 测试"))
	//data := struct {
	//	Code int
	//	Msg string
	//}{ 200,"验证成功"}
	//msg, _ := json.Marshal(data)
	rw.Write([]byte("test"))

	fmt.Println("test 测试")
}

//func ClientIP(r *http.Request) string {
//	xForwardedFor := r.Header.Get("X-Forwarded-For")
//	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
//	if ip != "" {
//		return ip
//	}
//
//	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
//	if ip != "" {
//		return ip
//	}
//
//	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
//		return ip
//	}
//
//	return ""
//}
