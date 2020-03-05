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
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/JessonChan/canlog"
)

const emptyPrefix = ""

// todo route by controller and method Name???
func (can *Can) Route(uris ...URI) *Can {
	return can.RouteWithPrefix(emptyPrefix, uris...)
}

// todo route with suffix and simplify
func (can *Can) RouteWithPrefix(prefix string, uris ...URI) *Can {
	for _, uri := range uris {
		can.route(prefix, uri)
	}
	return can
}

func (can *Can) route(prefix string, uri URI) {
	typ := toPtrKind(uri)
	can.ctrlEntryMap[prefix+typ.String()] = ctrlEntry{prefix: prefix, kind: reflect.Ptr, ctrl: uri, tim: time.Now().Unix()}
}

func (can *Can) RouteFunc(fns ...interface{}) *Can {
	return can.RouteFuncWithPrefix(emptyPrefix, fns...)
}
func (can *Can) RouteFuncWithPrefix(prefix string, fns ...interface{}) *Can {
	for _, fn := range fns {
		can.routeFunc(prefix, fn)
	}
	return can
}
func (can *Can) routeFunc(prefix string, fn interface{}) {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		canlog.CanInfo("can't forwarder func with ", fv.Kind())
		return
	}
	funcMethod := reflect.Method{
		Name:    runtime.FuncForPC(fv.Pointer()).Name(),
		PkgPath: "",
		Type:    fv.Type(),
		Func:    fv,
		Index:   0,
	}
	can.ctrlEntryMap[prefix+fv.String()] = ctrlEntry{prefix: prefix, kind: reflect.Func, fn: funcMethod, tim: time.Now().Unix()}
}

var uriRegMap = map[URI]bool{}

// todo with prefix???
// todo with can app Name ???
func RegisterURI(uri URI) bool {
	uriRegMap[uri] = true
	return true
}

type ctrlEntry struct {
	prefix string
	kind   reflect.Kind
	ctrl   interface{}
	fn     reflect.Method
	tim    int64
}

type sortCtrlEntry []ctrlEntry

func (s sortCtrlEntry) Len() int           { return len(s) }
func (s sortCtrlEntry) Less(i, j int) bool { return s[i].tim < s[j].tim }
func (s sortCtrlEntry) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (can *Can) buildRoute() {
	for uri, _ := range uriRegMap {
		can.route("", uri)
	}
	var ces []ctrlEntry
	for _, ce := range can.ctrlEntryMap {
		ces = append(ces, ce)
	}
	sort.Sort(sortCtrlEntry(ces))
	for _, ce := range ces {
		can.buildSingleRoute(ce)
	}
}

// todo 如果static存在的文件非常多的时候，这种实现方式会成为巨大的问题
// todo 这里也有个优势，就是防止被恶意请求，请求某个不存在的文件，会直接在调用io之前被拒绝
func (can *Can) buildStaticRoute() {
	_ = filepath.Walk(can.staticRootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		can.route(filepath.Clean("/"+strings.TrimPrefix(path, can.rootPath)), &staticController{})
		return nil
	})
	// todo 特殊处理favicon.ico和robots.txt
	can.route("/favicon.ico", &staticController{})
	can.route("/robots.txt", &staticController{})
}

func (can *Can) buildSingleRoute(ce ctrlEntry) {
	switch ce.kind {
	case reflect.Ptr:
		hs := factory(ce.ctrl)
		ctrlTagPaths, ctlName := urlStr(hs.typ.Elem())
		for _, hm := range hs.fns {
			can.routeMethod(invokeByReceiver, ce.prefix, hm.fn, ctlName+"."+hm.fn.Name, ctrlTagPaths)
		}
	case reflect.Func:
		can.routeMethod(invokeBySelf, ce.prefix, ce.fn, "RouteFunc."+ce.fn.Name, nil)
	}
}

// todo use factory to clean code
func (can *Can) routeMethod(invokeByWho int, prefix string, m reflect.Method, routerName string, ctrlTagPaths []string) {
	hm := factoryMethod(m)
	if hm == nil {
		return
	}
	for _, hp := range hm.patterns {
		route := can.routeMux.NewForwarder(routerName)
		can.methodMap[routerName] = &invoker{kind: invokeByWho, Method: &m}
		for _, path := range combinePaths(prefix, ctrlTagPaths, hp.path) {
			// default method is get
			httpMethods := defaultHttpMethods
			if len(hp.httpMethods) > 0 {
				httpMethods = hp.httpMethods
			}
			route.PathMethods(path, httpMethods...)
			canlog.CanDebug(routerName, path, httpMethods)
		}
	}
}

func combinePaths(prefix string, ctrlTagPaths []string, methodTagPath string) (paths []string) {
	if len(ctrlTagPaths) == 0 {
		ctrlTagPaths = []string{""}
	}
	for _, ctrlTagPaths := range ctrlTagPaths {
		paths = append(paths, filepath.Clean(strings.Join([]string{prefix, ctrlTagPaths, methodTagPath}, "/")))
	}
	return paths
}
