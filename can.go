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
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/JessonChan/jsun"
	"github.com/gorilla/schema"
)

type Can struct {
	srv                 *http.Server
	rootPath            string
	tplRootPath         string
	staticRootPath      string
	staticRequestPrefix string
	tplSuffix           []string
	debugTpl            bool

	rootRouter CanMux
	methodMap  map[string]reflect.Method
	filterMap  map[string][]Filter
	ctrlMap    map[string]ctrlEntry
	tplFuncMap map[string]interface{}
	tplNameMap map[string]bool
}

var defaultAddr = Addr{Host: "", Port: 8080}

func NewCan() *Can {
	return &Can{
		srv:        &http.Server{Addr: defaultAddr.String()},
		rootRouter: newGorillaMux(),
		methodMap:  map[string]reflect.Method{},
		filterMap:  map[string][]Filter{},
		ctrlMap:    map[string]ctrlEntry{},
		tplFuncMap: map[string]interface{}{},
		tplNameMap: map[string]bool{},
	}
}

const (
	tplDir    = "/views"
	staticDir = "/static"
)

type Addr struct {
	Host string
	Port int
}

type Opts struct {
	RootPath string
}

type StaticOpts struct {
	RequestPrefix string
	TplSuffix     []string
	Debug         bool
}

func (addr Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

type responseTypeHandler func(interface{}) ([]byte, error)

var responseJsonHandler responseTypeHandler = func(v interface{}) (bytes []byte, err error) { return jsun.Marshal(v, jsun.LowerCamelStyle) }

func (can *Can) SetMux(mux CanMux) {
	can.rootRouter = mux
}

func (can *Can) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rt, statusCode := can.serve(rw, r)
	if rt == nil {
		rw.WriteHeader(int(statusCode))
		_, _ = rw.Write([]byte(nil))
		return
	}
	switch rt.(type) {
	case ModelView:
		mv := rt.(ModelView)
		tpl := can.lookupTpl(mv.Tpl)
		if tpl == nil {
			canError("template not find", mv.Tpl, mv.Model)
			return
		}
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tpl.Execute(rw, mv.Model)
		if err != nil {
			canError("template error", err)
		}
		return
	case Redirect:
		code := rt.(Redirect).Code
		if code == 0 {
			code = http.StatusFound
		}
		http.Redirect(rw, r, rt.(Redirect).Url, code)
		return
	case Content:
		code := rt.(Content).Code
		if code == 0 {
			code = http.StatusOK
		}
		rw.WriteHeader(code)
		_, _ = rw.Write([]byte(rt.(Content).String))
		return
	case StaticFile:
		http.ServeFile(rw, r, rt.(StaticFile).Path)
		return
	}

	rw.WriteHeader(int(statusCode))
	bs, err := responseJsonHandler(rt)
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
	can.rootPath = getRootPath(as)
	can.tplRootPath = can.rootPath + tplDir
	can.staticRootPath = can.rootPath + staticDir
	staticOpts := getStaticOpts(as)
	can.staticRequestPrefix = staticOpts.RequestPrefix
	can.tplSuffix = staticOpts.TplSuffix
	can.debugTpl = staticOpts.Debug
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
	dir, err := os.Getwd()
	if err == nil {
		return dir
	}
	abs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return os.Args[0]
	}
	return filepath.Dir(abs)
}

func getStaticOpts(as []interface{}) StaticOpts {
	for _, v := range as {
		if opts, ok := v.(StaticOpts); ok {
			if opts.RequestPrefix == "" {
				opts.RequestPrefix = "/static"
			}
			opts.RequestPrefix = filepath.Clean(opts.RequestPrefix)
			if strings.HasPrefix(opts.RequestPrefix, "/") == false {
				opts.RequestPrefix = "/" + opts.RequestPrefix
			}
			if len(opts.TplSuffix) == 0 {
				opts.TplSuffix = []string{".tpl"}
			}
			return opts
		}
	}
	return StaticOpts{RequestPrefix: "/static", TplSuffix: []string{".tpl"}}
}

type StatusCode int

var decoder = schema.NewDecoder()

func (can *Can) serve(rw http.ResponseWriter, req *http.Request) (interface{}, StatusCode) {
	// todo static file should use router match also
	if strings.HasPrefix(req.URL.Path, can.staticRequestPrefix) {
		return StaticFile{Path: filepath.Clean(can.rootPath + req.URL.Path)}, http.StatusOK
	}

	match := can.rootRouter.Match(req)
	if match.Error() != nil {
		// todo sure?
		req.URL.Path = filepath.Clean(req.URL.Path)
		match = can.rootRouter.Match(req)
		if match.Error() != nil {
			return nil, http.StatusNotFound
		}
	}
	fs, _ := can.filterMap[match.Route().GetName()]
	if len(fs) > 0 {
		for _, f := range fs {
			ri := f.PreHandle(req)
			if rt, ok := ri.(bool); ok {
				// 返回为false 这个之后注册的filter失效
				if rt == false {
					break
				}
			} else {
				// 如果不是bool类型，提前结束
				return ri, http.StatusOK
			}
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
	uriFiled := ct.Elem().FieldByName(uriName)
	context := newContext(rw, req)
	if uriFiled.IsValid() {
		uriFiled.Set(reflect.ValueOf(context))
	}
	// method
	mt := reflect.New(m.Type.In(1)).Elem()
	uriFiled = mt.FieldByName(uriName)
	if uriFiled.IsValid() {
		uriFiled.Set(reflect.ValueOf(context))
	}

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
