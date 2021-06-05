package static

import (
	json2 "encoding/json"
	"html/template"
	"net/http"
	code "tommy.com/types"
	"tommy.com/utils"
)

/**
capitalize the first Letter  for public
*/
var T map[string]*template.Template

/**
初始化页面，利用template方式，减少代码
*/
func init() {
	basePath := "./static/base.html"
	T = make(map[string]*template.Template)
	temp := template.Must(template.ParseFiles(basePath, "./static/user.html"))
	T["user.html"] = temp
	temp = template.Must(template.ParseFiles(basePath, "./static/login.html"))
	T["login.html"] = temp
}

func Login(w http.ResponseWriter, r *http.Request) {
	T["login.html"].ExecuteTemplate(w, "base", nil)
}
func Userinfo(w http.ResponseWriter, r *http.Request) {
	usernname, _ := utils.GetSession(r, "username")
	if usernname != nil {
		u := code.GetUser(usernname.(string))
		T["user.html"].ExecuteTemplate(w, "base", u)
	} else {
		json, _ := json2.Marshal(NewRes("403", "login failed"))
		w.Write(json)
	}

}

type resJson struct {
	code, msg string
}

func NewRes(code string, msg string) resJson {
	return resJson{code: code, msg: msg}
}
