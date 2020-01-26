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
