package util

import (
	"be/option"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

type cookieMgr struct {
	s *securecookie.SecureCookie
}

var CM *cookieMgr

func InitCM() {
	s := securecookie.New([]byte(*option.Cookie01), []byte(*option.Cookie02))
	CM = &cookieMgr{s: s}
}

func (c *cookieMgr) Set(key string, value string, res http.ResponseWriter) string {
	expiration := time.Now().Add(time.Duration(24) * time.Hour)
	if encoded, err := c.s.Encode(key, value); err == nil {
		cookie := &http.Cookie{
			Name:    key,
			Value:   encoded,
			Expires: expiration,
			Path:    "/",
		}
		http.SetCookie(res, cookie)
		return encoded
	} else {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("Cookie Set失败")
		return ""
	}
}

func (c *cookieMgr) Get(key string, req *http.Request) (string, error) {
	// if *option.Mode == "DEV" {
	// 	return "ADMIN-DEV-TOKEN", nil
	// }

	// if cookie, err := req.Cookie(key); err == nil {
	// 	value := ""
	// 	if err = c.s.Decode(key, cookie.Value, &value); err == nil {
	// 		return value, nil
	// 	} else {
	// 		log.WithFields(log.Fields{
	// 			"err": err.Error(),
	// 		}).Error("Cookie Get失败")
	// 		return "", fmt.Errorf("Cookie Get失败")
	// 	}
	// } else {
	// 	log.WithFields(log.Fields{
	// 		"err": err.Error(),
	// 	}).Error("Cookie Get失败")
	// 	return "", fmt.Errorf("Cookie Get失败")
	// }
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
	}

	type Request struct {
		Token string `json:"token"`
	}

	request := &Request{}
	if err := ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
	}
	return request.Token, nil
}

func (c *cookieMgr) Remove(key string, res http.ResponseWriter) {
	// 设置为空字符串即认为删除
	c.Set(key, "", res)
}
