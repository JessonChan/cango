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
	mapDispatcher struct {
		routers  map[string]*mapForwarder
		pathName map[string]string
	}
	mapForwarder struct {
		name       string
		innerMux   *mapDispatcher
		paths      []string
		patternMap map[string]*mapPattern
	}
	mapMatcher struct {
		innerRouter *mapForwarder
		err         error
	}
	//  每个注册都是由一条路径和其对应的方法组成
	//  如 /api/movie/list.json [GET][POST]
	mapPattern struct {
		path    string
		methods map[string]bool
	}
)

func newMapMux() *mapDispatcher {
	return &mapDispatcher{routers: map[string]*mapForwarder{}, pathName: map[string]string{}}
}

func (m *mapDispatcher) NewRouter(name string) forwarder {
	mr, ok := m.routers[name]
	if ok {
		return mr
	}
	mr = &mapForwarder{name: name, innerMux: m, patternMap: map[string]*mapPattern{}}
	m.routers[name] = mr
	return mr
}
func (m *mapDispatcher) Match(req *http.Request) matcher {
	r, ok := m.routers[m.pathName[req.URL.Path]]
	if ok && r.patternMap[req.URL.Path].methods[req.Method] {
		return &mapMatcher{innerRouter: r}
	}
	return &mapMatcher{err: errors.New("map dispatch can't find the path")}
}

func (m *mapForwarder) Path(ps ...string) {
	for _, p := range ps {
		m.innerMux.pathName[p] = m.name
	}
	m.paths = ps
}

func (m *mapForwarder) Methods(ms ...string) {
	for _, path := range m.paths {
		pattern := m.patternMap[path]
		if pattern == nil {
			pattern = &mapPattern{path: path}
			pattern.methods = map[string]bool{}
		}
		for _, m := range ms {
			pattern.methods[m] = true
		}
		m.patternMap[path] = pattern
	}
}
func (m *mapForwarder) GetName() string {
	return m.name
}
func (m *mapForwarder) GetMethods() []string {
	return []string{}
}
func (m *mapForwarder) GetPath() string {
	return strings.Join(m.paths, ";")
}

func (m *mapMatcher) Error() error {
	return m.err
}

func (m *mapMatcher) Route() forwarder {
	return m.innerRouter
}
func (m *mapMatcher) GetVars() map[string]string {
	return map[string]string{}
}
