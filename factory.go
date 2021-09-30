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
	"reflect"
	"strings"
)

type (
	handlerStruct struct {
		typ reflect.Type
		fns []*handlerMethod
	}
	handlerMethod struct {
		fn       reflect.Method
		patterns []*handlePath
	}

	handlePath struct {
		path        string
		httpMethods []string
	}
)

var cacheStruct = map[reflect.Type]*handlerStruct{}
var cacheMethod = map[reflect.Method]*handlerMethod{}

func factory(i interface{}) *handlerStruct {
	return factoryType(toPtrKind(i))
}
func toPtrKind(i interface{}) reflect.Type {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		// 转换成指针类型
		typ = reflect.New(typ).Type()
	}
	return typ
}

func implements(in, typ reflect.Type) reflect.Type {
	if in.Implements(typ) {
		if in.Kind() == reflect.Ptr {
			in = in.Elem()
		}
		return in
	}
	return nil
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
		hm := factoryMethod(m, invokeBySelf)
		if hm != nil {
			hs.fns = append(hs.fns, hm)
		}

	}
	return hs
}

func factoryMethod(m reflect.Method, invokeByWho int) *handlerMethod {
	if hm, ok := cacheMethod[m]; ok {
		return hm
	}
	if uriInterfaceContains(m) {
		return nil
	}
	// 只接受参数列表最后一个做为uri,如果在参数中不包含uri，则放弃该路由
	// todo 支持根据路由名和方法名进行路由定义？？？
	for j := m.Type.NumIn(); j > invokeByWho; j-- {
		in := implements(m.Type.In(j-1), uriType)
		if in == nil {
			continue
		}
		switch in.Kind() {
		// 此时，in的类型为cango.URI，形如这种
		// func (c *Controller)Ping(cango.URI)interface{}
		case reflect.Interface:
			hm := &handlerMethod{
				fn: m,
				patterns: func() (pms []*handlePath) {
					return []*handlePath{{
						path:        "",
						httpMethods: defaultHTTPMethods,
					}}
				}(),
			}
			cacheMethod[m] = hm
			return hm
		case reflect.Struct:
			uriFiled, ok := in.FieldByName(uriName)
			if !ok {
				continue
			}
			paths := tagUriParse(uriFiled.Tag)
			methods := func() (foundHttpMethods []string) {
				for k, v := range httpMethodMap {
					if _, ok := in.FieldByName(k.Name()); ok {
						foundHttpMethods = append(foundHttpMethods, v)
					}
				}
				return
			}()
			hm := &handlerMethod{fn: m}
			// 没有在方法定义路径，需要用空的字段把方法带出去
			if len(paths) == 0 {
				paths = []string{""}
			}
			// 依然有改进空间
			for _, path := range paths {
				hm.patterns = append(hm.patterns, &handlePath{
					path:        path,
					httpMethods: methods,
				})
			}
			cacheMethod[m] = hm
			return hm
		}
	}
	return nil
}

// urlStr get uri from tag value
func urlStr(typ reflect.Type) ([]string, string) {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		// 不为空，则为不可导出变量
		if f.PkgPath != "" {
			continue
		}
		if f.Type == uriType {
			return tagUriParse(f.Tag), typ.PkgPath() + "." + typ.Name()
		}
	}
	return []string{}, ""
}

func tagUriParse(tag reflect.StructTag) []string {
	if s := tag.Get(uriTagName); s != "" {
		return strings.Split(s, ";")
	}
	return []string{}
}
