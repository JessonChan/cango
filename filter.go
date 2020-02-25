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

var filterType = reflect.TypeOf((*Filter)(nil)).Elem()
var filterName = filterType.Name()

var filterRegMap = map[Filter]bool{}
var uriFilterMap = map[reflect.Type][]reflect.Type{}

func RegisterFilter(filter Filter) bool {
	filterRegMap[filter] = true
	return true
}

func (can *Can) buildFilter() {
	for filter, _ := range filterRegMap {
		can.Filter(filter)
	}
	for flt, typArr := range uriFilterMap {
		dsp, ok := can.filterMuxMap[flt]
		if !ok {
			dsp = newCanMux()
		}
		can.filterMuxMap[flt] = dsp
		paths, methods := getPaths(flt)
		for _, path := range paths {
			buildSingleFilter(dsp, can.filterMap[flt], filepath.Clean(path), methods)
		}
		for _, typ := range typArr {
			hs := factoryType(typ)
			urls, _ := urlStr(typ.Elem())
			for _, fn := range hs.fns {
				for _, pattern := range fn.patterns {
					if len(urls) > 0 {
						for _, url := range urls {
							buildSingleFilter(dsp, can.filterMap[flt], filepath.Clean(url+"/"+pattern.path), pattern.methods)
						}
					} else {
						buildSingleFilter(dsp, can.filterMap[flt], filepath.Clean(pattern.path), pattern.methods)
					}
				}
			}
		}
	}
}

func getPaths(typ reflect.Type) ([]string, []string) {
	if typ.Implements(filterType) {
		ff, ok := typ.Elem().FieldByName(filterName)
		if ok {
			paths := tagUriParse(ff.Tag)
			if len(paths) == 0 {
				return []string{}, []string{}
			}
			var httpMethods []string
			for i := 0; i < typ.Elem().NumField(); i++ {
				if m, ok := httpMethodMap[typ.Elem().Field(i).Type]; ok {
					httpMethods = append(httpMethods, m)
				}
			}
			if len(httpMethods) == 0 {
				httpMethods = allHttpMethods
			}
			return paths, httpMethods
		}
	}
	return []string{}, []string{}
}

func buildSingleFilter(dsp dispatcher, f Filter, path string, methods []string) {
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
	name := fv.Elem().Type().Name()
	dsp.NewForwarder(name).PathMethods(path, methods...)
}

func (can *Can) filter(f Filter, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("filter controller must be ptr")
	}
	ts := uriFilterMap[reflect.TypeOf(f)]
	contain := false
	for _, t := range ts {
		if t == rp.Type() {
			contain = true
			break
		}
	}
	if contain {
		return
	}
	uriFilterMap[reflect.TypeOf(f)] = append(ts, rp.Type())
	can.filterMap[reflect.TypeOf(f)] = f
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
	if len(uris) == 0 {
		if len(uriFilterMap[reflect.TypeOf(f)]) == 0 {
			uriFilterMap[reflect.TypeOf(f)] = nil
			can.filterMap[reflect.TypeOf(f)] = f
		}
	}

	for _, uri := range uris {
		can.filter(f, uri)
	}
	return can
}
