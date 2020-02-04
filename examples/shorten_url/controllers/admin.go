// Copyright 2020 Cango Author.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	model := map[string]interface{}{
		"Name":  params.UrlName,
		"Olink": params.UrlLink,
	}
	if su == nil {
		model["Message"] = "人品太差"
		return cango.ModelView{Tpl: "/views/fail.tpl", Model: model}
	}
	model["Slink"] = su.UniqueId
	return cango.ModelView{Tpl: "/views/success.tpl", Model: model}
}
