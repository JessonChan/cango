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
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// decode a map[string]string to a struct.
// The first parameter must be a reflect.Ptr to a struct.
// The second parameter is a map
// The third parameter is optional,used to generate the holder's key based on the struct's Field
func decode(holder map[string]string, rv reflect.Value, filedName ...func(reflect.StructField) ([]string, entityType)) {
	doDecode(rv, func(s string, et entityType) *entityValue {
		return &entityValue{
			enc:   stringFlag,
			key:   s,
			value: holder[s],
		}
	}, append(filedName, noTagName)[0])
}

// decodeForm decodes a map[string][]string to a struct.
// The first parameter must be a reflect.Ptr to a struct.
// The second parameter is a map, typically url.Values from an HTTP request.
// The third parameter is optional,used to generate the holder's key based on the struct's Field
func decodeForm(holder map[string][]string, rv reflect.Value, filedName ...func(reflect.StructField) ([]string, entityType)) {
	doDecode(rv, func(s string, et entityType) *entityValue {
		return &entityValue{
			enc:   strSliceFlag,
			key:   s,
			value: holder[s],
		}
	}, append(filedName, noTagName)[0])
}

// decodeForm decodes holder func to a struct.
// The first parameter must be a reflect.Ptr to a struct.
// The second parameter is a func,which in args is string-key and out-args is string/[]string/gob bytes
// The third parameter is optional,used to generate the holder's key based on the struct's Field
func doDecode(rv reflect.Value, holder func(string, entityType) *entityValue, filedNameFn func(field reflect.StructField) ([]string, entityType)) {
	if rv.IsValid() == false || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = reflect.Indirect(rv)
	setValue(holder, rv, filedNameFn)
}

type encodeType int

const (
	stringFlag   encodeType = 0
	strSliceFlag            = 1
	gobBytes                = 2
)

type entityType int

const (
	defaultEntity entityType = iota
	cookieEntity
	sessionEntity
	headerEntity
)

type entityValue struct {
	enc   encodeType
	typ   entityType
	key   string
	value any
}

// setValue sets key-value to a struct
func setValue(holder func(string, entityType) *entityValue, rv reflect.Value, filedName func(field reflect.StructField) ([]string, entityType)) {
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
		if f.CanSet() {
			names, et := filedName(rv.Type().Field(i))
			if func() bool {
				kind := f.Kind()
				if kind == reflect.Slice {
					return false
				}
				if f.Type() == timeType {
					kind = timeTypeKind
				}
				// 返回值表示是否找到对应的caster
				if caster, ok := casterMap[kind]; ok {
					for _, name := range names {
						if v := holder(name, et); v != nil {
							switch v.enc {
							case stringFlag:
								f.Set(caster(v.value.(string)))
							case strSliceFlag:
								f.Set(caster(v.value.([]string)[0]))
							case gobBytes:
								_ = gob.NewDecoder(bytes.NewReader(v.value.([]byte))).DecodeValue(f)
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
				if v := holder(key, et); v != nil {
					values := v.value.([]string)
					if len(values) == 0 {
						continue
					}
					slice := reflect.MakeSlice(f.Type(), len(values), len(values))
					kind := slice.Index(0).Kind()
					if slice.Index(0).Type() == timeType {
						kind = timeTypeKind
					}
					if converter, ok := casterMap[kind]; ok {
						for idx, vs := range values {
							slice.Index(idx).Set(converter(vs))
						}
					}
					f.Set(slice)
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

func fieldTagHolder(field reflect.StructField, name string) []string {
	if strings.Contains(string(field.Tag), name) {
		tagValue := field.Tag.Get(name)
		if tagValue != "" {
			if tagValue == "~" {
				return []string{lowerCase(field.Name), field.Name, underScore(field.Name)}
			} else {
				return []string{tagValue}
			}
		}
	}
	return nil
}

var tagNames = [...]string{cookieTagName, sessionTagName, headerTagName}
var entityTypes = [...]entityType{cookieEntity, sessionEntity, headerEntity}

func fieldTagNames(field reflect.StructField) ([]string, entityType) {
	if field.Tag != "" {
		for idx, tag := range tagNames {
			names := fieldTagHolder(field, tag)
			if len(names) > 0 {
				return names, entityTypes[idx]
			}
		}
	}
	return []string{lowerCase(field.Name), field.Name, underScore(field.Name)}, defaultEntity
}

func noTagName(f reflect.StructField) ([]string, entityType) {
	return filedName(f, ""), defaultEntity
}

// lowerCase = AbcDef -> abcDef
func lowerCase(s string) string {
	if s == "" {
		return s
	}
	bs := []rune(s)
	if 'A' <= bs[0] && bs[0] <= 'Z' {
		bs[0] += 'a' - 'A'
	}
	return string(bs)
}

// underScore = AbcDef -> abc_def
func underScore(s string) string {
	if s == "" {
		return s
	}
	bs := make([]rune, 0, 2*len(s))
	for _, s := range s {
		if 'A' <= s && s <= 'Z' {
			s += 'a' - 'A'
			bs = append(bs, '_')
		}
		bs = append(bs, s)
	}
	if bs[0] == '_' {
		bs = bs[1:]
	}
	return string(bs)
}
