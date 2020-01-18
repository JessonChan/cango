package main

import (
	"github.com/JessonChan/cango"
)

type PageController struct {
	cango.URI `value:"/blog/{blogName}/article/{articleId}-{pageId}.json"`
	name      string
}

func main() {
	can := cango.Can{}
	can.Route(&PageController{})
	// can.Run(cango.Addr{Port: 8081})
}
