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
	"path/filepath"
	"reflect"
)

// todo 为什么filter 不使用和URI一样的方式进行注册
type Filter interface {
	PreHandle(req *http.Request) interface{}
	// todo
	// PostHandle()
	// AfterHandled()
}

type MyFilter struct {
	Filter `value:"/api/*"`
}

var filterRegMap = map[Filter]bool{}
var filterStructMap = map[reflect.Type]Filter{}

func RegisterFilter(filter Filter) bool {
	filterRegMap[filter] = true
	return true
}

func (can *Can) buildFilter() {
	for filter, _ := range filterRegMap {
		can.Filter(filter)
	}
	for typ, flt := range filterStructMap {
		hs := factoryType(typ)
		urls, _ := urlStr(typ.Elem())
		for _, fn := range hs.fns {
			for _, pattern := range fn.patterns {
				if len(urls) > 0 {
					for _, url := range urls {
						can.buildSingleFilter(flt, filepath.Clean(url+"/"+pattern.path), pattern.methods)
					}
				} else {
					can.buildSingleFilter(flt, filepath.Clean(pattern.path), pattern.methods)
				}
			}
		}
	}
}

func (can *Can) buildSingleFilter(f Filter, path string, methods []string) {
	if path == "" {
		return
	}
	if len(methods) == 0 {
		methods = defaultHttpMethods
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Ptr {
		panic("filter must be ptr")
	}

	name := fv.Type().Name()
	can.filterMux.NewRouter(name).PathMethods(path, methods...)
	can.filterMap[name] = f
}

func (can *Can) filter(f Filter, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("filter controller must be ptr")
	}
	filterStructMap[rp.Type()] = f
}

func (can *Can) Filter(f Filter, uris ...URI) *Can {
	rp := reflect.ValueOf(f)
	if rp.Kind() == reflect.Ptr {
		rp = rp.Elem()
	}
	for i := 0; i < rp.NumField(); i++ {
		tp := rp.Field(i).Type()
		if tp.Kind() != reflect.Ptr {
			// todo warning
			continue
		}
		if tp.Implements(uriType) {
			uris = append(uris, rp.Field(i).Interface().(URI))
		}
	}

	for _, uri := range uris {
		can.filter(f, uri)
	}
	return can
}
