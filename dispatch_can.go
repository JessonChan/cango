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
		NewRouter(name string) forwarder
		Match(req *http.Request) matcher
	}

	forwarder interface {
		PathMethods(path string, ms ...string)
		GetName() string
	}

	matcher interface {
		Error() error
		Route() forwarder
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
		mapRouter  forwarder
		fastRouter forwarder
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

func (m *canDispatcher) NewRouter(name string) forwarder {
	return &canForwarder{mapRouter: m.mapMux.NewRouter(name), fastRouter: m.fastMux.NewRouter(name)}
}
func (m *canDispatcher) Match(req *http.Request) matcher {
	for _, v := range m.muxSlice {
		if cm := v.Match(req); cm.Error() == nil {
			return cm
		}
	}
	return &mapMatcher{err: errors.New("can dispatch can't find the path")}
}

func isVarPattern(path string) bool {
	return strings.Contains(path, "{") || strings.Contains(path, "*")
}

func (m *canForwarder) PathMethods(path string, ms ...string) {
	m.mapRouter.PathMethods(path, ms...)
	if isVarPattern(path) {
		m.fastRouter.PathMethods(path, ms...)
	}
}
func (m *canForwarder) GetName() string {
	return m.mapRouter.GetName()
}
