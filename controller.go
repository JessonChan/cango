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

type WebRequest struct {
	http.ResponseWriter
	*http.Request
}

/*
	httpMethods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	for _, m := range httpMethods {
		fmt.Println("// Is"+m+" check the request method is http.Method"+m)
		fmt.Println("func (wr *WebRequest) Is"+m+"() bool {")
		fmt.Println("\treturn wr.Method == http.Method"+m)
		fmt.Println("}")
	}
*/

// IsGet check the request method is http.MethodGet
func (wr *WebRequest) IsGet() bool {
	return wr.Method == http.MethodGet
}

// IsPost check the request method is http.MethodPost
func (wr *WebRequest) IsPost() bool {
	return wr.Method == http.MethodPost
}

// IsHead check the request method is http.MethodHead
func (wr *WebRequest) IsHead() bool {
	return wr.Method == http.MethodHead
}

// IsPut check the request method is http.MethodPut
func (wr *WebRequest) IsPut() bool {
	return wr.Method == http.MethodPut
}

// IsPatch check the request method is http.MethodPatch
func (wr *WebRequest) IsPatch() bool {
	return wr.Method == http.MethodPatch
}

// IsDelete check the request method is http.MethodDelete
func (wr *WebRequest) IsDelete() bool {
	return wr.Method == http.MethodDelete
}

// IsOptions check the request method is http.MethodOptions
func (wr *WebRequest) IsOptions() bool {
	return wr.Method == http.MethodOptions
}

// IsTrace check the request method is http.MethodTrace
func (wr *WebRequest) IsTrace() bool {
	return wr.Method == http.MethodTrace
}

type URI interface {
	Request() *WebRequest
}

var uriType = reflect.TypeOf((*URI)(nil)).Elem()
var uriName = uriType.Name()
var uriMethods = func() (ms []reflect.Method) {
	ut := reflect.TypeOf(&uriImpl{})
	for i := 0; i < ut.NumMethod(); i++ {
		ms = append(ms, ut.Method(i))
	}
	return
}()

// 判断某个方法是否是接口中包含的方法，接口中的方法不参与路由的构建
func uriInterfaceContains(m reflect.Method) bool {
	for _, v := range uriMethods {
		if v.Name == m.Name && m.Type.NumIn() == v.Type.NumIn() && m.Type.NumOut() == v.Type.NumOut() {
			return true
		}
	}
	return false
}

type uriImpl struct {
	request *WebRequest
}

func (c *uriImpl) Request() *WebRequest {
	return c.request
}
func newContext(request *WebRequest) *uriImpl {
	// todo sync.Pool
	return &uriImpl{request: request}
}

// todo  可以根据路由的uri/方法名等自动查找tpl
type ModelView struct {
	Tpl   string
	Model interface{}
}

type Redirect struct {
	Url string
}

func (r Redirect) WithCode(code int) *RedirectWithCode {
	return &RedirectWithCode{Url: r.Url, Code: code}
}

type RedirectWithCode struct {
	Url  string
	Code int
}
type Content struct {
	String string
}

func (c Content) WithCode(code int) *ContentWithCode {
	return &ContentWithCode{String: c.String, Code: code}
}

type ContentWithCode struct {
	String string
	Code   int
}

type StaticFile struct {
	Path string
}

type DoNothing struct {
}

var doNothing = DoNothing{}

type httpMethod interface {
}

var defaultHTTPMethods = []string{http.MethodGet}
var allHTTPMethods = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}

/*
	httpMethods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	fmt.Println("type (")
	for _, m := range httpMethods {
		fmt.Printf("\t%sMethod httpMethod\n", m)
	}
	fmt.Println(")")
*/

// Common HTTP methods.
type (
	GetMethod     httpMethod
	HeadMethod    httpMethod
	PostMethod    httpMethod
	PutMethod     httpMethod
	PatchMethod   httpMethod
	DeleteMethod  httpMethod
	OptionsMethod httpMethod
	TraceMethod   httpMethod
)

/*
	httpMethods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	for _, m := range httpMethods {
		fmt.Printf("	reflect.TypeOf((*%sMethod)(nil)).Elem(): http.Method%s,\n", m, m)
	}
*/

var httpMethodMap = map[reflect.Type]string{
	reflect.TypeOf((*GetMethod)(nil)).Elem():     http.MethodGet,
	reflect.TypeOf((*PostMethod)(nil)).Elem():    http.MethodPost,
	reflect.TypeOf((*HeadMethod)(nil)).Elem():    http.MethodHead,
	reflect.TypeOf((*PutMethod)(nil)).Elem():     http.MethodPut,
	reflect.TypeOf((*PatchMethod)(nil)).Elem():   http.MethodPatch,
	reflect.TypeOf((*DeleteMethod)(nil)).Elem():  http.MethodDelete,
	reflect.TypeOf((*OptionsMethod)(nil)).Elem(): http.MethodOptions,
	reflect.TypeOf((*TraceMethod)(nil)).Elem():   http.MethodTrace,
}
