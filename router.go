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
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/JessonChan/canlog"
)

type (
	CanMux interface {
		NewRouter(name string) CanRouter
		Match(req *http.Request) CanMatcher
	}

	CanRouter interface {
		Path(ps ...string)
		Methods(ms ...string)
		GetName() string
		GetMethods() []string
		GetPath() string
	}

	CanMatcher interface {
		Error() error
		Route() CanRouter
		GetVars() map[string]string
	}
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
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("route controller must be ptr")
	}
	can.ctrlMap[prefix+rp.String()] = ctrlEntry{prefix: prefix, vl: rp, ctrl: uri, tim: time.Now().Unix()}
}

type ctrlEntry struct {
	prefix string
	vl     reflect.Value
	ctrl   URI
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
	for _, ce := range can.ctrlMap {
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
	_ = filepath.Walk(can.rootPath+staticDir, func(path string, info os.FileInfo, err error) error {
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
	prefix := ce.prefix
	rp := ce.vl
	uri := ce.ctrl

	strUrls, ctlName := can.urlStr(reflect.Indirect(rp).Type())
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
			route := can.rootMux.NewRouter(routerName)
			var httpMethods []string
			for j := 0; j < in.NumField(); j++ {
				f := in.Field(j)
				if f.PkgPath != "" {
					canlog.CanError("could not use unexpected filed in param:" + f.Name)
				}
				switch f.Type {
				case uriType:
					var paths []string
					for _, path := range tagUriParse(f.Tag) {
						for _, strUrl := range strUrls {
							paths = append(paths, filepath.Clean(strings.Join([]string{prefix, strUrl, path}, "/")))
						}
					}
					route.Path(paths...)
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
			canlog.CanDebug(route.GetName(), route.GetPath(), route.GetMethods())
		}
	}
}

// urlStr get uri from tag value
func (can *Can) urlStr(typ reflect.Type) ([]string, string) {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if f.Type == uriType {
			return tagUriParse(f.Tag), typ.Name()
		}
	}
	return []string{}, ""
}

func tagUriParse(tag reflect.StructTag) []string {
	return strings.Split(tag.Get(uriTagName), ";")
}
