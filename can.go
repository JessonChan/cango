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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/JessonChan/canlog"
	"github.com/JessonChan/jsun"
)

type Can struct {
	srv            *http.Server
	rootPath       string
	tplRootPath    string
	staticRootPath string
	tplSuffix      []string
	debugTpl       bool

	routeMux     dispatcher
	filterMuxMap map[reflect.Type]dispatcher
	methodMap    map[string]reflect.Method
	filterMap    map[reflect.Type]Filter
	ctrlEntryMap map[string]ctrlEntry
	tplFuncMap   map[string]interface{}
	tplNameMap   map[string]bool
}

var defaultAddr = Addr{Host: "", Port: 8080}

func NewCan() *Can {
	return &Can{
		srv:          &http.Server{Addr: defaultAddr.String()},
		routeMux:     newCanMux(),
		filterMuxMap: map[reflect.Type]dispatcher{},
		methodMap:    map[string]reflect.Method{},
		filterMap:    map[reflect.Type]Filter{},
		ctrlEntryMap: map[string]ctrlEntry{},
		tplFuncMap:   map[string]interface{}{},
		tplNameMap:   map[string]bool{},
	}
}

type Addr struct {
	Host string
	Port int
}

type Opts struct {
	RootPath  string
	TplDir    string
	StaticDir string
	TplSuffix []string
	DebugTpl  bool
}

var defaultTplSuffix = []string{".tpl", ".html"}

var defaultOpts = Opts{
	RootPath:  getRootPath(),
	TplDir:    "/view",
	StaticDir: "/static",
	TplSuffix: defaultTplSuffix,
	DebugTpl:  false,
}

func (addr Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

type responseTypeHandler func(interface{}) ([]byte, error)

var responseJsonHandler responseTypeHandler = func(v interface{}) (bytes []byte, err error) { return jsun.Marshal(v, jsun.LowerCamelStyle) }

// InitLogger 用来初始化cango的日志writer
func InitLogger(rw io.Writer) {
	canlog.SetWriter(rw, "CANGO")
}

var uriRegMap = map[URI]bool{}

// todo with prefix???
// todo with can app Name ???
func RegisterURI(uri URI) bool {
	uriRegMap[uri] = true
	return true
}

func (can *Can) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			var buf [1024 * 10]byte
			runtime.Stack(buf[:], false)
			canlog.CanError(string(buf[0:]))
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	rt, statusCode := can.serve(rw, r)
	// todo nil是不是可以表示已经在函数内完成了？
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
			canlog.CanError("template not find", mv.Tpl, mv.Model)
			return
		}
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tpl.Execute(rw, mv.Model)
		if err != nil {
			canlog.CanError("template error", err)
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
		_, err := rw.Write([]byte(rt.(Content).String))
		if err != nil {
			canlog.CanError(err)
		}
		return
	case StaticFile:
		var err error
		path := rt.(StaticFile).Path
		if path[0] != '/' {
			path = "/" + path
		}
		// todo 更好的实现 filepath.Clean的性能问题
		paths := [3]string{can.rootPath + path, path, can.staticRootPath + path}
		for _, p := range paths {
			_, err = os.Stat(p)
			if err != nil {
				continue
			}
			path = p
			break
		}
		if err != nil {
			canlog.CanDebug(err, "can't find the file", rt)
			// todo 404 default page
		}
		http.ServeFile(rw, r, path)
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
	can.srv.Addr = addr.String()
	can.srv.Handler = can
	can.srv.ErrorLog = canlog.GetLogger()
	opts := getOpts(as)
	can.rootPath = opts.RootPath
	can.tplRootPath = filepath.Clean(can.rootPath + "/" + opts.TplDir)
	can.staticRootPath = filepath.Clean(can.rootPath + "/" + opts.StaticDir)
	can.tplSuffix = opts.TplSuffix
	can.debugTpl = opts.DebugTpl
	can.buildStaticRoute()
	can.buildRoute()
	// 务必要先构建route再去构建filter
	can.buildFilter()

	startChan := make(chan error, 1)
	go func() {
		canlog.CanInfo("cango start success @ " + addr.String())
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

func getOpts(as []interface{}) Opts {
	newOpts := copyOpts()
	for _, v := range as {
		if opts, ok := v.(Opts); ok {
			if opts.RootPath != "" {
				newOpts.RootPath = opts.RootPath
			}
			if opts.TplDir != "" {
				newOpts.TplDir = opts.TplDir
			}
			if opts.StaticDir != "" {
				newOpts.StaticDir = opts.StaticDir
			}
			if len(opts.TplSuffix) != 0 {
				newOpts.TplSuffix = opts.TplSuffix
			}
			newOpts.DebugTpl = opts.DebugTpl
		}
	}
	return newOpts
}
func copyOpts() Opts {
	return Opts{
		RootPath:  defaultOpts.RootPath,
		TplDir:    defaultOpts.TplDir,
		StaticDir: defaultOpts.StaticDir,
		TplSuffix: append(defaultTplSuffix),
		DebugTpl:  false,
	}
}
func getRootPath() string {
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

type StatusCode int

func doubleMatch(mux dispatcher, req *http.Request) matcher {
	match := mux.Match(req)
	if match.Error() != nil {
		// todo 这里有性能问题
		originalPath := req.URL.Path
		req.URL.Path = filepath.Clean(originalPath)
		if originalPath == req.URL.Path {
			return match
		}
		match = mux.Match(req)
		req.URL.Path = originalPath
		return match
	}
	return match
}

func (can *Can) serve(rw http.ResponseWriter, req *http.Request) (interface{}, StatusCode) {
	// todo filter should not here?????
	for typ, dsp := range can.filterMuxMap {
		match := doubleMatch(dsp, req)
		if match.Error() == nil {
			ri := can.filterMap[typ].PreHandle(req)
			if rt, ok := ri.(bool); ok {
				// todo 这样的设计是不是合理？？？？
				// 返回为false 这个之后注册的filter失效
				if rt == false {
					break
				}
			} else {
				// 如果不是bool类型，提前结束
				return ri, http.StatusOK
			}
		}
	}

	match := doubleMatch(can.routeMux, req)
	if match.Error() != nil {
		canlog.CanError(req.Method, req.URL.Path, match.Error())
		return nil, http.StatusNotFound
	}
	m, ok := can.methodMap[match.Route().GetName()]
	if ok == false {
		// error,match failed
		return nil, http.StatusMethodNotAllowed
	}
	if m.Type.NumIn() == 1 {

		if m.Type.In(0) == uriType {
			vs := m.Func.Call([]reflect.Value{reflect.ValueOf(newContext(rw, req))})
			if len(vs) == 0 {
				return nil, http.StatusMethodNotAllowed
			}
			if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
				return vs[0].Elem().Interface(), http.StatusOK
			}
			return vs[0].Interface(), http.StatusOK
		}
		if m.Type.In(0).Implements(uriType) {
			mt := reflect.New(m.Type.In(0)).Elem()
			uriFiled := mt.FieldByName(uriName)
			if uriFiled.IsValid() {
				uriFiled.Set(reflect.ValueOf(newContext(rw, req)))
			}
			if func(methods string) bool {
				switch methods {
				case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch:
					return true
				default:
					return false
				}
			}(req.Method) {
				_ = req.ParseForm()
				if len(req.Form) > 0 {
					decodeForm(req.Form, mt.Addr().Interface())
				}
			}
			if len(match.GetVars()) > 0 {
				decode(match.GetVars(), mt.Addr().Interface())
			}
			vs := m.Func.Call([]reflect.Value{mt})
			if len(vs) == 0 {
				return nil, http.StatusMethodNotAllowed
			}
			if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
				return vs[0].Elem().Interface(), http.StatusOK
			}
			return vs[0].Interface(), http.StatusOK
		}
	}

	if m.Type.NumIn() == 1 && m.Type.In(0) == uriType {
		vs := m.Func.Call([]reflect.Value{reflect.ValueOf(newContext(rw, req))})
		if len(vs) == 0 {
			return nil, http.StatusMethodNotAllowed
		}
		if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
			return vs[0].Elem().Interface(), http.StatusOK
		}
		return vs[0].Interface(), http.StatusOK
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
	if func(httpMethod string) bool {
		switch httpMethod {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch:
			return true
		default:
			return false
		}
	}(req.Method) {
		_ = req.ParseForm()
		if len(req.Form) > 0 {
			decodeForm(req.Form, ct.Interface())
			decodeForm(req.Form, mt.Addr().Interface())
		}
	}
	if len(match.GetVars()) > 0 {
		decode(match.GetVars(), ct.Interface())
		decode(match.GetVars(), mt.Addr().Interface())
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
