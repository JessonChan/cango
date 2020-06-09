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
)

type (
	mapDispatcher struct {
		forwarders map[string]*mapForwarder
		// pathName 用来标记路由和对应的处理标识之间的关系
		pathName map[string][]string
	}
	mapForwarder struct {
		name       string
		invoker    *Invoker
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
	return &mapDispatcher{forwarders: map[string]*mapForwarder{}, pathName: map[string][]string{}}
}

func (m *mapDispatcher) NewForwarder(name string, invoker *Invoker) forwarder {
	mr, ok := m.forwarders[name]
	if ok {
		return mr
	}
	mr = &mapForwarder{name: name, innerMux: m, invoker: invoker, patternMap: map[string]*mapPattern{}}
	m.forwarders[name] = mr
	return mr
}
func (m *mapDispatcher) Match(req *http.Request) matcher {
	names := m.pathName[req.URL.Path]
	for _, name := range names {
		r, ok := m.forwarders[name]
		if ok && r.patternMap[req.URL.Path].methods[req.Method] {
			return &mapMatcher{innerRouter: r}
		}
	}
	return &mapMatcher{err: errors.New("map dispatch can't find the path")}
}

func (m *mapForwarder) PathMethods(path string, ms ...string) {
	names := m.innerMux.pathName[path]
	if len(names) == 0 {
		m.innerMux.pathName[path] = []string{m.name}
	} else {
		contains := false
		for _, v := range names {
			if v == m.name {
				contains = true
				break
			}
		}
		if contains {
		} else {
			m.innerMux.pathName[path] = append(names, m.name)
		}
	}
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

func (m *mapForwarder) GetInvoker() *Invoker {
	return m.invoker
}

func (m *mapMatcher) Error() error {
	return m.err
}

func (m *mapMatcher) Forwarder() forwarder {
	return m.innerRouter
}

func (m *mapMatcher) GetVars() map[string]string {
	return map[string]string{}
}
