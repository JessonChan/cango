package main

import (
	"fmt"

	"github.com/JessonChan/cango"
	"github.com/JessonChan/cango/cango-cli/shorten_url/controllers"
)

func main() {
	port := 8088
	cango.NewCan().
		RegTplFunc("localHost", func() string { return fmt.Sprintf("127.0.0.1:%d", port) }).
		Route(&controllers.AdminController{}, &controllers.ShorterController{}).
		// ListCtrl 已经在代码中注册
		Run(cango.Addr{Port: port})
}
