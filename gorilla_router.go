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

	"github.com/gorilla/mux"
)

type gorillaMatcher struct {
	*mux.RouteMatch
	CanRouter
}

func (gm *gorillaMatcher) Error() error {
	return gm.RouteMatch.MatchErr
}
func (gm *gorillaMatcher) Route() CanRouter {
	return gm.CanRouter
}

func (gm *gorillaMatcher) GetVars() map[string][]string {
	return toValues(gm.Vars)
}

type gorillaMux struct {
	*mux.Router
	routerMap map[*mux.Route]CanRouter
}
type gorillaRouter struct {
	*mux.Route
	path    string
	methods []string
}

func (gr *gorillaRouter) Path(path string) {
	gr.path = path
	gr.Route.Path(path)
}
func (gr *gorillaRouter) GetPath() string {
	return gr.path
}
func (gr *gorillaRouter) GetMethods() []string {
	return gr.methods
}

func (gr *gorillaRouter) Methods(ms ...string) {
	gr.Route.Methods(ms...)
	gr.methods = ms
}

func newGorillaMux() *gorillaMux {
	return &gorillaMux{Router: mux.NewRouter(), routerMap: map[*mux.Route]CanRouter{}}
}

func (gm *gorillaMux) NewRouter(name string) CanRouter {
	gr := &gorillaRouter{Route: gm.Name(name)}
	gm.routerMap[gr.Route] = gr
	return gr
}
func (gm *gorillaMux) Match(req *http.Request) CanMatcher {
	match := &mux.RouteMatch{}
	gm.Router.Match(req, match)
	return &gorillaMatcher{match, gm.routerMap[match.Route]}
}
