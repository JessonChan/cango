package cango

import (
	"fmt"
	"net/http"
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

func (can *Can) Route(string, *Controller) {
}
func (can *Can) route([]string, *Controller) {
}
