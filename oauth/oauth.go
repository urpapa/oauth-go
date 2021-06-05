package oauth

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"tommy.com/constants"
	code "tommy.com/types"
	"tommy.com/utils"

	"strings"
)

/**
此方法前保证clientId,clientSecret 登录等已经校验
*/
func Authorize(para *AuthenStruct) *code.Code {
	if para.ResponseType == "code" {
		c := code.NewCode()
		utils.C.Store(c.Code,&c)
		return &c
	}
	return nil
}

func AccessToken(w http.ResponseWriter, r *http.Request) (token code.AccessToken, err error) {
	authorize := r.Header["Authorization"][0]
	clientid, clientsecret, err := DecodeAuthHeader(authorize)
	if err != nil {
	return
	}
	r.ParseMultipartForm(1024)
	tokenStruct:= TokenStruct{ClientId: clientid, ClientSecret: clientsecret,UserName: r.FormValue(constants.SESSIOIN_USERNAME),Password: r.FormValue("password"), Code: r.FormValue("code"),GrantType: r.FormValue("grant_type")}
	switch tokenStruct.GrantType {
	case constants.TYPECODE:
		log.Println("type is authorization_code ")
		c := utils.GetCode(tokenStruct.Code)
		if c == "" {
			err = errors.New("code is valid")
		} else {
			token = code.NewAccessToken()
			token.UserName = tokenStruct.UserName
			//将token存入到缓存中
			utils.U.Store(token.Token,&token)
		}

		break

	case constants.TYPEPASSWORD:
		log.Println("grant_type is password")
		if strings.Contains(tokenStruct.UserName,"user") && tokenStruct.Password == "000000" {
			token = code.NewAccessToken()
			token.UserName = tokenStruct.UserName
			//将token存入到缓存中
		     utils.U.Store(token.Token,&token)
			//将登录信息写入到cookie中方便测试
			utils.SetSession(&w,r,"username",tokenStruct.UserName)
		} else {
			err = errors.New("username or password invalid")
		}


		break

	default:
		err = errors.New("the grantType is invalid")
		break
	}
	return
}

func Profile(token string) (user code.User, err error) {
	t, eer := utils.GetToken(token)
	if eer == nil {
		err = errors.New("token is valid")
	}
	user = code.User{
		UserName: t.UserName,
	}
	return
}

type TokenStruct struct {
	GrantType, ClientId, ClientSecret, UserName, Password, Scope, State, Code string
}

type AuthenStruct struct {
	ResponseType, ClientId, RedirectUri, State, ClientSecret, Scope string
}

/**
解析authentication 参数。 Basic test:test
test:test需要进行basse64加密
*/
func DecodeAuthHeader(authentication string) (clientId string, ClientSecret string, err error) {
	if authentication != "" && strings.HasPrefix(authentication, "Basic ") {
		tmps := strings.Split(authentication, " ")
		cs, err := base64.StdEncoding.DecodeString(tmps[1])
		if err != nil {
			return "", "", err
		}
		strs := strings.Split(string(cs), ":")
		return strs[0], strs[1], nil
	}
	return "", "", errors.New("authorization header is not exists")
}

/**
解析authorize方法参数
*/
func AuthenPara(m http.ResponseWriter, r *http.Request) (*AuthenStruct, error) {
	authorize := r.Header["Authorization"][0]
	clientid, clientsecret, err := DecodeAuthHeader(authorize)
	if err != nil {
		return nil, err
	}
	values := r.URL.Query()

	return &AuthenStruct{ResponseType: "code", ClientId: clientid, ClientSecret: clientsecret, State: values.Get(constants.OATUT_STATE), Scope: values.Get(constants.OATUT_SCOPE)}, nil

}

/**
解析token接口的参数
*/
func TokenPara(w http.ResponseWriter, r *http.Request) (*TokenStruct, error) {
	authorize := r.Header["Authorization"][0]
	clientid, clientsecret, err := DecodeAuthHeader(authorize)
	if err != nil {
		return nil, err
	}
	r.ParseForm()
	return &TokenStruct{ClientId: clientid, ClientSecret: clientsecret,
		UserName: r.PostFormValue(constants.SESSIOIN_USERNAME),
		Password: r.PostFormValue("password"), Code: r.PostFormValue("code")}, nil

}

/**
校验当前session是否已经登录
*/
func IsLogedIn(w http.ResponseWriter, r *http.Request) (string, bool) {
	username, error := utils.GetSession(r, constants.SESSIOIN_USERNAME)
	if error != nil {
		return "", false
	}
	return username.(string), true
}
