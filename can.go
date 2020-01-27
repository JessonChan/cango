package cango

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Can struct {
	srv          *http.Server
	viewRootPath string
	rootRouter   *mux.Router
	methodMap    map[string]reflect.Method
	filterMap    map[string][]Filter
	ctrlMap      map[string]ctrlEntry
}

type ctrlEntry struct {
	prefix string
	vl     reflect.Value
	ctrl   URI
}

var defaultAddr = Addr{Host: "", Port: 8080}

func NewCan() *Can {
	return &Can{
		srv:        &http.Server{Addr: defaultAddr.String()},
		filterMap:  make(map[string][]Filter, 1),
		methodMap:  map[string]reflect.Method{},
		rootRouter: mux.NewRouter(),
		ctrlMap:    map[string]ctrlEntry{},
	}
}

type Addr struct {
	Host string
	Port int
}

type View struct {
	RootPath string
}

func (addr Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

type responseTypeHandler func(interface{}) ([]byte, error)

var responseHandler responseTypeHandler = json.Marshal

func (can *Can) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
			_ = tpl.Execute(rw, mv.Model)
		}
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
	can.viewRootPath = getViewRootPath(as)
	can.buildRoute()

	startChan := make(chan interface{}, 1)
	var err error
	go func() { err = can.srv.ListenAndServe() }()
	if err != nil {
		fmt.Println("cango start failed @ " + addr.String())
		panic(err)
	}
	fmt.Println("cango start success @ " + addr.String())
	<-startChan
}

func getAddr(as []interface{}) Addr {
	for _, v := range as {
		if addr, ok := v.(Addr); ok {
			return addr
		}
	}
	return defaultAddr
}

func getViewRootPath(as []interface{}) string {
	for _, v := range as {
		if view, ok := v.(View); ok {
			return view.RootPath
		}
	}
	abs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return os.Args[0]
	}
	return filepath.Dir(abs)
}

var uriType = reflect.TypeOf((*URI)(nil)).Elem()

/*
	methods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	for _, m := range methods {
		fmt.Printf("	reflect.TypeOf((*%sMethod)(nil)).Elem(): http.Method%s,\n", m, m)
	}
*/

var httpMethodMap = map[reflect.Type]string{
	reflect.TypeOf((*GetMethod)(nil)).Elem():     http.MethodGet,
	reflect.TypeOf((*PostMethod)(nil)).Elem():    http.MethodPost,
	reflect.TypeOf((*HeadMethod)(nil)).Elem():    http.MethodHead,
	reflect.TypeOf((*PutMethod)(nil)).Elem():     http.MethodPut,
	reflect.TypeOf((*PatchMethod)(nil)).Elem():   http.MethodPatch,
	reflect.TypeOf((*DeleteMethod)(nil)).Elem():  http.MethodDelete,
	reflect.TypeOf((*OptionsMethod)(nil)).Elem(): http.MethodOptions,
	reflect.TypeOf((*TraceMethod)(nil)).Elem():   http.MethodTrace,
}

// urlStr get uri from tag value
func (can *Can) urlStr(uri interface{}) (string, string) {
	typ := reflect.TypeOf(uri)
	if typ.Kind() != reflect.Struct {
		panic(uri)
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if f.Type == uriType {
			return f.Tag.Get("value"), typ.Name()
		}
	}
	return "", ""
}

type StatusCode int

var decoder = schema.NewDecoder()

func (can *Can) serve(req *http.Request) (interface{}, StatusCode) {
	match := &mux.RouteMatch{}
	can.rootRouter.Match(req, match)
	if match.MatchErr != nil {
		return nil, http.StatusNotFound
	}
	fs, _ := can.filterMap[match.Route.GetName()]
	if len(fs) > 0 {
		for _, f := range fs {
			if f == nil {
				continue
			}
			f.PreHandle(req)
		}
		// do filter
	}

	m, ok := can.methodMap[match.Route.GetName()]
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
	if func(methods []string, err error) bool {
		if err != nil {
			return false
		}
		for _, m := range methods {
			switch m {
			case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch:
				return true
			default:
				return false
			}
		}
		return false
	}(match.Route.GetMethods()) {
		_ = req.ParseForm()
		_ = decoder.Decode(mt.Addr().Interface(), req.Form)
	}
	if len(match.Vars) > 0 {
		vars := toValues(match.Vars)
		_ = decoder.Decode(ct.Interface(), vars)
		_ = decoder.Decode(mt.Addr().Interface(), vars)
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

func lowerCase(str string) string {
	if str == "" {
		return str
	}
	bs := []byte(str)
	if 'A' <= bs[0] && bs[0] <= 'Z' {
		bs[0] += 'a' - 'A'
	}
	return string(bs)
}

func upperCase(str string) string {
	if str == "" {
		return str
	}
	bs := []byte(str)
	if 'a' <= bs[0] && bs[0] <= 'z' {
		bs[0] -= 'a' - 'A'
	}
	return string(bs)
}

func (can *Can) RouteWithPrefix(prefix string, uris ...URI) *Can {
	for _, uri := range uris {
		can.route(prefix, uri)
	}
	return can
}

const emptyPrefix = ""

func (can *Can) Route(uris ...URI) *Can {
	return can.RouteWithPrefix(emptyPrefix, uris...)
}

func (can *Can) route(prefix string, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("route controller must be ptr")
	}
	can.ctrlMap[prefix+rp.String()] = ctrlEntry{prefix: prefix, vl: rp, ctrl: uri}
}

func (can *Can) buildRoute() {
	for _, ce := range can.ctrlMap {
		can.buildSingleRoute(ce)
	}
}

func (can *Can) buildSingleRoute(ce ctrlEntry) {
	prefix := ce.prefix
	rp := ce.vl
	uri := ce.ctrl

	urlStr, ctlName := can.urlStr(reflect.Indirect(rp).Interface())
	urlStr = prefix + urlStr
	tvp := reflect.TypeOf(uri)
	for i := 0; i < tvp.NumMethod(); i++ {
		m := tvp.Method(i)
		if m.PkgPath != "" {
			continue
		}
		for i := 0; i < m.Type.NumIn(); i++ {
			in := m.Type.In(i)
			if in.Kind() != reflect.Struct {
				continue
			}
			routerName := ctlName + "." + m.Name
			route := can.rootRouter.Name(routerName)
			var httpMethods []string
			for j := 0; j < in.NumField(); j++ {
				f := in.Field(j)
				if f.PkgPath != "" {
					panic("could not use unexpected filed in param:" + f.Name)
				}
				switch f.Type {
				case uriType:
					route.Path(urlStr + f.Tag.Get("value"))
					can.methodMap[routerName] = m
					can.filterMap[routerName] = can.filterMap[rp.String()]
				}
				m, ok := httpMethodMap[f.Type]
				if ok {
					httpMethods = append(httpMethods, m)
				}
			}
			// default method is get
			if len(httpMethods) == 0 {
				httpMethods = append(httpMethods, http.MethodGet)
			}
			route.Methods(httpMethods...)
		}
	}
}
