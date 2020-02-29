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
	"testing"
	"time"
)

type Person struct {
	Name     string
	Age      int
	Height   float32
	Birthday time.Time `cookie:"birth"`
	School
}

type Persons struct {
	Name     []string
	Age      []int
	Height   float32
	Birthday time.Time `cookie:"birth"`
	School
}

type School struct {
	IsGood bool
}

func Test_decode(t *testing.T) {
	type args struct {
		holder    map[string]string
		v         interface{}
		filedName func(field reflect.StructField) []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "person",
			args: args{
				holder: map[string]string{
					"Name":   "Cango",
					"Age":    "1",
					"Height": "1.5",
					"birth":  time.Now().Format(longSimpleTimeFormat),
					"IsGood": "true",
				},
				v: &Person{},
				filedName: func(field reflect.StructField) []string {
					return filedName(field, "cookie")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decode(tt.args.holder, reflect.ValueOf(tt.args.v), tt.args.filedName)
			t.Log(tt.args.v)
		})
	}
}

func Test_decodeForm(t *testing.T) {
	type args struct {
		holder    map[string][]string
		v         interface{}
		filedName func(field reflect.StructField) []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "persons",
			args: args{
				holder: map[string][]string{
					"Name":   []string{"Cango", "Golang"},
					"Age":    []string{"1", "2"},
					"Height": []string{"1.5"},
					"birth":  []string{time.Now().Format(longSimpleTimeFormat)},
					"IsGood": []string{"true"},
				},
				v: &Persons{},
				filedName: func(field reflect.StructField) []string {
					return filedName(field, "cookie")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeForm(tt.args.holder, reflect.ValueOf(tt.args.v), tt.args.filedName)
			t.Log(tt.args.v)
		})
	}
}
