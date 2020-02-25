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
	"strings"
)

/**
struct 一对多 method
method 一对多  path
path 一对多 httpMethod

查找时(path,httpMethod)
*/

type pathMethods struct {
	path    string
	methods []string
}

type handlerMethod struct {
	fn       reflect.Method
	patterns []*pathMethods
}

type handlerStruct struct {
	typ reflect.Type
	fns []*handlerMethod
}

var cacheStruct = map[reflect.Type]*handlerStruct{}

func factory(i interface{}) *handlerStruct {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		panic("must be a ptr")
	}
	return factoryType(typ)
}
func factoryType(typ reflect.Type) *handlerStruct {
	hs, ok := cacheStruct[typ]
	if ok {
		return hs
	}
	hs = &handlerStruct{typ: typ}
	cacheStruct[typ] = hs
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.PkgPath != "" {
			continue
		}
		// todo 现在只接受一个参数???
		for j := 0; j < m.Type.NumIn(); j++ {
			in := m.Type.In(j)
			if in.Implements(uriType) == false {
				continue
			}
			switch in.Kind() {
			case reflect.Interface:
				hs.fns = append(hs.fns, &handlerMethod{
					fn: m,
					patterns: func() (pms []*pathMethods) {
						return []*pathMethods{{
							path:    "",
							methods: []string{http.MethodGet},
						}}
					}(),
				})
			case reflect.Struct:
				uriFiled, ok := in.FieldByName(uriName)
				if !ok {
					continue
				}
				paths := tagUriParse(uriFiled.Tag)
				methods := func() (methods []string) {
					for k, v := range httpMethodMap {
						if _, ok := in.FieldByName(k.Name()); ok {
							methods = append(methods, v)
						}
					}
					return
				}()
				hm := &handlerMethod{fn: m}
				hs.fns = append(hs.fns, hm)
				for _, path := range paths {
					hm.patterns = append(hm.patterns, &pathMethods{
						path:    path,
						methods: methods,
					})
				}
			}
		}
	}
	return hs
}

// urlStr get uri from tag value
func urlStr(typ reflect.Type) ([]string, string) {
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
