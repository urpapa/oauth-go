package utils

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"tommy.com/constants"
)

func TestSetGetSession(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/post/", FakePost)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/post/1", nil)
	mux.ServeHTTP(writer, request)
}

func setGet(w http.ResponseWriter, r *http.Request, attr string, obj interface{}, wg *sync.WaitGroup) {

	for i := 0; i < 10; i++ {
		c := http.Cookie{
			Name: constants.SESSION_KEY, Value: uuid.NewString(), HttpOnly: true,
		}
		r.AddCookie(&c)
		SetSession(&w, r, attr+strconv.Itoa(i), obj)
	}
	for i := 0; i < 10000; i++ {
		obj, err := GetSession(r, attr+strconv.Itoa(i))
		fmt.Printf("obj is %v, the err info is %v \n", obj, err)
	}

}
func FakePost(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go setGet(w, r, "go1", "go1-obj", &wg)
	//go setGet(w, &temp2, "go2", "go2-obj",&wg)
	//go setGet(w, &temp3, "go3", "go3-obj",&wg)
	//go setGet(w, &temp4, "go4", "go4-obj",&wg)
	wg.Wait()
}

func TestHashInt(t *testing.T) {
	log.Printf("test:  %v",HashInt("测试一下") %runtime.NumCPU())
}