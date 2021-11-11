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
	"path/filepath"
	"reflect"
)

// todo 为什么filter 不使用和URI一样的方式进行注册
type Filter interface {
	// PreHandle is used to perform operations before sending the request to the controller.
	// This method should return true to continue the request serve.
	// If this method returns false,the request will stop.
	// If this method returns cango-return type(Redirect/ModelView...),the request will response with the type
	PreHandle(request *WebRequest) interface{}
	// PostHandle is used to perform operations before sending the response to the client.
	// This method should return true.
	PostHandle(request *WebRequest) interface{}
	// todo
	// AfterHandled()
}

var _ = Filter(&emptyFilter{})

type FilterType reflect.Type

type filterDispatcher struct {
	filter     Filter
	dispatcher dispatcher
	uriTypes   []reflect.Type
}

func newFilterDispatcher(filter Filter) *filterDispatcher {
	fv := reflect.ValueOf(filter).Elem()
	ffv := fv.FieldByName(filterName)
	if ffv.CanSet() && ffv.IsNil() {
		ffv.Set(filterImpl)
	}
	return &filterDispatcher{filter: filter, dispatcher: newCanMux()}
}

type emptyFilter struct {
}

func (*emptyFilter) PreHandle(request *WebRequest) interface{} {
	return true
}
func (*emptyFilter) PostHandle(request *WebRequest) interface{} {
	return true
}

var filterImpl = reflect.ValueOf(&emptyFilter{})

var filterType = reflect.TypeOf((*Filter)(nil)).Elem()
var filterName = filterType.Name()

var filterRegMap = map[Filter]bool{}

func RegisterFilter(filter Filter) bool {
	filterRegMap[filter] = true
	return true
}

func (can *Can) buildFilter() {
	for filter, _ := range filterRegMap {
		can.Filter(filter)
	}

	for flt, fd := range can.filterDispatcher {
		paths, methods := getPaths(flt)
		for _, path := range paths {
			buildSingleFilter(fd.dispatcher, fd.filter, filepath.Clean(path), methods)
		}
		for _, typ := range fd.uriTypes {
			hs := factoryType(typ)
			urls, _ := urlStr(typ.Elem())
			if len(urls) == 0 {
				urls = []string{""}
			}
			for _, fn := range hs.fns {
				for _, pattern := range fn.patterns {
					for _, url := range urls {
						buildSingleFilter(fd.dispatcher, fd.filter, filepath.Clean(url+"/"+pattern.path), pattern.httpMethods)
					}
				}
			}
		}
	}
}

func getPaths(typ reflect.Type) ([]string, []string) {
	if typ.Implements(filterType) {
		if ff, ok := typ.Elem().FieldByName(filterName); ok {
			var httpMethods []string
			for i := 0; i < typ.Elem().NumField(); i++ {
				if m, ok := httpMethodMap[typ.Elem().Field(i).Type]; ok {
					httpMethods = append(httpMethods, m)
				}
			}
			if len(httpMethods) == 0 {
				httpMethods = allHTTPMethods
			}
			return tagUriParse(ff.Tag), httpMethods
		}
	}
	return []string{}, []string{}
}

func buildSingleFilter(dsp dispatcher, f Filter, path string, methods []string) {
	if path == "" {
		return
	}
	if len(methods) == 0 {
		methods = defaultHTTPMethods
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Ptr {
		panic("filter must be ptr")
	}
	name := fv.Elem().Type().Name()
	dsp.NewForwarder(name, &Invoker{kind: invokeByFilter, filter: f}).PathMethods(path, methods...)
}

func (can *Can) filter(f Filter, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("filter controller must be ptr")
	}
	typeOf := reflect.TypeOf(f)
	fd := can.filterDispatcher[typeOf]
	if fd == nil {
		fd = newFilterDispatcher(f)
	}
	contain := false
	for _, t := range fd.uriTypes {
		if t == rp.Type() {
			contain = true
			break
		}
	}
	if contain {
		return
	}
	fd.uriTypes = append(fd.uriTypes, rp.Type())
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
	// 只是注册Filter,路由使用Filter的Value字段
	if len(uris) == 0 {
		typeOf := reflect.TypeOf(f)
		fd := can.filterDispatcher[typeOf]
		if fd == nil {
			fd = newFilterDispatcher(f)
		}
		can.filterDispatcher[typeOf] = fd
	} else {
		for _, uri := range uris {
			can.filter(f, uri)
		}
	}
	return can
}
