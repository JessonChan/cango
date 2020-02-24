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
	"net/http"
	"reflect"
	"time"
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

func RegisterFilter(filter Filter) bool {
	filterRegMap[filter] = true
	return true
}

type filterEntry struct {
	method reflect.Method
	f      Filter
}

func (can *Can) buildFilter() {
	for filter, _ := range filterRegMap {
		can.Filter(filter)
	}

	for _, fe := range can.filterEntryMap {
		pm := routeInfoMap[fe.method.Name]
		if pm != nil {
			can.buildSingleFilter(fe.f, pm.paths, pm.methods)
		}
	}
}

// 会存在一种情况  就是一个path 对应多个filter 这个应该如何体现呢？
// 也没有关系，我们的

func (can *Can) buildSingleFilter(f Filter, paths []string, methods []string) {
	if len(paths) == 0 {
		return
	}
	if len(methods) == 0 {
		methods = []string{http.MethodGet}
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Ptr {
		panic("filter must be ptr")
	}

	// 随机生成名字，可以保证不会重复，因为某个filter可以被多次注册
	fwdName := fmt.Sprint(time.Now().UnixNano())
	fwd := can.filterMux.NewRouter(fwdName)
	fwd.Path(paths...)
	fwd.Methods(methods...)
}

func (can *Can) filter(f Filter, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("filter controller must be ptr")
	}
	for i := 0; i < rp.Type().NumMethod(); i++ {
		// todo bug here method 会重复
		can.filterEntryMap[rp.Type().Method(i).Name] = filterEntry{
			method: rp.Type().Method(i),
			f:      f,
		}
	}
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
