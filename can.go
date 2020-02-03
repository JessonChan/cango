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
package cango

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gorilla/schema"
)

type Can struct {
	srv                 *http.Server
	tplRootPath         string
	rootRouter          CanMux
	methodMap           map[string]reflect.Method
	filterMap           map[string][]Filter
	ctrlMap             map[string]ctrlEntry
	staticRequestPrefix string
	tplFuncMap          map[string]interface{}
}

var defaultAddr = Addr{Host: "", Port: 8080}

func NewCan() *Can {
	return &Can{
		srv:        &http.Server{Addr: defaultAddr.String()},
		filterMap:  map[string][]Filter{},
		methodMap:  map[string]reflect.Method{},
		rootRouter: newGorillaMux(),
		ctrlMap:    map[string]ctrlEntry{},
		tplFuncMap: map[string]interface{}{},
	}
}

type Addr struct {
	Host string
	Port int
}

type Opts struct {
	RootPath string
}

type StaticOpts struct {
	RequestPrefix string
}

func (addr Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

type responseTypeHandler func(interface{}) ([]byte, error)

var responseHandler responseTypeHandler = json.Marshal

func (can *Can) SetMux(mux CanMux) {
	can.rootRouter = mux
}

func (can *Can) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, can.staticRequestPrefix) {
		http.ServeFile(rw, r, can.tplRootPath+r.URL.Path)
		return
	}
	rt, statusCode := can.serve(r)
	if rt == nil {
		rw.WriteHeader(int(statusCode))
		_, _ = rw.Write([]byte(nil))
		return
	}
	switch rt.(type) {
	case ModelView:
		mv := rt.(ModelView)
		tpl := can.lookupTpl(mv.Tpl)
		if tpl != nil {
			err := tpl.Execute(rw, mv.Model)
			if err != nil {
				canError("template error", err)
			}
		}
		canError("template not find", mv.Tpl, mv.Model)
		return
	case Redirect:
		http.Redirect(rw, r, rt.(Redirect).Url, rt.(Redirect).Code)
		return
	}

	rw.WriteHeader(int(statusCode))
	bs, err := responseHandler(rt)
	if err == nil {
		_, _ = rw.Write(bs)
	} else {
		_, _ = rw.Write([]byte("{}"))
	}
}

func (can *Can) Run(as ...interface{}) {
	addr := getAddr(as)
	if can.srv == nil {
		can.srv = &http.Server{}
	}
	can.srv.Addr = addr.String()
	can.srv.Handler = can
	can.tplRootPath = getRootPath(as)
	can.staticRequestPrefix = getStaticRequestPrefix(as)
	can.buildRoute()

	startChan := make(chan error, 1)
	go func() {
		canInfo("cango start success @ " + addr.String())
		err := can.srv.ListenAndServe()
		if err != nil {
			startChan <- err
		}
	}()
	err := <-startChan
	if err != nil {
		panic(err)
	}
}

func getAddr(as []interface{}) Addr {
	for _, v := range as {
		if addr, ok := v.(Addr); ok {
			return addr
		}
	}
	return defaultAddr
}

func getRootPath(as []interface{}) string {
	for _, v := range as {
		if opts, ok := v.(Opts); ok {
			return opts.RootPath
		}
	}
	abs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return os.Args[0]
	}
	return filepath.Dir(abs)
}

func getStaticRequestPrefix(as []interface{}) string {
	for _, v := range as {
		if opts, ok := v.(StaticOpts); ok {
			opts.RequestPrefix = filepath.Clean(opts.RequestPrefix)
			if strings.HasPrefix(opts.RequestPrefix, "/") {
				return opts.RequestPrefix
			}
			return "/" + opts.RequestPrefix
		}
	}
	return "/static"
}

type StatusCode int

var decoder = schema.NewDecoder()

func (can *Can) serve(req *http.Request) (interface{}, StatusCode) {
	match := can.rootRouter.Match(req)
	if match.Error() != nil {
		return nil, http.StatusNotFound
	}
	fs, _ := can.filterMap[match.Route().GetName()]
	if len(fs) > 0 {
		for _, f := range fs {
			f.PreHandle(req)
		}
		// do filter
	}

	m, ok := can.methodMap[match.Route().GetName()]
	if ok == false {
		// error,match failed
		return nil, http.StatusMethodNotAllowed
	}
	if m.Type.NumIn() != 2 {
		// error,method only have one arg
		return nil, http.StatusMethodNotAllowed
	}
	// controller
	ct := reflect.New(m.Type.In(0).Elem())
	// method
	mt := reflect.New(m.Type.In(1)).Elem()

	// todo 是否做如下区分 get=>Form, post/put/patch=>PostForm
	// todo 是否需要在此类方法上支持更多的特性，如自定义struct来区分pathValue和formValue
	// todo 性能提升
	if func(methods []string) bool {
		for _, m := range methods {
			switch m {
			case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch:
				return true
			default:
				return false
			}
		}
		return false
	}(match.Route().GetMethods()) {
		_ = req.ParseForm()
		if len(req.Form) > 0 {
			_ = decoder.Decode(mt.Addr().Interface(), req.Form)
		}
	}
	if len(match.GetVars()) > 0 {
		_ = decoder.Decode(ct.Interface(), match.GetVars())
		_ = decoder.Decode(mt.Addr().Interface(), match.GetVars())
	}

	vs := ct.MethodByName(m.Name).Call([]reflect.Value{mt})
	if len(vs) == 0 {
		return nil, http.StatusMethodNotAllowed
	}
	if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
		return vs[0].Elem().Interface(), http.StatusOK
	}
	return vs[0].Interface(), http.StatusOK
}
func toValues(m map[string]string) map[string][]string {
	mm := make(map[string][]string, len(m))
	for k, v := range m {
		mm[k] = make([]string, 1)
		mm[k][0] = v
	}
	return mm
}
