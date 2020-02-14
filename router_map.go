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
	mapMux struct {
		routers  map[string]*mapRouter
		pathName map[string]string
	}
	mapRouter struct {
		name        string
		innerMux    *mapMux
		paths       []string
		methods     map[string]bool
		methodSlice []string
	}
	mapMatch struct {
		innerRouter *mapRouter
		err         error
	}
)

func newMapMux() *mapMux {
	return &mapMux{routers: map[string]*mapRouter{}, pathName: map[string]string{}}
}

func (m *mapMux) NewRouter(name string) CanRouter {
	mr := &mapRouter{name: name, innerMux: m, methods: map[string]bool{}}
	m.routers[name] = mr
	return mr
}
func (m *mapMux) Match(req *http.Request) CanMatcher {
	r, ok := m.routers[m.pathName[req.URL.Path]]
	if ok {
		return &mapMatch{innerRouter: r}
	}
	return &mapMatch{err: errors.New("can't find the path")}
}

func (m *mapRouter) Path(ps ...string) {
	for _, p := range ps {
		m.innerMux.pathName[p] = m.name
	}
	m.paths = ps
}

func (m *mapRouter) Methods(ms ...string) {
	for _, v := range ms {
		m.methods[v] = true
	}
	m.methodSlice = ms
}
func (m *mapRouter) GetName() string {
	return m.name
}
func (m *mapRouter) GetMethods() []string {
	return m.methodSlice
}
func (m *mapRouter) GetPath() string {
	return strings.Join(m.paths, ",")
}

func (m *mapMatch) Error() error {
	return m.err
}

func (m *mapMatch) Route() CanRouter {
	return m.innerRouter
}
func (m *mapMatch) GetVars() map[string][]string {
	return map[string][]string{}
}
