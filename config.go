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
	"io/ioutil"
	"reflect"
	"strings"
	"sync"

	"github.com/JessonChan/canlog"
)

var envs = map[string]string{}
var envOnce sync.Once

// Env 从配置文件中取配置项为key的值
func Env(key string) string {
	envOnce.Do(initIniConfig)
	return envs[key]
}

// Envs 对传入的对象进行赋值
// i 须为指针类型
func Envs(i interface{}) {
	decode(envs, reflect.ValueOf(i), func(field reflect.StructField) []string {
		return filedName(field, nameTagName)
	})
}

func initIniConfig() {
	configPath := getRootPath() + "/conf/cango.ini"
	for _, line := range strings.Split(func() string {
		bs, err := ioutil.ReadFile(configPath)
		if err != nil {
			canlog.CanError(err)
			return ""
		}
		return string(bs)
	}(), "\n") {
		idx := strings.Index(line, "=")
		if idx == -1 {
			canlog.CanError("can't parse config line", line)
			continue
		}
		line := strings.TrimSpace(line)
		commentIdx := strings.Index(line, "#")
		if commentIdx == 0 {
			continue
		}
		if commentIdx == -1 {
			commentIdx = len(line)
		}
		envs[strings.TrimSpace(line[:idx])] = strings.TrimSpace(line[idx+1 : commentIdx])
	}
}
