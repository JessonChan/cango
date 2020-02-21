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
	canMux struct {
		muxSlice []CanMux
		fastMux  *fastMux
		mapMux   *mapMux
	}
	canRouter struct {
		mapRouter  CanRouter
		fastRouter CanRouter
	}
)

func newCanMux() *canMux {
	mm := newMapMux()
	gm := newFastMux()
	return &canMux{
		mapMux:   mm,
		fastMux:  gm,
		muxSlice: []CanMux{mm, gm},
	}
}

func (m *canMux) NewRouter(name string) CanRouter {
	return &canRouter{mapRouter: m.mapMux.NewRouter(name), fastRouter: m.fastMux.NewRouter(name)}
}
func (m *canMux) Match(req *http.Request) CanMatcher {
	for _, v := range m.muxSlice {
		if cm := v.Match(req); cm.Error() == nil {
			return cm
		}
	}
	return &mapMatch{err: errors.New("can't find the path")}
}

func (m *canRouter) Path(ps ...string) {
	m.mapRouter.Path(ps...)
	m.fastRouter.Path(ps...)
}

func (m *canRouter) Methods(ms ...string) {
	m.mapRouter.Methods(ms...)
	m.fastRouter.Methods(ms...)
}
func (m *canRouter) GetName() string {
	return m.mapRouter.GetName()
}
func (m *canRouter) GetMethods() []string {
	return m.mapRouter.GetMethods()
}
func (m *canRouter) GetPath() string {
	return m.mapRouter.GetPath()
}
