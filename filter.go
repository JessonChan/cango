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
	// for typ, flt := range filterStructMap {
	//
	// }
}

// func (can *Can) buildSingleFilter(f Filter, paths []string, methods []string) {
// 	if len(paths) == 0 {
// 		return
// 	}
// 	if len(methods) == 0 {
// 		methods = []string{http.MethodGet}
// 	}
// 	fv := reflect.ValueOf(f)
// 	if fv.Kind() != reflect.Ptr {
// 		panic("filter must be ptr")
// 	}
//
// 	fwd := can.filterMux.NewRouter(fv.Type().Name())
// 	fwd.Path(paths...)
// 	fwd.Methods(methods...)
// 	can.filterMap[fwd.GetName()] = f
// }

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
