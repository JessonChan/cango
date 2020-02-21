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

type ListCtrl struct {
	cango.URI `value:"/;/list"`
}

var _ = cango.RegisterURI(&ListCtrl{})

func (a *ListCtrl) Get(params struct {
	cango.URI
}) interface{} {
	return cango.ModelView{Tpl: "/views/list.tpl", Model: map[string]interface{}{"Slice": models.GetAll()}}
}
