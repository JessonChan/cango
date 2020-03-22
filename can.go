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

// Opts 程序启动的配置参数
type Opts struct {
	// 监听的主机
	Host string
	// 监听的端口
	Port int
	// 模板文件文件夹，相对程序运行路径
	TplDir string
	// 静态文件文件夹，相对程序运行路径
	StaticDir string
	// 模板文件后续名，默认为 .tpl 和 .html
	TplSuffix []string
	// 是否调试页面，true 表示每次都重新加载模板
	DebugTpl bool
	// 日志文件位置，绝对路径
	CanlogPath string
	// gorilla cookie store 的key
	CookieSessionKey string
	// gorilla cookie store 加密使用的key
	CookieSessionSecure string
}

var defaultTplSuffix = []string{".tpl", ".html"}

var defaultOpts = Opts{
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

func (can *Can) Run(as ...interface{}) {
	// 优先从配置文件中读取
	// 之后从传入参数中读取
	Envs(&defaultAddr)
	Envs(&defaultOpts)

	if loggerInitialed == false && defaultOpts.CanlogPath != "" && defaultOpts.CanlogPath != "console" {
		InitLogger(canlog.NewFileWriter(defaultOpts.CanlogPath))
	} else {
		_, _ = canlog.GetLogger().Writer().Write(cangoMark)
	}

	can.rootPath = getRootPath()
	addr := getAddr(as)
	can.srv.Addr = addr.String()
	can.srv.Handler = can
	can.srv.ErrorLog = canlog.GetLogger()
	opts := getOpts(as)
	can.tplRootPath = filepath.Clean(can.rootPath + "/" + opts.TplDir)
	can.staticRootPath = filepath.Clean(can.rootPath + "/" + opts.StaticDir)
	can.tplSuffix = opts.TplSuffix
	can.debugTpl = opts.DebugTpl
	can.buildStaticRoute()
	can.buildRoute()
	// 务必要先构建route再去构建filter
	can.buildFilter()

	// 初化session
	if opts.CookieSessionKey != "" && gorillaStore == nil {
		gorillaStore = newCookieSession(opts.CookieSessionKey, opts.CookieSessionSecure)
	}

	// 当没有开启session功能时，需要进行空赋值
	if gorillaStore == nil {
		gorillaStore = &emptyGorillaStore{}
	}

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

func getOpts(as []interface{}) *Opts {
	optsPtr := &defaultOpts
	for _, v := range as {
		if opts, ok := v.(Opts); ok {
			if opts.TplDir != "" {
				optsPtr.TplDir = opts.TplDir
			}
			if opts.StaticDir != "" {
				optsPtr.StaticDir = opts.StaticDir
			}
			if len(opts.TplSuffix) != 0 {
				optsPtr.TplSuffix = opts.TplSuffix
			}
			optsPtr.DebugTpl = opts.DebugTpl
		}
	}
	return optsPtr
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
	request := &WebRequest{
		ResponseWriter: rw,
		Request:        r,
	}

	// todo filter should not be here?????
	needStop := false
	var filterReturn interface{}
	var needHandle = true
	for typ, dsp := range can.filterMuxMap {
		match := deepMatch(dsp, r)
		if match.Error() == nil {
			ri := can.filterMap[typ].PreHandle(request)
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
				if filterReturn != nil {
					// 不需要判断 r (*http.Request) 因为他的改变会在函数内生效（指针）
					if reflect.TypeOf(filterReturn).Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
						rw = filterReturn.(http.ResponseWriter)
						request.ResponseWriter = rw
					} else {
						needHandle = false
					}
				}
			}
		}
	}
	if needStop {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	var handleReturn interface{}
	var statusCode int

	if needHandle {
		handleReturn, statusCode = can.serve(request)
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
		http.Redirect(rw, r, handleReturn.(RedirectWithCode).Url, code)
	case Content:
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(handleReturn.(Content).String))
		if err != nil {
			canlog.CanError(err)
		}
	case ContentWithCode:
		code := handleReturn.(ContentWithCode).Code
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
		// todo 可能有安全隐患
		paths := [...]string{can.rootPath + path, can.staticRootPath + path, path}
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
			errorHandleMap[404](rw, r)
		} else {
			http.ServeFile(rw, r, path)
		}
	case DoNothing:
		if statusCode == http.StatusNotFound {
			errorHandleMap[404](rw, r)
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
		match := deepMatch(dsp, r)
		if match.Error() == nil {
			can.filterMap[typ].PostHandle(request)
		}
	}
}

func deepMatch(mux dispatcher, req *http.Request) matcher {
	match := mux.Match(req)
	if match.Error() != nil {
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

func (can *Can) serve(request *WebRequest) (interface{}, int) {
	req := request.Request
	match := deepMatch(can.routeMux, req)
	if match.Error() != nil {
		canlog.CanError(req.Method, req.URL.Path, match.Error())
		return nil, http.StatusNotFound
	}

	invoker, ok := can.methodMap[match.Forwarder().GetName()]
	if ok == false {
		// error,match failed
		return nil, http.StatusMethodNotAllowed
	}

	uriRequestValue := reflect.ValueOf(newContext(request))
	callerIn := make([]reflect.Value, invoker.Type.NumIn())
	cookies := req.Cookies()
	_ = req.ParseForm()
	gs, _ := gorillaStore.Get(request.Request, cangoSessionKey)
	for i := 0; i < len(callerIn); i++ {
		in := invoker.Type.In(i)
		callerIn[i] = newValue(in)
		if in.Implements(uriType) {
			if in == uriType {
				callerIn[i].Set(uriRequestValue)
			} else {
				uriFiled := value(callerIn[i]).FieldByName(uriName)
				if uriFiled.IsValid() && uriFiled.CanSet() {
					uriFiled.Set(uriRequestValue)
				}
			}
		}
		// 先解析form
		// 再赋值path value，如果form中包含和path中相同的变量，被path覆盖
		// 读取cookie，只赋值有cookie标签的变量
		// 解析session，赋值有session标签的变量
		reqHolder := func(key string) (interface{}, int, bool) {
			holderKey := key[0:holderLen]
			valueKey := key[holderLen:]
			switch holderKey {
			case cookieHolderKey:
				for _, cookie := range cookies {
					if cookie.Name == valueKey {
						return cookie.Value, stringFlag, true
					}
				}
				return nil, stringFlag, false
			case sessionHolderKey:
				if i, ok := gs.Values[valueKey]; ok {
					return i, gobBytes, true
				}
				return nil, gobBytes, false
			default:
				if v, ok := match.GetVars()[valueKey]; ok {
					return v, stringFlag, true
				}
				if v, ok := req.Form[valueKey]; ok {
					return v, strSliceFlag, true
				}
				return nil, stringFlag, false
			}
		}

		doDecode(addr(callerIn[i]), reqHolder, fieldTagNames)

		if in.Implements(constructorType) {
			uriFiled := value(callerIn[i]).FieldByName(constructorTypeName)
			if uriFiled.IsValid() && uriFiled.CanSet() {
				uriFiled.Set(valueOfEmptyConstructor)
			}
			addr(callerIn[i]).Interface().(Constructor).Construct(request)
		}
	}
	return call(*invoker.Method, callerIn)
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
	// todo 是否支持无出参的函数
	vs := m.Func.Call(values)
	if len(vs) == 0 {
		return nil, http.StatusMethodNotAllowed
	}
	if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
		return vs[0].Elem().Interface(), http.StatusOK
	}
	return vs[0].Interface(), http.StatusOK
}

var cangoMark = []byte{10, 10, 10, 32, 32, 95, 95, 95, 95, 95, 95, 32, 32, 32, 32, 32, 95, 95, 95, 32, 32, 32, 32, 32, 32, 46, 95, 95, 32, 32, 32, 95, 95, 46, 32, 32, 32, 95, 95, 95, 95, 95, 95, 95, 32, 32, 32, 95, 95, 95, 95, 95, 95, 32, 32, 32, 32, 95, 95, 32, 32, 10, 32, 47, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 47, 32, 32, 32, 92, 32, 32, 32, 32, 32, 124, 32, 32, 92, 32, 124, 32, 32, 124, 32, 32, 47, 32, 32, 95, 95, 95, 95, 95, 124, 32, 47, 32, 32, 95, 95, 32, 32, 92, 32, 32, 124, 32, 32, 124, 32, 10, 124, 32, 32, 44, 45, 45, 45, 45, 39, 32, 32, 47, 32, 32, 94, 32, 32, 92, 32, 32, 32, 32, 124, 32, 32, 32, 92, 124, 32, 32, 124, 32, 124, 32, 32, 124, 32, 32, 95, 95, 32, 32, 124, 32, 32, 124, 32, 32, 124, 32, 32, 124, 32, 124, 32, 32, 124, 32, 10, 124, 32, 32, 124, 32, 32, 32, 32, 32, 32, 47, 32, 32, 47, 95, 92, 32, 32, 92, 32, 32, 32, 124, 32, 32, 46, 32, 96, 32, 32, 124, 32, 124, 32, 32, 124, 32, 124, 95, 32, 124, 32, 124, 32, 32, 124, 32, 32, 124, 32, 32, 124, 32, 124, 32, 32, 124, 32, 10, 124, 32, 32, 96, 45, 45, 45, 45, 46, 47, 32, 32, 95, 95, 95, 95, 95, 32, 32, 92, 32, 32, 124, 32, 32, 124, 92, 32, 32, 32, 124, 32, 124, 32, 32, 124, 95, 95, 124, 32, 124, 32, 124, 32, 32, 96, 45, 45, 39, 32, 32, 124, 32, 124, 95, 95, 124, 32, 10, 32, 92, 95, 95, 95, 95, 95, 95, 47, 95, 95, 47, 32, 32, 32, 32, 32, 92, 95, 95, 92, 32, 124, 95, 95, 124, 32, 92, 95, 95, 124, 32, 32, 92, 95, 95, 95, 95, 95, 95, 124, 32, 32, 92, 95, 95, 95, 95, 95, 95, 47, 32, 32, 40, 95, 95, 41, 32, 10, 10, 10}
