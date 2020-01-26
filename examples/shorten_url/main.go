package main

import (
	"github.com/JessonChan/cango"
	"github.com/JessonChan/cango/examples/shorten_url/controllers"
)

func main() {
	cango.NewCan().
		Route(&controllers.AdminController{}, &controllers.ShorterController{}).
		Run(cango.Addr{Port: 8088}, cango.View{RootPath: "/Users/jessonchan/code/go/path/src/github.com/JessonChan/cango/examples/shorten_url"})
}
