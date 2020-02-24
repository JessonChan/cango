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
	fastDispatcher struct {
		routers            map[string]*fastForwarder
		routerArr          []*fastForwarder
		methodRouterArrMap map[string][]*fastPatten
		pathNameMap        map[string]string
	}
	fastForwarder struct {
		innerMux *fastDispatcher
		name     string
		pattens  []*fastPatten
		paths    []string
		methods  []string
	}
	fastPatten struct {
		name          string
		pattern       string
		words         []word
		varIdx        []int
		hasVar        bool
		methodMap     map[string]bool
		isWildcard    bool
		wildcardLeft  string
		wildcardRight string
	}
	fastMatcher struct {
		err    error
		router *fastForwarder
		vars   map[string]string
	}
)

type word struct {
	key        string
	isVar      bool
	isWildcard bool
}

func newFastMux() *fastDispatcher {
	return &fastDispatcher{routers: map[string]*fastForwarder{}, methodRouterArrMap: map[string][]*fastPatten{}, pathNameMap: map[string]string{}}
}

func (fm *fastDispatcher) NewRouter(name string) forwarder {
	fr := &fastForwarder{innerMux: fm, name: name}
	fm.routers[name] = fr
	fm.routerArr = append(fm.routerArr, fr)
	return fr
}

func (fm *fastDispatcher) doMatch(method, url string) *fastMatcher {
	elements := parsePath(url)
	stopMap := map[*fastPatten]bool{}

	for k, elem := range elements {
		for _, pattern := range fm.methodRouterArrMap[method] {
			if stopMap[pattern] {
				continue
			}
			// todo wildcard逻辑要加进来
			// todo /login/*.html 如果为 /log/*html则不行，通配符必须在两个分隔符之前
			// todo 为了简化，先约定只允许有一个通配符，并且通配符pattern不允许再有路径变量
			// todo 通配符逻辑统一
			if pattern.isWildcard {
				// todo 效率提升
				if strings.HasPrefix(url, pattern.wildcardLeft) && strings.HasSuffix(url, pattern.wildcardRight) {
					continue
				}
			}
			if len(pattern.words) != len(elements) {
				stopMap[pattern] = true
				continue
			}

			// 如果是变量，肯定符合
			if pattern.words[k].isVar {
				continue
			}
			// 如果不是变量，判断是不是相等
			if pattern.words[k].key == elem {
				continue
			}
			// 两种情况都不是，不符合
			stopMap[pattern] = true
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
			vars:   vars,
		}
		return fm
	default:
		// 找到多个
		// todo 选最接近的
		// todo 这里可以随机选一个
		panic("multi should handle")
		return nil
	}
}
func (fm *fastDispatcher) Match(req *http.Request) matcher {
	cm := fm.doMatch(req.Method, req.URL.Path)
	if cm == nil {
		return &fastMatcher{err: errors.New("fast dispatch can't find the path")}
	}
	return cm
}

func (fr *fastForwarder) Path(ps ...string) {
	for k, path := range ps {
		routerName := fmt.Sprintf("%v-%d", fr.name, k)
		fr.innerMux.pathNameMap[routerName] = fr.name
		patten := &fastPatten{name: routerName, pattern: path}
		patten.words, patten.varIdx, patten.hasVar, patten.isWildcard = elementsToWords(parsePath(path))
		if patten.isWildcard {
			ss := strings.Split(path, "*")
			patten.wildcardLeft = ss[0]
			patten.wildcardRight = ss[1]
		}
		fr.pattens = append(fr.pattens, patten)
	}
	fr.paths = ps
}

func (fr *fastForwarder) Methods(ms ...string) {
	for _, v := range ms {
		fr.innerMux.methodRouterArrMap[v] = append(fr.innerMux.methodRouterArrMap[v], fr.pattens...)
	}
	fr.methods = ms
}
func (fr *fastForwarder) GetName() string {
	return fr.name
}
func (fr *fastForwarder) GetMethods() []string {
	return fr.methods
}
func (fr *fastForwarder) GetPath() string {
	return strings.Join(fr.paths, ";")
}

func (fm *fastMatcher) Error() error {
	return fm.err
}
func (fm *fastMatcher) Route() forwarder {
	return fm.router
}
func (fm *fastMatcher) GetVars() map[string]string {
	return fm.vars
}

func elementsToWords(elements []string) ([]word, []int, bool, bool) {
	words := make([]word, len(elements))
	hasVar := false
	idx := make([]int, len(elements))
	isWildcard := false
	j := 0
	for i, elem := range elements {
		if elem == "*" {
			words[i] = word{isWildcard: true}
			isWildcard = true
			continue
		}
		// todo 如果这个word就是{}呢？
		// todo 也就是说地址是/a/{}/b/c 这种的话不会被当做变量
		// todo 如果真实需要注册的地址就是/a/{name}/b/c 应该怎么办？
		if strings.HasPrefix(elem, "{") && strings.HasSuffix(elem, "}") {
			hasVar = true
			words[i] = word{key: elem[1 : len(elem)-1], isVar: true}
			idx[j] = i
			j++
			continue
		}
		words[i] = word{key: elem, isVar: false}
	}
	return words, idx, hasVar, isWildcard
}

func parsePath(url string) []string {
	// todo : clean & split in only one loop
	url = filepath.Clean(url)
	// todo : 如果是 /a/b/c/ 是不是要等价于 /a/b/c 呢？
	if url == "" || url == "/" {
		return []string{"/"}
	}
	// todo .这个分隔符应该是最后一个才需要，并且应该是生成两个pattern
	return strings.FieldsFunc(url, func(r rune) bool {
		if r == '/' || r == '.' {
			return true
		}
		return false
	})
}
