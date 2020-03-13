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
	"reflect"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

func decode(holder map[string]string, rv reflect.Value, filedName ...func(reflect.StructField) []string) {
	checkSet(func(s string) (interface{}, int, bool) {
		v, ok := holder[s]
		return v, stringFlag, ok
	}, rv, append(filedName, noTagName)[0])
}
func decodeForm(holder map[string][]string, rv reflect.Value, filedName ...func(field reflect.StructField) []string) {
	checkSet(func(s string) (interface{}, int, bool) {
		v, ok := holder[s]
		return v, strSliceFlag, ok
	}, rv, append(filedName, noTagName)[0])
}

func checkSet(holder func(string) (interface{}, int, bool), rv reflect.Value, filedNameFn func(field reflect.StructField) []string) {
	if rv.IsValid() == false || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = reflect.Indirect(rv)
	setValue(holder, rv, filedNameFn)
}

const (
	stringFlag   = 0
	strSliceFlag = 1
)

func setValue(holder func(string) (interface{}, int, bool), rv reflect.Value, filedName func(field reflect.StructField) []string) {
	if rv.Kind() == reflect.Interface {
		return
	}
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if f.Kind() == reflect.Ptr {
			f = reflect.Indirect(f)
		}
		if f.Kind() == reflect.Struct && f.Type() != timeType {
			setValue(holder, f, filedName)
		}
		// 返回值表示是否找到对应的caster
		names := filedName(rv.Type().Field(i))
		if f.CanSet() {
			if func() bool {
				kind := f.Kind()
				if kind == reflect.Slice {
					return false
				}
				if f.Type() == timeType {
					kind = timeTypeKind
				}
				if caster, ok := casterMap[kind]; ok {
					for _, name := range names {
						if str, flag, ok := holder(name); ok {
							switch flag {
							case stringFlag:
								f.Set(caster(str.(string)))
							case strSliceFlag:
								f.Set(caster(str.([]string)[0]))
							}
						}
					}
					return true
				}
				return false
			}() {
				continue
			}
			if f.Kind() != reflect.Slice {
				continue
			}
			for _, key := range names {
				if str, _, ok := holder(key); ok && len(str.([]string)) != 0 {
					sv := reflect.MakeSlice(f.Type(), len(str.([]string)), len(str.([]string)))
					kind := sv.Index(0).Kind()
					if sv.Index(0).Type() == timeType {
						kind = timeTypeKind
					}
					if converter, ok := casterMap[kind]; ok {
						for idx, vs := range str.([]string) {
							sv.Index(idx).Set(converter(vs))
						}
					}
					f.Set(sv)
					break
				}
			}
		}
	}
}

func filedName(f reflect.StructField, tagName string) []string {
	if tagName != "" {
		tag := f.Tag.Get(tagName)
		if tag != "" {
			if tag != "~" {
				return []string{tag}
			}
		}
	}
	return []string{lowerCase(f.Name), f.Name, underScore(f.Name)}
}

func fieldTagNames(field reflect.StructField) []string {
	if field.Tag != "" {
		if strings.Contains(string(field.Tag), cookieTagName) {
			tagValue := field.Tag.Get(cookieTagName)
			if tagValue != "" {
				if tagValue == "~" {
					return []string{cookieHolderKey + lowerCase(field.Name), cookieHolderKey + field.Name, cookieHolderKey + underScore(field.Name)}
				} else {
					return []string{cookieHolderKey + tagValue}
				}
			}
		}
	}
	return []string{formPathHolderKey + lowerCase(field.Name), formPathHolderKey + field.Name, formPathHolderKey + underScore(field.Name)}
}

func noTagName(f reflect.StructField) []string {
	return filedName(f, "")
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
