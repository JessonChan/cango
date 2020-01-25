package cango

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

type Can struct {
	srv          *http.Server
	viewRootPath string
}

var defaultAddr = Addr{Host: "", Port: 8080}

func NewCan() *Can {
	return &Can{srv: &http.Server{Addr: defaultAddr.String()}}
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
	rt, statusCode := can.Do(r)
	mv, ok := rt.(ModelView)
	if ok {
		tpl := can.lookupTpl(mv.Tpl)
		if tpl != nil {
			_ = tpl.Execute(rw, mv.Model)
		}
		return
	}
	rw.WriteHeader(int(statusCode))
	if rt == nil {
		_, _ = rw.Write([]byte("{}"))
		return
	}
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

var rootRouter = mux.NewRouter()
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

func assignValue(value reflect.Value, vars map[string]string) {
	for k, v := range vars {
		f := value.FieldByName(upperCase(k))
		if f.IsValid() == false {
			continue
		}
		switch f.Kind() {
		case reflect.String:
			f.Set(reflect.ValueOf(v))
		case reflect.Int:
			i, _ := strconv.ParseInt(v, 10, 0)
			f.SetInt(i)
		}
	}
}

type StatusCode int

func (can *Can) Do(req *http.Request) (interface{}, StatusCode) {
	match := &mux.RouteMatch{}
	rootRouter.Match(req, match)
	if match.MatchErr != nil {
		return nil, http.StatusNotFound
	}
	m, ok := methodMap[match.Route.GetName()]
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

	assignValue(ct.Elem(), match.Vars)
	assignValue(mt, match.Vars)
	vs := ct.MethodByName(m.Name).Call([]reflect.Value{mt})
	if len(vs) == 0 {
		return nil, http.StatusMethodNotAllowed
	}
	if vs[0].Kind() == reflect.Ptr || vs[0].Kind() == reflect.Interface {
		return vs[0].Elem().Interface(), http.StatusOK
	}
	return vs[0].Interface(), http.StatusOK
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

var methodMap = map[string]reflect.Method{}

func (can *Can) RouteWithPrefix(prefix string, uris ...URI) {
	for _, uri := range uris {
		can.route(prefix, uri)
	}
}

const emptyPrefix = ""

func (can *Can) Route(uris ...URI) {
	can.RouteWithPrefix(emptyPrefix, uris...)
}

func (can *Can) route(prefix string, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("route controller must be prt")
	}
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
			route := rootRouter.Name(routerName)
			var httpMethods []string
			for j := 0; j < in.NumField(); j++ {
				f := in.Field(j)
				if f.PkgPath != "" {
					panic("could not use unexpected filed in param:" + f.Name)
				}
				switch f.Type {
				case uriType:
					route.Path(urlStr + f.Tag.Get("value"))
					methodMap[routerName] = m
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
