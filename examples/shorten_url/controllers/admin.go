package controllers

import (
	"github.com/JessonChan/cango"
	"github.com/JessonChan/cango/examples/shorten_url/models"
)

type AdminController struct {
	cango.URI `value:"/dwz"`
}

func (a *AdminController) Get(params struct {
	cango.URI `value:"/admin_____login"`
}) interface{} {
	return cango.ModelView{Tpl: "/views/admin.tpl", Model: map[string]interface{}{"Slice": models.GetAll()}}
}

func (a *AdminController) Form(params struct {
	cango.URI
}) interface{} {
	return cango.ModelView{Tpl: "/views/index.tpl"}
}
func (a *AdminController) Submit(params struct {
	cango.URI `value:"/{submitId}"`
	cango.PostMethod
	UrlLink  string
	UrlName  string
	SubmitId int
	HelpMsg  string
}) interface{} {
	su := models.Insert(params.UrlLink, params.UrlName)
	var model interface{}
	if su != nil {
		model = map[string]interface{}{
			"Name":  su.Name,
			"Olink": su.Url,
			"Slink": su.UniqueId,
		}
	}
	return cango.ModelView{Tpl: "/views/success.tpl", Model: model}
}
