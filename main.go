package main

import (
	json2 "encoding/json"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"runtime"
	"tommy.com/oauth"
	"tommy.com/static"
	"tommy.com/utils"
)

//go Main func for Oauth-go Project
func main() {

	//scheduel clear the cache
	runtime.GOMAXPROCS(4)
	go utils.TimeClear()
	// starting up the server
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/oauth/authorize", authorization)
	mux.HandleFunc("/oauth/token",token)
	mux.HandleFunc("/login",static.Login)
	mux.HandleFunc("/user",static.Userinfo)
	mux.HandleFunc("/dologin",doLoing)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server := &http.Server{
		Addr:           "127.0.0.1:9999",
		Handler:        mux,
	//	ReadTimeout:    time.Duration(8 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func staticFiles() {
	dir := http.Dir("./static")
	handler := http.StripPrefix("/static/", http.FileServer(dir))
	http.Handle("/static/", handler)
}

/**
auth授权接口
*/
func authorization(w http.ResponseWriter, r *http.Request) {
	authenParam, err := oauth.AuthenPara(w, r)
	if err == nil {
		if _, ok := oauth.IsLogedIn(w, r); ok {
			code := oauth.Authorize(authenParam)
			http.Redirect(w, r, authenParam.RedirectUri+"&code="+code.Code, 200)
		}
	} else {
		//判断本次请求是否已经登录
		json, _ := json2.Marshal(NewRes("403", err.Error()))
		w.Write(json)
	}


}
func token(w http.ResponseWriter, r *http.Request) {
	token,err:=oauth.AccessToken(w,r)
	if err==nil{
	json,_ :=json2.Marshal(token)
	w.Write(json)
	} else {
		//判断本次请求是否已经登录
		json, _ := json2.Marshal(NewRes("403", err.Error()))
		w.Write(json)
	}


}

func doLoing(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	username:=r.PostFormValue("username")
	password:=r.PostFormValue("password")
	if username=="test"&& password=="000000"{
		utils.SetSession(&w,r,"username",username)
		http.Redirect(w,r,"/user",302)
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
