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

type IniConfig struct {
	envs    map[string]string
	envForm map[string][]string
}

var configOnce sync.Once

var canIniConfig *IniConfig

func NewIniConfig(filePath string) *IniConfig {
	return initIniConfig(filePath)
}

// Env 从配置文件中取配置项为key的值
func Env(key string) string {
	configOnce.Do(initCangoIni)
	return canIniConfig.Env(key)
}
func initCangoIni() {
	canIniConfig = NewIniConfig(getRootPath() + "/conf/cango.ini")
}

// Envs 对传入的对象进行赋值
// i 须为指针类型
func Envs(i interface{}) {
	configOnce.Do(initCangoIni)
	canIniConfig.Envs(i)
}

// Env 从配置文件中取配置项为key的值
func (ic *IniConfig) Env(key string) string {
	return ic.envs[key]
}

// Envs 对传入的对象进行赋值
// i 须为指针类型
func (ic *IniConfig) Envs(i interface{}) {
	decodeForm(ic.envForm, reflect.ValueOf(i), func(field reflect.StructField) []string {
		return filedName(field, nameTagName)
	})
}

func initIniConfig(configPath string) *IniConfig {
	var envs = map[string]string{}
	var envForm = map[string][]string{}
	for _, line := range strings.Split(func() string {
		bs, err := ioutil.ReadFile(configPath)
		if err != nil {
			return ""
		}
		return string(bs)
	}(), "\n") {
		if line == "" {
			continue
		}
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
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		envs[key] = trimQuote(value)
		// 数组变量
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			strs := strings.Split(value[1:len(value)-1], ",")
			for i := 0; i < len(strs); i++ {
				strs[i] = trimQuote(strs[i])
			}
			envForm[key] = strs
		} else {
			envForm[key] = []string{trimQuote(value)}
		}
	}
	return &IniConfig{
		envs:    envs,
		envForm: envForm,
	}
}

func trimQuote(str string) string {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		return str[1 : len(str)-1]
	}
	if strings.HasPrefix(str, `'`) && strings.HasSuffix(str, `'`) {
		return str[1 : len(str)-1]
	}
	return str
}
