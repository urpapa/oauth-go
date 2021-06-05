package code

import (
	"fmt"
	"github.com/google/uuid"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
	"tommy.com/utils"
)

/**
unit test for the Expire Func
*/
func TestCodeExpire(t *testing.T) {
	createdTime, _ := time.Parse("2006-01-02 15:04:05", "2014-06-15 08:37:18")
	c := Code{
		Code:       "123",
		ExpireTIme: 7200,
		State:      "eyz",
		Created:    createdTime,
	}
	fmt.Printf("%v\n", createdTime)
	fmt.Printf("%v\n", time.Now())
	if !c.Expire(time.Now()) {
		t.Failed()
	}
}

func TestGloable(t *testing.T) {
	ch := make(chan map[string]*Session, 1)
	ch <- utils.SESSION
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(4)
	go putSessionInfo("gr1:", &wg, ch)
	go putSessionInfo("gr2:", &wg, ch)
	go putSessionInfo("gr3:", &wg, ch)
	go putSessionInfo("gr4:", &wg, ch)
	wg.Wait()
	fmt.Printf("---------------------*____________________________________\n")
	for key, value := range utils.SESSION {
		fmt.Printf("key is %v,value is %v", key, value)
	}
}

func putSessionInfo(num string, wg *sync.WaitGroup, ch chan map[string]*Session) {
	fmt.Printf("%s in \n", num)
	chs := <-ch
	fmt.Printf("%s get the value \n", num)
	for i := 0; i < 100; i++ {
		m := make(map[string]interface{})
		str := num + strconv.Itoa(i)
		m[num+strconv.Itoa(i)] = &str
		chs[num+uuid.NewString()] = &Session{Uuid: uuid.NewString(), CreatedAt: time.Now(), Name: strconv.Itoa(i), Attributes: &m}
	}
	for key, value := range utils.SESSION {
		Rattrap := *value.Attributes
		Rattrap[num] = &num
		fmt.Printf("grouptings:%v,key is %v,value is %v \n", num, key, value)

	}
	ch <- chs
	fmt.Printf("%s is finished \n", num)
	wg.Done()
}
