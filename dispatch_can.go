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
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	dispatcher interface {
		NewForwarder(name string, invoker *Invoker) forwarder
		Match(req *http.Request) matcher
	}

	forwarder interface {
		PathMethods(path string, ms ...string)
		GetInvoker() *Invoker
	}

	matcher interface {
		Error() error
		Forwarder() forwarder
		GetVars() map[string]string
	}
)

type (
	canDispatcher struct {
		muxSlice []dispatcher
		fastMux  *fastDispatcher
		mapMux   *mapDispatcher
	}
	canForwarder struct {
		mapForwarder  forwarder
		fastForwarder forwarder
	}
)

func newCanMux() *canDispatcher {
	mm := newMapDispatcher()
	fm := newFastDispatcher()
	return &canDispatcher{
		mapMux:   mm,
		fastMux:  fm,
		muxSlice: []dispatcher{mm, fm},
	}
}

func isVarPattern(path string) bool {
	return strings.Contains(path, "{") || strings.Contains(path, "*")
}

// Gin's path parameters pattern is /:name, Cango's is /{name}
var replacer = strings.NewReplacer("{", ":", "}", "")

type JSON struct {
	Data interface{}
}

// Render (JSON) writes data with custom ContentType.
func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

var jsonContentType = []string{"application/json; charset=utf-8"}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := responseJsonHandler(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func (m *canDispatcher) Gins() (ghs []*GinHandler) {
	for _, forwarder := range m.mapMux.forwarders {
		for _, pattern := range forwarder.patternMap {
			gh := &GinHandler{
				Url: func() string {
					if strings.Contains(pattern.path, "{") {
						return replacer.Replace(pattern.path)
					}
					return pattern.path
				}(),
				Handle: func(ctx *gin.Context) {
					handleReturn, code := serve(m, &WebRequest{
						ResponseWriter: ctx.Writer,
						Request:        ctx.Request,
					})
					switch hr := handleReturn.(type) {
					case Content:
						ctx.String(code, hr.String)
					case StaticFile:
						ctx.File(hr.Path)
					case Redirect:
						ctx.Redirect(code, hr.Url)
					case ContentWithCode:
						ctx.String(hr.Code, hr.String)
					default:
						ctx.Render(code, JSON{Data: handleReturn})
					}
				},
				HttpMethods: pattern.methods,
			}
			ghs = append(ghs, gh)
		}
	}
	return
}

func (m *canDispatcher) NewForwarder(name string, invoker *Invoker) forwarder {
	return &canForwarder{mapForwarder: m.mapMux.NewForwarder(name, invoker), fastForwarder: m.fastMux.NewForwarder(name, invoker)}
}
func (m *canDispatcher) Match(req *http.Request) matcher {
	for _, v := range m.muxSlice {
		if cm := v.Match(req); cm.Error() == nil {
			return cm
		}
	}
	return &mapMatcher{err: errors.New("can dispatch can't find the path")}
}

func (m *canForwarder) PathMethods(path string, ms ...string) {
	m.mapForwarder.PathMethods(path, ms...)
	if isVarPattern(path) {
		m.fastForwarder.PathMethods(path, ms...)
	}
}

func (m *canForwarder) GetInvoker() *Invoker {
	return nil
}
