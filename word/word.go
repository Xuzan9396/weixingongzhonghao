package word

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
	"sync"
	"weixingongzhonghao/redix2"
)

type SensiInfo struct {
	Word   string `json:"word"`
	Msg    string `json:"msg"`
	IsText string `json:"is_text"`
}

type ReplyInfo struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	IsText    string `json:"is_text"`
	MediaId   string `json:"media_id"`
	WeiXinUrl string `json:"wei_xin_url"`
}

type SensiPlat struct {
	sync.RWMutex
	arrs      []*SensiInfo
	word      map[string]*SensiInfo
	subStr    string
	menu      string
	ios_menu  string
	reply     map[string]*ReplyInfo
	ios_reply map[string]*ReplyInfo
}

var g_SensiPlat *SensiPlat

func GetSensiPlat() *SensiPlat {
	if g_SensiPlat == nil {
		g_SensiPlat = &SensiPlat{}
		g_SensiPlat.SetString()
		g_SensiPlat.SetSubString()
		g_SensiPlat.SetMenuString()
		g_SensiPlat.SetIosMenuString()
		g_SensiPlat.SetReplyString("ios")
		g_SensiPlat.SetReplyString("all")

	}

	return g_SensiPlat
}

// 设置关键字
func (s *SensiPlat) SetString() {

	str, err := redis.String(redix2.RedisX.Get("WEIXINGONGZHONGHAO:WORD"))
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Lock()
	defer s.Unlock()
	s.arrs = make([]*SensiInfo, 0)
	s.word = make(map[string]*SensiInfo)
	err = json.Unmarshal([]byte(str), &s.arrs)
	if err != nil {
		return
	}
	lens := len(s.arrs)
	if lens > 0 {
		for _, value := range s.arrs {
			keys := strings.Split(value.Word, ",")
			if len(keys) > 0 {
				for _, v := range keys {
					v = strings.ToLower(v)
					s.word[v] = value
				}
			}
		}
	}

	fmt.Println(s.word)
}

// 检查当前关键字
func (s *SensiPlat) CheckString(word string) (bools, is_img bool, str string) {
	word = strings.ToLower(word)
	s.RLock()
	defer s.RUnlock()
	strs, bools := s.word[word]
	if !bools {
		return
	}
	if strs.IsText == "0" {
		is_img = true
	}
	str = strs.Msg

	return
}

func (s *SensiPlat) SetSubString() {

	str, err := redis.String(redix2.RedisX.Get("WEIXINGONGZHONGHAO:SUBSCRIBE"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("订阅", str)
	s.subStr = str

}

func (s *SensiPlat) SetIosMenuString() {
	str, err := redis.String(redix2.RedisX.Get("WEIXINGONGZHONGHAO:IOS_MENU"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)
	s.ios_menu = str

}

func (s *SensiPlat) SetReplyString(act string) {
	var str []byte
	var err error
	if act == "ios" {
		s.Lock()
		defer s.Unlock()
		s.ios_reply = make(map[string]*ReplyInfo)
		str, err = redis.Bytes(redix2.RedisX.Get("WEIXINGONGZHONGHAO:IOS_REPLY"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(str))

		if err := json.Unmarshal(str, &s.ios_reply); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s.ios_reply)

	} else {
		s.Lock()
		defer s.Unlock()
		s.reply = make(map[string]*ReplyInfo)
		str, err = redis.Bytes(redix2.RedisX.Get("WEIXINGONGZHONGHAO:REPLY"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(str))

		if err := json.Unmarshal(str, &s.reply); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s.reply)
	}

}

func (s *SensiPlat) SetMenuString() {
	str, err := redis.String(redix2.RedisX.Get("WEIXINGONGZHONGHAO:MENU"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)
	s.menu = str

}

func (s *SensiPlat) GetActString(act string) string {
	var str string
	switch act {
	case "ios":

		str = s.ios_menu
	case "all":
		str = s.menu
	case "sub":
		str = s.subStr

	}
	return str
}

func (s *SensiPlat) GetActReply(act string) (str map[string]*ReplyInfo) {

	switch act {
	case "ios_reply":
		str = s.ios_reply
	case "reply":
		str = s.reply
	}

	return str
}

//
//func (s *SensiPlat)CheckStr(text string)string{
//	defer func() {
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	text = strings.ToLower(text)
//	ns:=int64(1000000)
//	CurrTime:=time.Now().UnixNano()/ns
//	rs := []byte(text)
//	var findChineseCount = 10*3
//	var findMinCount = 2
//	for i:=0;i<= len(text)-findMinCount;i++ {
//		for j:=i+findMinCount+findChineseCount;j>=i+findMinCount;j--{
//			if j>len(rs){
//				j = len(rs)
//			}
//			sub:=string(rs[i:j])
//			dsub := sub
//			//清理空格----------------------
//			reg := regexp.MustCompile("\\s+")
//			dsub = reg.ReplaceAllString(sub, "")
//			//------------------------
//			subf:=dsub
//			bHave,rStr:=s.IsHaveStr(subf)
//			//glog.Infoln("CheckStr ",subf)
//			if bHave {
//				return rStr
//			}
//		}
//	}
//	LastTime:=time.Now().UnixNano()/ns
//	if LastTime-CurrTime > 200{
//		glog.Warning("call proc result:", "file check"," over 200 mill,",LastTime-CurrTime)
//		//fmt.Println("call proc result:", procName," over 100 mill,",LastTime-CurrTime)
//	}
//	return ""
//}
//
//func (s *SensiPlat)IsHaveStr(str string)(bool,string){
//	//str1 := `,`+str
//	//str = str+`,`
//	s.Lock()
//	defer s.Unlock()
//	if s.arrs!=nil{
//		iLen:=len(s.arrs)
//		for i:=0;i<iLen;i++{
//			strArrs:=strings.Split(s.arrs[i].Word, ",")
//			if strArrs!=nil{
//				kLen:=len(strArrs)
//				for k:=0;k<kLen;k++{
//					if strArrs[k] == str{
//						return true,s.arrs[i].Msg
//					}
//				}
//			}
//			/*if strings.Contains(s.arrs[i].Word,str)||strings.Contains(s.arrs[i].Word,str1){
//				return true,s.arrs[i].Msg
//			}*/
//		}
//	}
//	return false,""
//}
//
//func(s *SensiPlat)ReadFromRedis(){
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	str,err:=redis.GetRedisApp().GetV2("SENSITIVE_WORD")
//	if str!=""&&err==nil{
//		s.SetString(str)
//	}
//}
//
//func (s *SensiPlat)ReadUserChatInfo(userId int){
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	if s.IsNeedRead(userId)==false{
//		return
//	}
//	keyStr:="MEMBER:AUTH:"+strconv.Itoa(userId)
//	ret,err:=redis.GetRedisApp().HGetV2(keyStr,"speak_auth")
//	if err==nil&&ret!=""{
//		iValue,_:=strconv.ParseInt(ret,10,64)
//		s.SetChatValue(userId,iValue)
//	}else{
//		s.SetChatValue(userId,0)
//	}
//}
//
//func (s *SensiPlat)ReadUserInfo(userId int){
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	keyStr:="MEMBER:AUTH:"+strconv.Itoa(userId)
//	ret,err:=redis.GetRedisApp().HGetV2(keyStr,"speak_auth")
//	if err==nil&&ret!=""{
//		iValue,_:=strconv.ParseInt(ret,10,64)
//		s.SetChatValue(userId,iValue)
//	}else{
//		s.SetChatValue(userId,0)
//	}
//}
//
//func (s *SensiPlat)IsOfficial(userId int)bool{
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//
//	keyStr:="MEMBER:AUTH:"+strconv.Itoa(userId)
//	ret,err:=redis.GetRedisApp().HGetV2(keyStr,"is_official_auth")
//	if err==nil&&ret!=""{
//		iValue,_:=strconv.ParseInt(ret,10,64)
//		if iValue == 1{
//			return true
//		}
//	}
//	return false
//}
//
//func (s *SensiPlat)SetChatValue(userId int,value int64){
//	s.Lock()
//	defer s.Unlock()
//	ret,ok:=s.arrChats[userId]
//	if ok{
//		ret.ReadTime = time.Now().Unix()
//		ret.Value = value
//	}else{
//		var vv UserChatInfo
//		vv.ReadTime = time.Now().Unix()
//		vv.Value = value
//		vv.userId = userId
//		s.arrChats[userId] = &vv
//	}
//}
//
//func (s *SensiPlat)IsNeedRead(userId int)bool{
//	s.Lock()
//	defer s.Unlock()
//	ret,ok:=s.arrChats[userId]
//	if !ok{
//		return true
//	}
//	if time.Now().Unix()-ret.ReadTime>180{
//		return true
//	}
//	return false
//}
//
//func (s *SensiPlat)IsBanInfo(userId int)(bool,int64){
//	s.Lock()
//	defer s.Unlock()
//	ret,ok:=s.arrChats[userId]
//	if !ok{
//		return false,0
//	}
//	if ret.Value!=0{
//		if ret.Value==-1{
//			return true,ret.Value
//		}
//		currTime:=time.Now().Unix()
//		if currTime < ret.Value{
//			return true,ret.Value
//		}
//	}
//	return false,0
//}
//
//func (s *SensiPlat)IsBanChat(userId int)(bool,int64){
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	s.ReadUserChatInfo(userId)
//	if ok,bTime:=s.IsBanInfo(userId);ok{
//		return true,bTime
//	}
//	return false,0
//}
//
//func (s *SensiPlat)GetChatUserInfo()[]UserChatInfo{
//	s.Lock()
//	defer s.Unlock()
//	arrs:=make([]UserChatInfo,0)
//	if s.arrChats==nil{
//		return arrs
//	}
//	var count int = 0
//	for _,v:=range s.arrChats{
//		var vv UserChatInfo
//		vv=*v
//		arrs = append(arrs,vv)
//		count++
//		if count>300{
//			break
//		}
//	}
//	if len(s.arrChats)>5000{
//		s.arrChats = make(map[int]*UserChatInfo)
//		return nil
//	}
//	return arrs
//}
//
//func (s *SensiPlat)DelChatInfo(userId int){
//	s.Lock()
//	defer s.Unlock()
//	if s.arrChats == nil{
//		return
//	}
//	delete(s.arrChats,userId)
//}
//
////清理长时间未登陆用户
//func (s *SensiPlat)ClearChatInfo(){
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	arrs:=s.GetChatUserInfo()
//	if arrs==nil{
//		return
//	}
//	iLen:=len(arrs)
//	if iLen <= 0{
//		return
//	}
//	currTime:=time.Now().Unix()
//	for i:=0;i<iLen;i++{
//		if currTime - arrs[i].ReadTime>3600*24*1{
//			s.DelChatInfo(arrs[i].userId)
//		}
//	}
//}

//func (s *SensiPlat)Run(){
//	defer func(){
//		if err := recover(); err != nil {
//			glog.Errorln(string(debug.Stack()))
//		}
//	}()
//	s.ReadFromRedis()
//	tick := time.Tick(30 * time.Second) //2秒执行一次
//	for {
//		select {
//		case <-tick:
//			go s.ReadFromRedis()
//			s.ClearChatInfo()
//		}
//	}
//}
