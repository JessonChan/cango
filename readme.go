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
	"strings"
)

func optsToTable() {
	// 需要更新时进行替换
	lines := strings.Split(`	// 监听的主机
	Host string
	// 监听的端口
	Port int
	// 配置文件、模板文件及静态文件的根目录
	RootPath string
	// 模板文件文件夹，相对RootPath路径
	TplDir string
	// 静态文件文件夹，相对RootPath路径
	StaticDir string
	// 模板文件后续名，默认为 .tpl 和 .html
	TplSuffix []string
	// 是否调试页面，true 表示每次都重新加载模板
	DebugTpl bool
	// 日志文件位置，绝对路径
	CanlogPath string
	// gorilla cookie store 的key
	CookieSessionKey string
	// gorilla cookie store 加密使用的key
	CookieSessionSecure string`, "\n")
	if len(lines)%2 != 0 {
		panic("必须是偶数")
	}

	lowerCase := func(s string) string {
		bs := []rune(s)
		if 'A' <= bs[0] && bs[0] <= 'Z' {
			bs[0] += 'a' - 'A'
		}
		return string(bs)
	}
	underScore := func(s string) string {
		bs := make([]rune, 0, 2*len(s))
		for _, s := range s {
			if 'A' <= s && s <= 'Z' {
				s += 'a' - 'A'
				bs = append(bs, '_')
			}
			bs = append(bs, s)
		}
		if bs[0] == '_' {
			s = string(bs[1:])
		} else {
			s = string(bs)
		}
		return s
	}

	names := func(name string) string {
		l := lowerCase(name)
		u := underScore(name)
		s := u
		if l != u {
			s += "/" + l
		}
		if name != u {
			s += "/" + name
		}
		return s
	}

	format := func(words []string) {
		s := "| "
		for i := 0; i < len(words); i++ {
			word := words[i]
			if i == 0 {
				word = names(words[i])
			}
			s += word + " |"
		}
		fmt.Println(s)
	}
	format([]string{"参数名", "参数类型", "参数说明"})
	format([]string{"----", "----", "----"})
	for i := 0; i < len(lines); i = i + 2 {
		j := i + 1
		key := strings.TrimSpace(lines[j])
		comment := strings.TrimSpace(lines[i])
		format(append(strings.Split(key, " "), strings.Replace(comment, "// ", "", 1)))
	}
}
