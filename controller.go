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

type URI interface {
	Request() *WebRequest
}

var uriType = reflect.TypeOf((*URI)(nil)).Elem()
var uriName = uriType.Name()

type uriImpl struct {
	request *WebRequest
}

func (c *uriImpl) Request() *WebRequest {
	return c.request
}
func newContext(rw http.ResponseWriter, req *http.Request) *uriImpl {
	// todo sync.Pool
	return &uriImpl{&WebRequest{Request: req, ResponseWriter: rw}}
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
	return &RedirectWithCode{Redirect: &r, Code: code}
}

type RedirectWithCode struct {
	*Redirect
	Code int
}
type Content struct {
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

var defaultHttpMethods = []string{http.MethodGet}
var allHttpMethods = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}

/*
	methods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	for _, m := range methods {
		fmt.Printf("type %sMethod httpMethod\n", m)
	}
*/

type GetMethod httpMethod
type HeadMethod httpMethod
type PostMethod httpMethod
type PutMethod httpMethod
type PatchMethod httpMethod
type DeleteMethod httpMethod
type OptionsMethod httpMethod
type TraceMethod httpMethod

/*
	methods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	for _, m := range methods {
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

var httpMethodNameMap = map[string]string{
	"GetMethod":     "GET",
	"HeadMethod":    "HEAD",
	"PostMethod":    "POST",
	"PutMethod":     "PUT",
	"PatchMethod":   "PATCH",
	"DeleteMethod":  "DELETE",
	"OptionsMethod": "Options",
	"TraceMethod":   "Trace",
}
