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
	Constructor
}

var cookieType = reflect.TypeOf((*Cookie)(nil)).Elem()
var cookieTypeName = cookieType.Name()

const constructFlagFieldName = "isConstruct"

type emptyCookieConstructor struct {
	isConstruct bool
}

func (e *emptyCookieConstructor) Construct(r *http.Request) {
	e.isConstruct = true
}

func cookieConstruct(r *http.Request, cs Cookie) {
	if reflect.TypeOf(cs) == cookieType {
		return
	}
	cs.Construct(r)
	csv := reflect.ValueOf(cs)
	// true -> 表示使用的是emptyCookieConstructor，也就是默认的构造器。
	// 证明没有重新实现Construct方法，需要进行处理
	elem := csv.Elem()
	elem = elem.FieldByName(cookieTypeName).Elem()
	if elem.Elem().FieldByName(constructFlagFieldName).Bool() {
		cookies := r.Cookies()
		checkSet(stringFlag, cookieHolder(cookies), csv, cookieFiledName())
	}
}
