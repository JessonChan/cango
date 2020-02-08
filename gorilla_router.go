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
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type (
	gorillaMux struct {
		*mux.Router
		routerMap map[string]*gorillaRouter
		// map path to name
		pathNameMap map[string]string
	}
	gorillaRouter struct {
		innerMux *gorillaMux
		name     string
		routes   []*mux.Route
		path     []string
		methods  []string
	}
	gorillaMatcher struct {
		*mux.RouteMatch
		*gorillaRouter
	}
)

func (gm *gorillaMatcher) Error() error {
	return gm.RouteMatch.MatchErr
}
func (gm *gorillaMatcher) Route() CanRouter {
	return gm.gorillaRouter
}

func (gm *gorillaMatcher) GetVars() map[string][]string {
	return toValues(gm.Vars)
}

func (gr *gorillaRouter) Path(paths ...string) {
	for k, path := range paths {
		routerName := fmt.Sprintf("%s-%d", gr.name, k)
		route := gr.innerMux.Name(routerName)
		route.Path(path)
		gr.routes = append(gr.routes, route)
		gr.innerMux.pathNameMap[routerName] = gr.name
	}
	gr.path = paths
}
func (gr *gorillaRouter) Methods(ms ...string) {
	for _, r := range gr.routes {
		r.Methods(ms...)
	}
	gr.methods = ms
}

func (gr *gorillaRouter) GetName() string {
	return gr.name
}
func (gr *gorillaRouter) GetPath() string {
	return strings.Join(gr.path, ";")
}
func (gr *gorillaRouter) GetMethods() []string {
	return gr.methods
}

func newGorillaMux() *gorillaMux {
	return &gorillaMux{Router: mux.NewRouter(), routerMap: map[string]*gorillaRouter{}, pathNameMap: map[string]string{}}
}

func (gm *gorillaMux) NewRouter(name string) CanRouter {
	gr := &gorillaRouter{name: name, innerMux: gm}
	gm.routerMap[name] = gr
	return gr
}
func (gm *gorillaMux) Match(req *http.Request) CanMatcher {
	match := &mux.RouteMatch{}
	gm.Router.Match(req, match)
	return &gorillaMatcher{match, gm.routerMap[gm.pathNameMap[match.Route.GetName()]]}
}
