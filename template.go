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
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var rootTpl *template.Template

var tplOnce sync.Once

func (can *Can) RegTplFunc(name string, fn interface{}) *Can {
	can.tplFuncMap[name] = fn
	return can
}

func (can *Can) initTpl() {
	rootTpl = template.New("")
	_ = filepath.Walk(can.tplRootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			canDebug("walk tpl files ", path, "dir skip")
			return nil
		}
		if func() bool {
			for _, suffix := range can.tplSuffix {
				if strings.HasSuffix(path, suffix) {
					return true
				}
			}
			return false
		}() {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				canError(err)
				return err
			}
			names := []string{info.Name(), strings.TrimPrefix(path, can.tplRootPath), strings.TrimPrefix(path, can.rootPath)}
			cnt := len(names)
			for i := 0; i < cnt; i++ {
				name := names[i]
				if strings.HasPrefix(name, "/") {
					names = append(names, name[1:])
				} else {
					names = append(names, "/"+name)
				}
			}
			for _, name := range names {
				if can.tplNameMap[name] {
					continue
				}
				can.tplNameMap[name] = true
				canDebug("walk tpl files ", path, name)
				_, err = rootTpl.New(name).Funcs(can.tplFuncMap).Parse(string(bs))
				if err != nil {
					canError(err)
					continue
				}
			}
		}
		return nil
	})
}

func (can *Can) lookupTpl(name string) *template.Template {
	if can.debugTpl {
		rootTpl = nil
		can.tplNameMap = map[string]bool{}
		can.initTpl()
	} else {
		tplOnce.Do(func() {
			can.initTpl()
		})
	}
	return rootTpl.Lookup(name)
}
