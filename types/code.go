package code

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"
)

/**
code one time use
*/
type Code struct {
	Code, State string
	ExpireTIme  int
	Created     time.Time
}

var config Configuration

func (c *Code) Expire(t time.Time) bool {
	duration, _ := time.ParseDuration(strconv.Itoa(c.ExpireTIme) + "s")
	if c.Created.Add(duration).After(t) {
		return true
	}
	return false

}

//AccessToken
type AccessToken struct {
	Token        string
	ExpireTime   int
	RefreshToken string
	CreateTime   time.Time
	UserName     string
}

/**
user info
*/
type User struct {
	UserName, UserCode, Password string
}

type Configuration struct {
	ClientId, ClientSecret string
}

type Expire interface {
	Expire(t time.Time) bool
}

type Session struct {
	Uuid       string
	Name       string
	CreatedAt  time.Time
	Attributes *map[string]interface{}
}

func newSession(userName string) Session {
	return Session{Uuid: uuid.NewString(), Name: userName, CreatedAt: time.Now()}
}

func (token *AccessToken) Expire(t time.Time) bool {
	duration, _ := time.ParseDuration(strconv.Itoa(token.ExpireTime) + "s")
	if token.CreateTime.Add(duration).After(t) {
		return true
	}
	return false
}

/**
创建新的code 并存放到缓存中
*/
func NewCode() Code {
	return Code{
		Code:       uuid.NewString(),
		ExpireTIme: 7200,
		Created:    time.Now(),
	}
}

func NewAccessToken() AccessToken {
	return AccessToken{
		Token: uuid.NewString(), ExpireTime: 7200,
		CreateTime:   time.Now(),
		RefreshToken: uuid.NewString(),
	}

}

/**
load the config.json info
*/
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

/**
初始化配置文件信息
*/
func init() {
	loadConfig()
}

/**
校验clientId,clientSecret是有效
*/
func CheckClient(clientId string, clientSecret string) bool {
	if clientId == config.ClientId && clientSecret == config.ClientSecret {
		return true
	}
	return false
}
