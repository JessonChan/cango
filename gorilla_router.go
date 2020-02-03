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
}

func (gm *gorillaMatcher) Error() error {
	return gm.RouteMatch.MatchErr
}
func (gm *gorillaMatcher) Route() CanRouter {
	return &gorillaRouter{gm.RouteMatch.Route}
}

func (gm *gorillaMatcher) GetVars() map[string][]string {
	return toValues(gm.Vars)
}

type gorillaMux struct {
	*mux.Router
}
type gorillaRouter struct {
	*mux.Route
}

func (gr *gorillaRouter) Path(path string) {
	gr.Route.Path(path)
}

func (gr *gorillaRouter) Methods(ms ...string) {
	gr.Route.Methods(ms...)
}

func newGorillaMux() *gorillaMux {
	return &gorillaMux{mux.NewRouter()}
}

func (gm *gorillaMux) NewRouter(name string) CanRouter {
	return &gorillaRouter{gm.Name(name)}
}
func (gm *gorillaMux) Match(req *http.Request) CanMatcher {
	match := &mux.RouteMatch{}
	gm.Router.Match(req, match)
	return &gorillaMatcher{match}
}
