package utils

import (
	"errors"
	"github.com/google/uuid"
	"hash/crc32"
	"net/http"
	"runtime"
	"sync"
	"time"
	"tommy.com/constants"
	"tommy.com/types"
)

/**
Manange the the accessToken, code
the dataStruct  :
accessToken is  map[token]*AccessTokenUser
step0:request ,check clientid,clientSecret ,whether the user is loged in ,if not log in with userName and password and Generate the Cookie
setp1: if info is correct ,generate the code
step2: if the code is exists ,generate token and the user info ,delete the code
stpe3: get userInfo with the Given token
*/
//key is client_id  value is  types.code
// map[string]*code.Code

/**
性能优化问题：CODE,USER,SESSION当前都是现成安全型Map 如果多核环境下也只有一个gogroutine在处理，提升效率 需要按照当前系统的cpu core数量进行sync.map创建
增加并发能力，提升吞吐量
 */
//var C sync.Map

//key is token value is tokenInfo include AccessToken and UserInfo
// map[string]*code.AccessToken
//var U sync.Map

//map[string]*code.Session
//var SESSION sync.Map

var  codeSlice []*sync.Map
var  tokenSlice []*sync.Map
var  sessionSlice []*sync.Map


func initMap(){
	codeSlice=make([]*sync.Map,runtime.NumCPU(),runtime.NumCPU())
	tokenSlice=make([]*sync.Map,runtime.NumCPU(),runtime.NumCPU())
	sessionSlice=make([]*sync.Map,runtime.NumCPU(),runtime.NumCPU())
	for i:=0;i<runtime.NumCPU();i++{
		var s1 sync.Map
		var s2 sync.Map
		var s3 sync.Map
		codeSlice[i]=&s1
		tokenSlice[i]=&s2
		sessionSlice[i]=&s3
	}
}

func init() {
initMap()
}

/**
根据client obtain the code and remove the code
use only once
*/
func GetCode(c string) string {
	C:=codeSlice[HashInt(c)%runtime.NumCPU()]
	if tp, ok := C.Load(c); ok {
		//方法调用完成后清理code
		defer C.Delete(c)
		//返回code信息
		return tp.(*code.Code).Code
	}
	return ""
}
/**
get AccessToken,
*/
func GetToken(token string) (*code.AccessToken, error) {
	T:=tokenSlice[HashInt(token)%runtime.NumCPU()]
		if tp, ok := T.Load(token); ok {
			return tp.(*code.AccessToken), nil
	}
	return nil, errors.New("token is not exists")
}

func GetSessionMap(sessionkey string ) *sync.Map{
	return sessionSlice[HashInt(sessionkey)%runtime.NumCPU()]

}

func GetCodeMap(code string) *sync.Map{
	return codeSlice[HashInt(code)%runtime.NumCPU()]

}
func GetTokenMap(token string) *sync.Map{
	return tokenSlice[HashInt(token)%runtime.NumCPU()]

}

/**
循环，清除过期的code 以及token 信息
*/
func Clear() {
	for _,C:= range codeSlice{
		C.Range(func(key, value interface{}) bool {
			if value.(*code.Code).Expire(time.Now()) {
				C.Delete(key)
				return true
			}
			return false
		})
	}
   for _,T:=range tokenSlice{
	   T.Range(func(key, value interface{}) bool {
		   if value.(*code.AccessToken).Expire(time.Now()) {
			   T.Delete(key)
			   return true
		   }
		   return false
	   })
   }

}


/**
向当前的session中存入值，操作不走：
setp1:判断当前cookie中是否有sessionId如果没有说明首次登陆，创建sessionId并写入到本次请求的cookie中，需要考虑性能问题
setp2:利用sessionId从Session中获取对应的map如果存在则获取，如果不存在则创建并存入到当前session中
*/
func SetSession(w *http.ResponseWriter, r *http.Request, attr string, obj interface{}) {
	//根据key 获取cookie在本次请求中是否存在，如果不存在则需要创建cookie,并将cookie写入到返回中
	cookie, err := r.Cookie(constants.SESSION_KEY)
	var sessionId string
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:     constants.SESSION_KEY,
			Value:    uuid.NewString(),
			HttpOnly: true,
		}
		//将uuid写入到cookie中

		http.SetCookie(*w, cookie)
	}
	//sessionid的值
	sessionId = cookie.Value
	var attrMap map[string]interface{}
	syncmap:=GetSessionMap(sessionId)
	tmps, ok := syncmap.Load(sessionId)
	var s *code.Session
	if !ok {
		attrMap = make(map[string]interface{})
		s = &code.Session{Uuid: sessionId, Name: "", Attributes: &attrMap, CreatedAt: time.Now()}
		syncmap.Store(sessionId, s)
	} else {
		s = tmps.(*code.Session)
	}
	//将具体存session内容存入session中
	attrMap = *s.Attributes
	attrMap[attr] = obj

}

var NOSessionValue = errors.New("have no session value for the given key")

func GetSession(r *http.Request, attr string) (interface{}, error) {
	cookie, err := r.Cookie(constants.SESSION_KEY)
	if err == http.ErrNoCookie {
		return nil, NOSessionValue
	}
	syncmap:=GetSessionMap(cookie.Value)
	tmps, ok := syncmap.Load(cookie.Value)
	if ok {
		session := tmps.(*code.Session)
		attrMap := *session.Attributes
		value, ok2 := attrMap[attr]
		if ok2 {
			return value, nil
		}
	}
	return nil, NOSessionValue
}
/**
计算出一个非负的整数
 */
func HashInt(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}


/**
task to clear the cache info
*/
var tricker = time.NewTicker(2 * time.Minute)

/**
定时清理
*/
func TimeClear() {
	<-tricker.C
	Clear()
}
