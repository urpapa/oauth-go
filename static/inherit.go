package static

import (
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

func Login(w http.ResponseWriter,r *http.Request){
	T["login.html"].ExecuteTemplate(w,"base",nil)
}
func Userinfo(w http.ResponseWriter,r *http.Request){
	usernname,_:=utils.GetSession(r,"username")
	T["user.html"].ExecuteTemplate(w,"base",&code.User{UserName: usernname.(string),UserCode: ""})
}
