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

type ShorterController struct {
	cango.URI
}

func (s *ShorterController) Redirect(params struct {
	cango.URI `value:"/t/{Id}"`
	Id        string
}) interface{} {
	url := models.GetUrl(params.Id)
	if url == "" {
		url = "http://www.github.com"
	}
	return cango.Redirect{Code: 302, Url: url}
}
