package controller

import (
	"github.com/JessonChan/cango"
)

type CanCtrl struct {
	cango.URI `value:"/"`
}

var _ = cango.RegisterURI(&CanCtrl{})

func (c *CanCtrl) Index(ps struct {
	cango.URI `value:"/;/index/index.html"`
}) interface{} {
	return cango.ModelView{
		Tpl:   "index.html",
		Model: "Hello,cango!",
	}
}

func (c *CanCtrl) Api(ps struct {
	cango.URI `value:"/api;api.json"`
}) interface{} {
	return map[string]string{
		"hello": "cango",
	}
}
