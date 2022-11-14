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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/JessonChan/canlog"
	"github.com/JessonChan/jsun"
	"github.com/gin-gonic/gin"
)

type any interface{}
type Can struct {
	name           string
	srv            *http.Server
	rootPath       string
	tplRootPath    string
	staticRootPath string
	tplSuffix      []string
	debugTpl       bool

	routeMux         *routeDispatcher
	filterDispatcher map[FilterType]*filterDispatcher
	tplFuncMap       map[string]interface{}
	tplNameMap       map[string]bool
	fallbackHandler  http.Handler
}

var defaultAddr = Addr{Host: "", Port: 8080}

// NewCan 生成服务对象，name 为生成的对象名，可以为空
// 一般，只有我们在工程中需要注册多个web服务时才需要设置
func NewCan(name ...string) *Can {
	return &Can{
		name:             append(name, "")[0],
		srv:              &http.Server{Addr: defaultAddr.String()},
		routeMux:         &routeDispatcher{dispatcher: newCanMux(), ctrlEntryMap: map[string]ctrlEntry{}},
		filterDispatcher: map[FilterType]*filterDispatcher{},
		tplFuncMap:       map[string]interface{}{},
		tplNameMap:       map[string]bool{},
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

var responseJsonHandler responseTypeHandler = func(v interface{}) (bytes []byte, err error) { return jsun.Marshal(v) }

var loggerInitialed = false

// InitLogger 用来初始化cango的日志writer
func InitLogger(rw io.Writer) {
	loggerInitialed = true
	canlog.SetWriter(rw, "CANGO")
}

func (can *Can) FallbackHandler(handler http.Handler) *Can {
	can.fallbackHandler = handler
	return can
}

// SetJsonWriter will set a json writer
func (can *Can) SetJsonWriter(rth responseTypeHandler) *Can {
	responseJsonHandler = rth
	return can
}
func (can *Can) Mode(runMode string) *Can {
	cangoRunMode = runMode
	return can
}

func (can *Can) Run(as ...interface{}) error {
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

	canlog.CanInfo("cango start success @ " + addr.String())
	return can.srv.ListenAndServe()
}

type GinHandler struct {
	Url         string
	HttpMethods map[string]bool
	Handle      func(ctx *gin.Context)
}

func (can *Can) ToGins() []*GinHandler {
	can.buildStaticRoute()
	can.buildRoute()
	return can.routeMux.dispatcher.(*canDispatcher).Gins()
}
func (p *Can) Shutdown() *Can {
	p.srv.Shutdown(context.Background())
	return p
}

// GinRoute convert gin to cango
func (can *Can) GinRoute(eg *gin.Engine) {
	for _, gh := range can.ToGins() {
		for hm, _ := range gh.HttpMethods {
			eg.Handle(hm, gh.Url, gh.Handle)
		}
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

// 取得项目运行目录
func getRootPath() string {
	// 先判断当前路径是否有效
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	// 判断当前路径是否为绝对路径
	if abs, err := filepath.Abs(os.Args[0]); err == nil {
		return filepath.Dir(abs)
	}
	return os.Args[0]
}

func (can *Can) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := &WebRequest{
		ResponseWriter: rw,
		Request:        r,
	}

	defer func() {
		if err := recover(); err != nil {
			var buf [1024 * 10]byte
			runtime.Stack(buf[:], false)
			canlog.CanError(err)
			canlog.CanError(string(buf[0:]))
			request.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		}
	}()

	// todo filter should not be here?????
	needStop := false
	var filterReturn interface{}
	var needHandle = true
	var filterChan []Filter
	for _, dsp := range can.filterDispatcher {
		match := doubleMatch(dsp.dispatcher, r)
		if match.Error() == nil {
			filterChan = append(filterChan, dsp.filter)
			ri := dsp.filter.PreHandle(request)
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
					needHandle = false
				}
			}
		}
	}
	if needStop {
		request.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	var handleReturn interface{}
	var statusCode int

	if needHandle {
		handleReturn, statusCode = serve(can.routeMux, request)
		if statusCode == serveFallbackCode {
			if can.fallbackHandler != nil {
				can.fallbackHandler.ServeHTTP(request.ResponseWriter, r)
				return
			}
			statusCode = http.StatusNotFound
		}
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
			request.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := tpl.Execute(request.ResponseWriter, mv.Model)
			if err != nil {
				canlog.CanError("template error", err)
			}
		}
	case Redirect:
		http.Redirect(request.ResponseWriter, r, handleReturn.(Redirect).Url, http.StatusFound)
	case RedirectWithCode:
		code := handleReturn.(RedirectWithCode).Code
		if code == 0 {
			code = http.StatusFound
		}
		http.Redirect(request.ResponseWriter, r, handleReturn.(RedirectWithCode).Url, code)
	case Content:
		request.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		request.ResponseWriter.WriteHeader(http.StatusOK)
		_, err := request.ResponseWriter.Write([]byte(handleReturn.(Content).String))
		if err != nil {
			canlog.CanError(err)
		}
	case ContentWithCode:
		code := handleReturn.(ContentWithCode).Code
		if code == 0 {
			code = http.StatusOK
		}
		request.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		request.ResponseWriter.WriteHeader(code)
		_, err := request.ResponseWriter.Write([]byte(handleReturn.(Content).String))
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
		// todo here is a bug when file ends with index.html
		// should use ServeContent instead of ServeFile
		paths := [...]string{can.rootPath + path, can.staticRootPath + path, path}
		for _, p := range paths {
			_, err = os.Stat(p)
			if err != nil {
				continue
			}
			path = p
			break
		}
		// todo rw???
		// 这里使用rw 而没有使用request.ResponseWriter 应该保持统一才好
		// 但是在ServerFile中会直接写入一些数据，产生不可知的影响
		if err != nil {
			canlog.CanDebug(err, "can't find the file", handleReturn)
			errorHandleMap[404](rw, r)
		} else {
			http.ServeFile(rw, r, path)
		}
	case DoNothing:
		if statusCode == http.StatusNotFound {
			errorHandleMap[404](request.ResponseWriter, r)
		}
	default:
		request.ResponseWriter.WriteHeader(statusCode)
		bs, err := responseJsonHandler(handleReturn)
		if err == nil {
			_, _ = request.ResponseWriter.Write(bs)
		} else {
			canlog.CanError(err)
			_, _ = request.ResponseWriter.Write([]byte("{}"))
		}
	}
	// postHandle
	for _, f := range filterChan {
		f.PostHandle(request)
	}
}

func doubleMatch(mux dispatcher, req *http.Request) matcher {
	match := mux.Match(req)
	if match.Error() != nil {
		originalPath := req.URL.Path
		req.URL.Path = filepath.Clean(originalPath)
		// some url like "/api//user.json" and "/api/user.json",should match again
		if originalPath == req.URL.Path {
			return match
		}
		match = mux.Match(req)
		req.URL.Path = originalPath
		return match
	}
	return match
}

const serveFallbackCode = -1

func serve(mux dispatcher, request *WebRequest) (interface{}, int) {
	req := request.Request
	match := doubleMatch(mux, req)
	if match.Error() != nil {
		canlog.CanError(req.Method, req.URL.Path, match.Error())
		return nil, serveFallbackCode
	}
	invoker := match.Forwarder().GetInvoker()
	uriRequestValue := reflect.ValueOf(newContext(request))
	callerIn := make([]reflect.Value, invoker.Type.NumIn())
	cookies := req.Cookies()
	_ = req.ParseForm()
	gs, _ := gorillaStore.Get(request.Request, cangoSessionKey)
	var bodyBytes []byte
	var isParse bool

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
		// 读取header，只赋值有header标签的变量
		// 读取cookie，只赋值有cookie标签的变量
		// 解析session，赋值有session标签的变量
		reqHolder := func(valueKey string, et entityType) *entityValue {
			switch et {
			case cookieEntity:
				for _, cookie := range cookies {
					if cookie.Name == valueKey {
						return &entityValue{
							enc:   stringFlag,
							key:   cookie.Name,
							value: cookie.Value,
						}
					}
				}
				return nil
			case headerEntity:
				if i, ok := req.Header[valueKey]; ok {
					return &entityValue{
						enc:   stringFlag,
						key:   valueKey,
						value: i,
					}
				}
				return nil
			case sessionEntity:
				if i, ok := gs.Values[valueKey]; ok {
					return &entityValue{
						enc:   gobBytes,
						key:   valueKey,
						value: i,
					}
				}
				return nil
			default:
				if v, ok := match.GetVars()[valueKey]; ok {
					return &entityValue{
						enc:   stringFlag,
						key:   valueKey,
						value: v,
					}
				}
				if v, ok := req.Form[valueKey]; ok {
					return &entityValue{
						enc:   strSliceFlag,
						key:   valueKey,
						value: v,
					}
				}
				return nil
			}
		}

		doDecode(addr(callerIn[i]), reqHolder, fieldTagNames)
		// TODO redesign 这个接口
		if in.Implements(constructorType) {
			uriFiled := value(callerIn[i]).FieldByName(constructorTypeName)
			if uriFiled.IsValid() && uriFiled.CanSet() {
				uriFiled.Set(valueOfEmptyConstructor)
			}
			addr(callerIn[i]).Interface().(Constructor).Construct(request)
		}
		// TODO 如果是json的类型
		ct := request.Request.Header.Get("Content-Type")
		// 不太精准
		if strings.Contains(strings.ToLower(ct), "json") {
			if !isParse {
				isParse = true
				bs, err := io.ReadAll(request.Request.Body)
				if err != nil {
					canlog.CanError(err)
				}
				bodyBytes = bs
			}
			to := addr(callerIn[i]).Interface()
			err := json.Unmarshal(bodyBytes, to)
			if err != nil {
				canlog.CanError(err)
			}
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
	if vs[0].IsValid() == false {
		return nil, http.StatusMethodNotAllowed
	}
	if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
		if vs[0].Elem().IsValid() {
			return vs[0].Elem().Interface(), http.StatusOK
		}
	}
	return vs[0].Interface(), http.StatusOK
}

var cangoMark = []byte{10, 10, 10, 32, 32, 95, 95, 95, 95, 95, 95, 32, 32, 32, 32, 32, 95, 95, 95, 32, 32, 32, 32, 32, 32, 46, 95, 95, 32, 32, 32, 95, 95, 46, 32, 32, 32, 95, 95, 95, 95, 95, 95, 95, 32, 32, 32, 95, 95, 95, 95, 95, 95, 32, 32, 32, 32, 95, 95, 32, 32, 10, 32, 47, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 47, 32, 32, 32, 92, 32, 32, 32, 32, 32, 124, 32, 32, 92, 32, 124, 32, 32, 124, 32, 32, 47, 32, 32, 95, 95, 95, 95, 95, 124, 32, 47, 32, 32, 95, 95, 32, 32, 92, 32, 32, 124, 32, 32, 124, 32, 10, 124, 32, 32, 44, 45, 45, 45, 45, 39, 32, 32, 47, 32, 32, 94, 32, 32, 92, 32, 32, 32, 32, 124, 32, 32, 32, 92, 124, 32, 32, 124, 32, 124, 32, 32, 124, 32, 32, 95, 95, 32, 32, 124, 32, 32, 124, 32, 32, 124, 32, 32, 124, 32, 124, 32, 32, 124, 32, 10, 124, 32, 32, 124, 32, 32, 32, 32, 32, 32, 47, 32, 32, 47, 95, 92, 32, 32, 92, 32, 32, 32, 124, 32, 32, 46, 32, 96, 32, 32, 124, 32, 124, 32, 32, 124, 32, 124, 95, 32, 124, 32, 124, 32, 32, 124, 32, 32, 124, 32, 32, 124, 32, 124, 32, 32, 124, 32, 10, 124, 32, 32, 96, 45, 45, 45, 45, 46, 47, 32, 32, 95, 95, 95, 95, 95, 32, 32, 92, 32, 32, 124, 32, 32, 124, 92, 32, 32, 32, 124, 32, 124, 32, 32, 124, 95, 95, 124, 32, 124, 32, 124, 32, 32, 96, 45, 45, 39, 32, 32, 124, 32, 124, 95, 95, 124, 32, 10, 32, 92, 95, 95, 95, 95, 95, 95, 47, 95, 95, 47, 32, 32, 32, 32, 32, 92, 95, 95, 92, 32, 124, 95, 95, 124, 32, 92, 95, 95, 124, 32, 32, 92, 95, 95, 95, 95, 95, 95, 124, 32, 32, 92, 95, 95, 95, 95, 95, 95, 47, 32, 32, 40, 95, 95, 41, 32, 10, 10, 10}
