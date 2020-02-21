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
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

type (
	fastMux struct {
		routers            map[string]*fastRouter
		routerArr          []*fastRouter
		methodRouterArrMap map[string][]*fastPatten
		pathNameMap        map[string]string
	}
	fastRouter struct {
		innerMux *fastMux
		name     string
		pattens  []*fastPatten
		paths    []string
		methods  []string
	}
	fastPatten struct {
		name      string
		pattern   string
		words     []word
		varIdx    []int
		hasVar    bool
		methodMap map[string]bool
	}
	fastMatcher struct {
		err    error
		router *fastRouter
		vars   map[string][]string
	}
)

type word struct {
	key   string
	isVar bool
}

func newFastMux() *fastMux {
	return &fastMux{routers: map[string]*fastRouter{}, methodRouterArrMap: map[string][]*fastPatten{}, pathNameMap: map[string]string{}}
}

func (fm *fastMux) NewRouter(name string) CanRouter {
	fr := &fastRouter{innerMux: fm, name: name}
	fm.routers[name] = fr
	fm.routerArr = append(fm.routerArr, fr)
	return fr
}

func (fm *fastMux) doMatch(method, url string) *fastMatcher {
	elements := parsePath(url)
	stopMap := map[*fastPatten]bool{}

	for k, elem := range elements {
		for _, router := range fm.methodRouterArrMap[method] {
			if stopMap[router] {
				continue
			}
			// len(elements) 和 len(pattSeg) 如果在不是变量情况下，肯定是一样长的
			// 如果存在变量情况下呢？

			if len(router.words) != len(elements) {
				stopMap[router] = true
				continue
			}

			// 如果是变量，肯定符合
			if router.words[k].isVar {
				continue
			}
			// 如果不是变量，判断是不是相等
			if router.words[k].key == elem {
				continue
			}
			// 两种情况都不是，不符合
			stopMap[router] = true
		}
	}
	diff := len(fm.methodRouterArrMap[method]) - len(stopMap)
	switch diff {
	case 0:
		// 没有找到
		return nil
	case 1:
		// 找到了
		var router *fastPatten
		for _, v := range fm.methodRouterArrMap[method] {
			if stopMap[v] {
				continue
			}
			router = v
			break
		}
		if router == nil {
			return nil
		}
		vars := map[string]string{}
		for _, id := range router.varIdx {
			vars[router.words[id].key] = elements[id]
		}
		fm := &fastMatcher{
			router: fm.routers[fm.pathNameMap[router.name]],
			vars:   toValues(vars),
		}
		return fm
	default:
		// 找到多个
		// todo 选最接近的
		// todo 这里可以随机选一个
		return nil
	}
}
func (fm *fastMux) Match(req *http.Request) CanMatcher {
	cm := fm.doMatch(req.Method, req.URL.Path)
	if cm == nil {
		return &fastMatcher{err: errors.New("can't find the path")}
	}
	return cm
}

func (fr *fastRouter) Path(ps ...string) {
	for k, path := range ps {
		routerName := fmt.Sprintf("%v-%d", fr.name, k)
		fr.innerMux.pathNameMap[routerName] = fr.name
		patten := &fastPatten{name: routerName, pattern: path}
		patten.words, patten.varIdx, patten.hasVar = elementsToWords(parsePath(path))
		fr.pattens = append(fr.pattens, patten)
	}
	fr.paths = ps
}

func (fr *fastRouter) Methods(ms ...string) {
	for _, v := range ms {
		fr.innerMux.methodRouterArrMap[v] = append(fr.innerMux.methodRouterArrMap[v], fr.pattens...)
	}
	fr.methods = ms
}
func (fr *fastRouter) GetName() string {
	return fr.name
}
func (fr *fastRouter) GetMethods() []string {
	return fr.methods
}
func (fr *fastRouter) GetPath() string {
	return strings.Join(fr.paths, ";")
}

func (fm *fastMatcher) Error() error {
	return fm.err
}
func (fm *fastMatcher) Route() CanRouter {
	return fm.router
}
func (fm *fastMatcher) GetVars() map[string][]string {
	return fm.vars
}

func elementsToWords(elements []string) ([]word, []int, bool) {
	words := make([]word, len(elements))
	hasVar := false
	idx := make([]int, len(elements))
	j := 0
	for i := 0; i < len(elements); i++ {
		if []byte(elements[i])[0] == '{' && []byte(elements[i])[len(elements[i])-1] == '}' {
			hasVar = true
			words[i] = word{key: elements[i][1 : len(elements[i])-1], isVar: true}
			idx[j] = i
			j++
			continue
		}
		words[i] = word{key: elements[i], isVar: false}
	}
	return words, idx, hasVar
}

func parsePath(url string) []string {
	// todo : clean & split in only one loop
	url = filepath.Clean(url)
	// todo : 如果是 /a/b/c/ 是不是要等价于 /a/b/c 呢？
	if url == "" || url == "/" {
		return []string{"/"}
	}
	return strings.FieldsFunc(url, func(r rune) bool {
		if r == '/' || r == '.' {
			return true
		}
		return false
	})
}
