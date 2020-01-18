package cango

import (
	"fmt"
	"net/http"
	"reflect"
)

type Can struct {
	srv *http.Server
}

var defaultAddr = Addr{Host: "", Port: 8080}

func NewCan() *Can {
	return &Can{&http.Server{Addr: defaultAddr.String()}}
}

type Addr struct {
	Host string
	Port int
}

func (addr Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

type canHandlers struct {
}

func (ch *canHandlers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
}

func (can *Can) Run(as ...Addr) {
	addr := append(as, defaultAddr)[0]
	if can.srv == nil {
		can.srv = &http.Server{Addr: addr.String()}
	}
	can.srv.Handler = http.DefaultServeMux
	err := can.srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (can *Can) Route(uri URI) {
	rv := reflect.ValueOf(uri)
	if rv.Kind() != reflect.Ptr {
		panic(uri)
	}
	tv := reflect.Indirect(rv).Interface()
	fmt.Println(reflect.TypeOf(tv))
	tvs := reflect.TypeOf(tv)
	for i := 0; i < tvs.NumField(); i++ {
		f := tvs.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if fmt.Sprint(f.Type) == "cango.URI" {
			fmt.Println(f.Name + ":" + f.Type.Name() + ":" + f.Tag.Get("value"))
			fmt.Println(parseTag(f.Tag.Get("value")))
		}
	}
}
func parseTag(tag string) (vars map[int]string) {
	if tag == "" {
		return
	}
	vars = map[int]string{}
	segNum := 0
	lastBit := tag[0]
	segBegin := 0
	for i := 0; i < len(tag); i++ {
		v := tag[i]
		if lastBit == '{' {
			segNum++
			segBegin = i
		}
		if v == '}' {
			vars[segNum] = tag[segBegin:i]
		}
		lastBit = v
	}
	return
}
func (can *Can) route([]string, *Controller) {
}
