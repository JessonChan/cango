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

type Cookie interface {
	constructor
	CookiePlaceHolder()
}

var cookieType = reflect.TypeOf((*Cookie)(nil)).Elem()
var cookieTypeName = cookieType.Name()
var valueOfEmptyCookie = reflect.ValueOf(&emptyCookieConstructor{})

type emptyCookieConstructor struct {
}

func (e *emptyCookieConstructor) Construct(r *http.Request) {
}
func (e *emptyCookieConstructor) CookiePlaceHolder() {
}

func cookieConstruct(r *http.Request, cs Cookie) Cookie {
	if reflect.TypeOf(cs) == cookieType {
		return cs
	}
	csv := reflect.ValueOf(cs)
	if csv.Kind() == reflect.Ptr {
		csv = csv.Elem()
	}
	// 先从默认值进行赋值
	cookies := r.Cookies()
	checkSet(stringFlag, cookieHolder(cookies), csv, cookieNameWithTag)
	// 最后执行自定义函数的New方法
	// 如果没有自定义方法，则执行结束
	cs.Construct(r)
	return cs
}

func cookieNameWithTag(field reflect.StructField) []string {
	return filedName(field, cookieTagName)
}
