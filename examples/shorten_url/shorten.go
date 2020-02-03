package main

import (
	"fmt"

	"github.com/JessonChan/cango"
	"github.com/JessonChan/cango/examples/shorten_url/controllers"
)

func main() {
	port := 8088
	cango.NewCan().
		RegTplFunc("localHost", func() string { return fmt.Sprintf("127.0.0.1:%d", port) }).
		Route(&controllers.AdminController{}, &controllers.ShorterController{}).
		Run(cango.Addr{Port: port})
}
