package main

import (
	"log"
	"net/http"
	"os"

	// "github.com/JessonChan/canlog"

	"github.com/JessonChan/cango"
)

// CanCtrl AND  VisitFilter

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

type Page struct {
	Offset int
	Size   int
}

func (c *CanCtrl) Api(ps struct {
	cango.URI `value:"/{name}/{year}/api.json"`
	cango.PostMethod
	Name  string
	Year  int
	Color string
	Page
}) interface{} {
	return map[string]interface{}{
		"name":  ps.Name,
		"age":   ps.Year,
		"color": ps.Color,
		"size":  ps.Size,
	}
}

type VisitFilter struct {
	cango.Filter
	// controller
	*CanCtrl
}

var _ = cango.RegisterFilter(&VisitFilter{})

func (v *VisitFilter) PreHandle(req *http.Request) interface{} {
	// canlog.CanDebug(req.Method, req.URL.Path)
	log.Println(req.Method, req.URL.Path)
	return true
}

func main() {
	// better use canlog
	// canlog.SetWriter(canlog.NewFileWriter("/tmp/cango-app.log"), "App")
	// cango.InitLogger(canlog.GetLogger().Writer())
	cango.InitLogger(os.Stdout)
	cango.
		NewCan().
		RegTplFunc("buttonValue", func() string {
			return "click me!"
		}).
		Run(cango.Addr{Port: 8008}, cango.StaticOpts{TplSuffix: []string{".html"}, Debug: true})
}
