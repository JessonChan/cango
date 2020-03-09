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
	methodMap    map[string]*invoker
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
		methodMap:    map[string]*invoker{},
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
	Host       string
	Port       int
	RootPath   string
	TplDir     string
	StaticDir  string
	TplSuffix  []string
	DebugTpl   bool
	CanlogPath string
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

var loggerInitialed = false

// InitLogger 用来初始化cango的日志writer
func InitLogger(rw io.Writer) {
	loggerInitialed = true
	canlog.SetWriter(rw, "CANGO")
}

func (can *Can) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			var buf [1024 * 10]byte
			runtime.Stack(buf[:], false)
			canlog.CanError(err)
			canlog.CanError(string(buf[0:]))
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	// todo filter should not be here?????
	needStop := false
	var filterReturn interface{}
	for typ, dsp := range can.filterMuxMap {
		match := doubleMatch(dsp, r)
		if match.Error() == nil {
			ri := can.filterMap[typ].PreHandle(rw, r)
			if rt, ok := ri.(bool); ok {
				// todo 这样的设计是不是合理？？？？
				// 返回为false 这个之后注册的filter失效
				if rt == false {
					needStop = true
				}
			} else {
				// 如果不是bool类型，需要提前结束
				// todo 需要设计filter的执行顺序，优先生效最早返回的
				filterReturn = ri
			}
		}
	}
	if needStop {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	var handleReturn interface{}
	var statusCode int
	var needHandle = true
	if filterReturn != nil {
		// 不需要判断 r (*http.Request) 因为他的改变会在函数内生效（指针）
		if reflect.TypeOf(filterReturn).Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
			rw = filterReturn.(http.ResponseWriter)
		} else {
			needHandle = false
		}
	}
	if needHandle {
		handleReturn, statusCode = can.serve(rw, r)
		// todo nil是不是可以表示已经在函数内完成了？
		// todo 这样的设计返回是有问题的
		if handleReturn == nil {
			handleReturn = doNothing
		}
	} else {
		handleReturn = filterReturn
	}

	// todo 直接在type switch中判断
	if reflect.TypeOf(handleReturn).Kind() == reflect.Ptr {
		handleReturn = reflect.ValueOf(handleReturn).Elem().Interface()
	}

	switch handleReturn.(type) {
	case ModelView:
		mv := handleReturn.(ModelView)
		tpl := can.lookupTpl(mv.Tpl)
		if tpl == nil {
			canlog.CanError("template not find", mv.Tpl, mv.Model)
			return
		} else {
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := tpl.Execute(rw, mv.Model)
			if err != nil {
				canlog.CanError("template error", err)
			}
		}
	case Redirect:
		http.Redirect(rw, r, handleReturn.(Redirect).Url, http.StatusFound)
	case RedirectWithCode:
		code := handleReturn.(RedirectWithCode).Code
		if code == 0 {
			code = http.StatusFound
		}
		http.Redirect(rw, r, handleReturn.(RedirectWithCode).Url, http.StatusFound)
	case Content:
		code := handleReturn.(Content).Code
		if code == 0 {
			code = http.StatusOK
		}
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(code)
		_, err := rw.Write([]byte(handleReturn.(Content).String))
		if err != nil {
			canlog.CanError(err)
		}
	case StaticFile:
		var err error
		path := handleReturn.(StaticFile).Path
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
			canlog.CanDebug(err, "can't find the file", handleReturn)
			error404(rw, r)
		} else {
			http.ServeFile(rw, r, path)
		}
	case DoNothing:
		if statusCode == http.StatusNotFound {
			error404(rw, r)
		}
	default:
		rw.WriteHeader(statusCode)
		bs, err := responseJsonHandler(handleReturn)
		if err == nil {
			_, _ = rw.Write(bs)
		} else {
			_, _ = rw.Write([]byte("{}"))
		}
	}
	// postHandle
	for typ, dsp := range can.filterMuxMap {
		match := doubleMatch(dsp, r)
		if match.Error() == nil {
			can.filterMap[typ].PostHandle(rw, r)
		}
	}
}

func (can *Can) Run(as ...interface{}) {
	// 优先从配置文件中读取
	// 之后从传入参数中读取
	Envs(&defaultAddr)
	Envs(&defaultOpts)

	if loggerInitialed == false && defaultOpts.CanlogPath != "" && defaultOpts.CanlogPath != "console" {
		InitLogger(canlog.NewFileWriter(defaultOpts.CanlogPath))
	}

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
	host := defaultAddr.Host
	port := defaultAddr.Port
	for _, v := range as {
		if addr, ok := v.(Addr); ok {
			return addr
		}
		if h, ok := v.(string); ok {
			host = h
		}
		if p, ok := v.(int); ok {
			port = p
		}
	}
	return Addr{host, port}
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
		TplSuffix: append(defaultOpts.TplSuffix),
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

func (can *Can) serve(rw http.ResponseWriter, req *http.Request) (interface{}, int) {
	match := doubleMatch(can.routeMux, req)
	if match.Error() != nil {
		canlog.CanError(req.Method, req.URL.Path, match.Error())
		return nil, http.StatusNotFound
	}

	invoker, ok := can.methodMap[match.Forwarder().GetName()]
	if ok == false {
		// error,match failed
		return nil, http.StatusMethodNotAllowed
	}

	uriContext := reflect.ValueOf(newContext(rw, req))
	callerIn := make([]reflect.Value, invoker.Type.NumIn())
	var receiver reflect.Value
	if invoker.kind == invokeByReceiver {
		receiver = newValue(invoker.Type.In(0))
		uriFiled := receiver.Elem().FieldByName(uriName)
		if uriFiled.IsValid() && uriFiled.CanSet() {
			uriFiled.Set(uriContext)
		}
		callerIn[0] = receiver
	}
	foundURI := -1
	for j := len(callerIn); j > invoker.kind; j-- {
		idx := j - 1
		in := invoker.Type.In(idx)
		callerIn[idx] = newValue(in)
		// 只会对最后一个进行赋值
		if in.Implements(uriType) && foundURI == -1 {
			foundURI = idx
			if in == uriType {
				callerIn[idx].Set(uriContext)
			} else {
				uriFiled := callerIn[idx].FieldByName(uriName)
				if uriFiled.IsValid() && uriFiled.CanSet() {
					uriFiled.Set(uriContext)
				}
			}
		}
		if in.Implements(constructorType) {
			// todo cookie and form 同时成立 ？？？
			if in.Implements(cookieType) {
				value(callerIn[idx]).FieldByName(cookieTypeName).Set(valueOfEmptyCookie)
				cookies := req.Cookies()
				checkSet(stringFlag, cookieHolder(cookies), addr(callerIn[idx]), cookieNameWithTag)
				addr(callerIn[idx]).Interface().(CookieValue).Construct(req)
			} else if in.Implements(formValueType) {
				value(callerIn[idx]).FieldByName(formValueTypeName).Set(valueOfEmptyForm)
				// ParsForm可以多少调用，不会影响性能或副作用
				_ = req.ParseForm()
				decodeForm(req.Form, addr(callerIn[idx]), pathFormFn)
				addr(callerIn[idx]).Interface().(FormValue).Construct(req)
			} else if in.Implements(pathValueType) {
				value(callerIn[idx]).FieldByName(pathValueTypeName).Set(valueOfEmptyPath)
				decode(match.GetVars(), addr(callerIn[idx]), pathFormFn)
				addr(callerIn[idx]).Interface().(PathValue).Construct(req)
			} else {
				// todo ???
				value(callerIn[idx]).FieldByName(constructorTypeName).Set(valueOfEmptyConstructor)
				addr(callerIn[idx]).Interface().(FormValue).Construct(req)
			}
		}
	}
	var args0 = callerIn[foundURI]
	// 先解析form
	if shouldParseForm(req.Method) {
		_ = req.ParseForm()
		if len(req.Form) > 0 {
			if invoker.kind == invokeByReceiver {
				decodeForm(req.Form, addr(receiver))
			}
			if args0.Type() != uriType {
				decodeForm(req.Form, addr(args0), pathFormFn)
			}
		}
	}

	// 再赋值path value，如果form中包含和path中相同的变量，被path覆盖
	if len(match.GetVars()) > 0 {
		if invoker.kind == invokeByReceiver {
			decode(match.GetVars(), addr(receiver), pathFormFn)
		}
		if args0.Type() != uriType {
			decode(match.GetVars(), addr(args0), pathFormFn)
		}
	}

	// 最后读取cookie，只赋值有cookie标签的变量
	cookies := req.Cookies()
	if len(cookies) >= 0 {
		if invoker.kind == invokeByReceiver {
			checkSet(stringFlag, cookieHolder(cookies), addr(receiver), cookieNameWithTag)
		}
		if args0.Type() != uriType {
			checkSet(stringFlag, cookieHolder(cookies), addr(args0), cookieNameWithTag)
		}
	}
	return call(*invoker.Method, callerIn)
}

func cookieHolder(cookies []*http.Cookie) func(cookieName string) (interface{}, bool) {
	return func(cookieName string) (interface{}, bool) {
		for _, cookie := range cookies {
			if cookie.Name == cookieName {
				return cookie.Value, true
			}
		}
		return nil, false
	}
}
func pathFormFn(field reflect.StructField) []string {
	return filedName(field, pathFormName)
}

func value(value reflect.Value) reflect.Value {
	if value.Type().Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

func addr(value reflect.Value) reflect.Value {
	if value.Type().Kind() == reflect.Ptr {
		return value
	}
	return value.Addr()
}

func newValue(typ reflect.Type) reflect.Value {
	// pointer,new typ.Elem()
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	// normal
	return reflect.New(typ).Elem()
}

// 执行函数
func call(m reflect.Method, values []reflect.Value) (interface{}, int) {
	vs := m.Func.Call(values)
	if len(vs) == 0 {
		return nil, http.StatusMethodNotAllowed
	}
	if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
		return vs[0].Elem().Interface(), http.StatusOK
	}
	return vs[0].Interface(), http.StatusOK
}

// todo 是否做如下区分 get=>Form, post/put/patch=>PostForm
// todo 是否需要在此类方法上支持更多的特性，如自定义struct来区分pathValue和formValue
// todo 性能提升
func shouldParseForm(httpMethod string) bool {
	switch httpMethod {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}
