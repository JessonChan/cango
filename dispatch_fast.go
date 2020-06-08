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
	"path/filepath"
	"strings"
)

type (
	fastDispatcher struct {
		forwarders   map[string]*fastForwarder
		forwarderArr []*fastForwarder
		patternMap   map[string]*fastPatten
	}
	fastForwarder struct {
		innerMux   *fastDispatcher
		name       string
		invoker    *Invoker
		pattens    []*fastPatten
		paths      []string
		methods    []string
		patternMap map[string]*fastPatten
	}
	fastPatten struct {
		forwarder     *fastForwarder
		pattern       string
		words         []word
		varIdx        []int
		methodMap     map[string]bool
		isWildcard    bool
		wildcardLeft  string
		wildcardRight string
	}
	fastMatcher struct {
		err       error
		forwarder *fastForwarder
		vars      map[string]string
	}
)

type word struct {
	key        string
	isVar      bool
	isWildcard bool
}

var _ = forwarder(&fastForwarder{})

func newFastMux() *fastDispatcher {
	return &fastDispatcher{forwarders: map[string]*fastForwarder{}, patternMap: map[string]*fastPatten{}}
}

func (fm *fastDispatcher) NewForwarder(name string, invoker *Invoker) forwarder {
	ff, ok := fm.forwarders[name]
	if ok {
		return ff
	}
	ff = &fastForwarder{innerMux: fm, name: name, invoker: invoker, patternMap: map[string]*fastPatten{}}
	fm.forwarders[name] = ff
	fm.forwarderArr = append(fm.forwarderArr, ff)
	return ff
}

func (fm *fastDispatcher) copy() map[string]*fastPatten {
	searchMap := make(map[string]*fastPatten, len(fm.patternMap))
	for k, v := range fm.patternMap {
		searchMap[k] = v
	}
	return searchMap
}

func (fm *fastDispatcher) doMatch(method, url string) *fastMatcher {
	elements := parsePath(url)
	searchMap := fm.copy()

	for k, elem := range elements {
		for key, pattern := range searchMap {
			if pattern.methodMap[method] == false {
				delete(searchMap, key)
			}
			// todo wildcard逻辑要加进来
			// todo /login/*.html 如果为 /log/*html则不行，通配符必须在两个分隔符之前
			// todo 为了简化，先约定只允许有一个通配符，并且通配符pattern不允许再有路径变量
			// todo 通配符逻辑统一
			if pattern.isWildcard {
				if strings.HasPrefix(url, pattern.wildcardLeft) && strings.HasSuffix(url, pattern.wildcardRight) {
					continue
				}
			} else {
				if len(pattern.words) != len(elements) {
					delete(searchMap, key)
					continue
				}
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
			delete(searchMap, key)
		}
	}
	switch len(searchMap) {
	case 0:
		// 没有找到
		return nil
	case 1:
		// 找到了
		var pattern *fastPatten
		for _, v := range searchMap {
			pattern = v
		}
		// never happen
		if pattern == nil {
			return nil
		}
		fm := &fastMatcher{
			forwarder: pattern.forwarder,
			vars:      make(map[string]string, len(pattern.varIdx)),
		}
		if len(pattern.varIdx) > 0 {
			for _, idx := range pattern.varIdx {
				fm.vars[pattern.words[idx].key] = elements[idx]
			}
		}
		return fm
	default:
		// 找到多个
		// todo 选最接近的
		// todo 这里可以随机选一个
		// panic("multi should handle")
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

func (fr *fastForwarder) PathMethods(path string, ms ...string) {
	patten, ok := fr.patternMap[path]
	if !ok {
		patten = &fastPatten{forwarder: fr, pattern: path, methodMap: map[string]bool{}}
		patten.words, patten.varIdx, patten.isWildcard = elementsToWords(parsePath(path))
		if patten.isWildcard {
			ss := strings.Split(path, "*")
			patten.wildcardLeft = ss[0]
			patten.wildcardRight = ss[1]
		}
		fr.pattens = append(fr.pattens, patten)
		fr.patternMap[path] = patten
	}
	for _, m := range ms {
		if patten.methodMap[m] == true {
			continue
		}
		patten.methodMap[m] = true
	}
	fr.innerMux.patternMap[fr.name] = patten
}
func (fr *fastForwarder) GetName() string {
	return fr.name
}
func (fr *fastForwarder) GetInvoker() *Invoker {
	return fr.invoker
}

func (fm *fastMatcher) Error() error {
	return fm.err
}
func (fm *fastMatcher) Forwarder() forwarder {
	return fm.forwarder
}
func (fm *fastMatcher) GetVars() map[string]string {
	return fm.vars
}

func elementsToWords(elements []string) ([]word, []int, bool) {
	words := make([]word, len(elements))
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
			words[i] = word{key: elem[1 : len(elem)-1], isVar: true}
			idx[j] = i
			j++
			continue
		}
		words[i] = word{key: elem, isVar: false}
	}
	return words, idx[0:j], isWildcard
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
