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
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// todo 缓存v struct结构
func decode(holder map[string]string, v interface{}, filedName func(reflect.StructField) []string) {
	checkSet(stringFlag, func(s string) (interface{}, bool) {
		v, ok := holder[s]
		return v, ok
	}, v, filedName)
}
func decodeForm(holder map[string][]string, v interface{}, filedName func(field reflect.StructField) []string) {
	checkSet(strSliceFlag, func(s string) (interface{}, bool) {
		v, ok := holder[s]
		return v, ok
	}, v, filedName)
}

func checkSet(flag int, holder func(string) (interface{}, bool), v interface{}, filedName func(field reflect.StructField) []string) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = reflect.Indirect(rv)
	setValue(flag, holder, rv, filedName)
}

const (
	stringFlag   = 0
	strSliceFlag = 1
)

func setValue(flag int, holder func(string) (interface{}, bool), rv reflect.Value, filedName func(field reflect.StructField) []string) {
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		fmt.Printf("%#v,%v\n", f, rv.Type().Field(i))
		if f.Kind() == reflect.Ptr {
			f = reflect.Indirect(f)
		}
		if f.Type().Kind() == reflect.Struct && f.Type() != timeType {
			setValue(flag, holder, f, filedName)
		}
		// 返回值表示是否找到对应的converter
		name := filedName(rv.Type().Field(i))
		set := func(flag int) bool {
			kind := f.Kind()
			if f.Type() == timeType {
				kind = timeTypeKind
			}
			if kind == reflect.Slice {
				return false
			}
			if converter, ok := converters[kind]; ok {
				for _, key := range name {
					if str, ok := holder(key); ok {
						switch flag {
						case stringFlag:
							f.Set(converter(str.(string)))
						case strSliceFlag:
							f.Set(converter(str.([]string)[0]))
						}
						return true
					}
				}
				return true
			}
			return false
		}
		if f.CanSet() {
			switch flag {
			case stringFlag:
				set(stringFlag)
			case strSliceFlag:
				if set(strSliceFlag) == false {
					for _, key := range name {
						if str, ok := holder(key); ok && len(str.([]string)) != 0 {
							// 判断f 是何类型  基本类型的话 直接取str的第一个元素，结束
							// 如果是 slice的话，MakeSlices
							if f.Kind() == reflect.Slice {
								sv := reflect.MakeSlice(f.Type(), len(str.([]string)), len(str.([]string)))
								kind := sv.Index(0).Kind()
								if sv.Index(0).Type() == timeType {
									kind = timeTypeKind
								}
								if converter, ok := converters[kind]; ok {
									for idx, vs := range str.([]string) {
										sv.Index(idx).Set(converter(vs))
									}
								}
								f.Set(sv)
							}
							break
						}
					}
				}
			}
		}
	}
}

func filedName(f reflect.StructField, tagName string) []string {
	tag := f.Tag.Get(tagName)
	if tag != "" {
		return []string{tag}
	}
	return []string{lowerCase(f.Name), f.Name, underScore(f.Name)}
}

// lowercase
func lowerCase(s string) string {
	bs := []rune(s)
	if 'A' <= bs[0] && bs[0] <= 'z' {
		bs[0] += 'a' - 'A'
	}
	return string(bs)
}

// underScore
func underScore(s string) string {
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
