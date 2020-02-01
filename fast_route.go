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
	"path/filepath"
	"strings"
)

type fastMux struct {
	routers map[string]*fastRouter
}

type fastRouter struct {
	name      string
	pattern   string
	seg       []word
	methodMap map[string]bool
}

type word struct {
	key   string
	isVar bool
}

func newMux() *fastMux {
	return &fastMux{routers: map[string]*fastRouter{}}
}

func (fm *fastMux) newRouter(name string) *fastRouter {
	fr := &fastRouter{name: name, seg: []word{}, methodMap: map[string]bool{}}
	fm.routers[name] = fr
	return fr
}

type routerPair struct {
	router *fastRouter
	vars   map[string][]string
}

func (fm *fastMux) match(method string, url string) *routerPair {
	pair := &routerPair{vars: map[string][]string{}}
	frs := make([]*fastRouter, len(fm.routers))
	seg := parsePath(url)
	idx := 0
	// method 要对应上
	for _, fr := range fm.routers {
		// todo fix: url 和 pattern 有可能包含不一样的seg
		if fr.methodMap[method] && len(seg) == len(fr.seg) {
			frs[idx] = fr
			idx++
		}
	}
	if idx == 0 {
		return nil
	}
	for i := 0; i < len(frs); i++ {
		fr := frs[i]
		find := false
		for j := 0; j < len(seg); j++ {
			if fr.seg[j].isVar {
				pair.vars[fr.seg[j].key] = []string{seg[j]}
				find = true
				continue
			}
			if fr.seg[j].key == seg[j] {
				find = true
				continue
			}
			find = false
		}
		if find {
			pair.router = fr
			return pair
		}
	}
	return nil
}

func (fr *fastRouter) path(pattern string) {
	fr.pattern = pattern
	fr.seg = pathToWord(parsePath(pattern))
}

func (fr *fastRouter) methods(ms ...string) {
	for _, v := range ms {
		fr.methodMap[v] = true
	}
}

func pathToWord(seg []string) []word {
	words := make([]word, len(seg))
	for i := 0; i < len(seg); i++ {
		if []byte(seg[i])[0] == '{' && []byte(seg[i])[len(seg[i])-1] == '}' {
			words[i] = word{key: seg[i][1 : len(seg[i])-1], isVar: true}
			continue
		}
		words[i] = word{key: seg[i], isVar: false}
	}
	return words
}

func parsePath(url string) []string {
	// todo : clean & split in only one loop
	url = filepath.Clean(url)
	return strings.FieldsFunc(url, func(r rune) bool {
		if r == '/' || r == '.' {
			return true
		}
		return false
	})
}
