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
			return nil
		}
		if strings.HasSuffix(path, ".tpl") {
			name := strings.TrimPrefix(path, can.tplRootPath)
			bs, _ := ioutil.ReadFile(path)
			_, _ = rootTpl.New(name).Funcs(can.tplFuncMap).Parse(string(bs))
		}
		return nil
	})
}

func (can *Can) lookupTpl(name string) *template.Template {
	tplOnce.Do(func() {
		can.initTpl()
	})
	return rootTpl.Lookup(name)
}
