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
	mm := newMapMux()
	gm := newFastMux()
	return &canDispatcher{
		mapMux:   mm,
		fastMux:  gm,
		muxSlice: []dispatcher{mm, gm},
	}
}

func isVarPattern(path string) bool {
	return strings.Contains(path, "{") || strings.Contains(path, "*")
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
