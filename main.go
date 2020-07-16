package main

import (
	"fmt"
	"net/http"
	"weixingongzhonghao/config"
	"weixingongzhonghao/handle"
	"weixingongzhonghao/redix2"
	"weixingongzhonghao/word"
)

/*


{
    "media_id": "qIKXAZDnG4gaHeVts7sWd6Kpzasz0hIt_c8kGrCj0Do",
    "url": "http://mmbiz.qpic.cn/mmbiz_png/R3kNehdENMsibSJrCib0zibya76w5sDotfEsPOyvWaeSsc7t60c7v3jrZMZnVd1Afrd0PSp1gxaibiabpCpXIaU95rw/0?wx_fmt=png",
    "item": []
}

{
    "media_id": "qIKXAZDnG4gaHeVts7sWd0GGGgiiDwdomZQtKMGFDUI",
    "url": "http://mmbiz.qpic.cn/mmbiz_jpg/R3kNehdENMsibSJrCib0zibya76w5sDotfE4I03icTFd2OpUHFbsCHnuN9xN48soum9iax7gDe4HUVqCKWqLBicBs0Hw/0?wx_fmt=jpeg",
    "item": []
}

*/

// GOOS=linux go build -ldflags "-s -w"  -o weixingongzhonghao ./main.go

func main() {

	if err := config.InitConfigJson(); err != nil {
		fmt.Println(err, "配置项错误")
		return
	}

	redix2.InitRedis()
	word.GetSensiPlat() // 设置关键字

	http.HandleFunc("/go_weixin/check", handle.ServeWechat)
	http.HandleFunc("/go_weixin/test", handle.Test)
	http.HandleFunc("/go_weixin/menu", handle.GetMenu)
	http.HandleFunc("/go_weixin/custom_menu", handle.GeCustomtMenu) // 个性菜单列表
	http.HandleFunc("/go_weixin/material", handle.GetMaterial)      // 批量获取素材
	http.HandleFunc("/go_weixin/set_material", handle.SetMaterial)  // 上传素材
	//http.HandleFunc("/go_weixin/set_menu", handle.SetMenu) // 个性菜单列表
	go handle.RedisPush()
	fmt.Println("wechat server listener at", config.G_JsonConfig.Http.Port)
	err := http.ListenAndServe(config.G_JsonConfig.Http.Port, nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
