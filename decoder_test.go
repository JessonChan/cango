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

func Test_assign(t *testing.T) {
	now := time.Now()
	holder := map[string]string{
		"Name":   "Cango",
		"Age":    "1",
		"Height": "1.5",
		"birth":  now.Format(longSimpleTimeFormat),
		"IsGood": "true",
	}
	filedName := func(field reflect.StructField) []string {
		return filedName(field, "cookie")
	}
	var p Person
	var v = newValue(reflect.TypeOf(p))
	decode(holder, addr(v), filedName)
	p = v.Interface().(Person)
	if p.Name != "Cango" {
		t.Fail()
	}
	if p.Age != 1 {
		t.Fail()
	}
	if p.Height != 1.5 {
		t.Fail()
	}
	if p.IsGood != true {
		t.Fail()
	}
	if p.Birthday.Format(longSimpleTimeFormat) != now.Format(longSimpleTimeFormat) {
		t.Fail()
	}
}

func Test_decode(t *testing.T) {
	type args struct {
		holder    map[string]string
		v         interface{}
		filedName func(field reflect.StructField) []string
	}
	now := time.Now()
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
					"birth":  now.Format(longSimpleTimeFormat),
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
			switch tt.args.v.(type) {
			case *Person:
				v := tt.args.v.(*Person)
				if v.Name != "Cango" {
					t.Fail()
				}
				if v.Age != 1 {
					t.Fail()
				}
				if v.Height != 1.5 {
					t.Fail()
				}
				if v.IsGood != true {
					t.Fail()
				}
				if v.Birthday.Format(longSimpleTimeFormat) != now.Format(longSimpleTimeFormat) {
					t.Fail()
				}

			default:
				t.Fail()
			}
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
					"Name":   {"Cango", "Golang"},
					"Age":    {"1", "2"},
					"Height": {"1.5"},
					"birth":  {time.Now().Format(longSimpleTimeFormat)},
					"IsGood": {"true"},
				},
				v: &Persons{},
				filedName: func(field reflect.StructField) []string {
					return filedName(field, "cookie")
				},
			},
		},
	}
	now := time.Now()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeForm(tt.args.holder, reflect.ValueOf(tt.args.v), tt.args.filedName)
			switch tt.args.v.(type) {
			case *Persons:
				v := tt.args.v.(*Persons)
				if v.Name[0] != "Cango" {
					t.Fail()
				}
				if v.Age[0] != 1 {
					t.Fail()
				}
				if v.Height != 1.5 {
					t.Fail()
				}
				if v.IsGood != true {
					t.Fail()
				}
				if v.Birthday.Format(longSimpleTimeFormat) != now.Format(longSimpleTimeFormat) {
					t.Fail()
				}
			default:
				t.Fail()

			}
		})
	}
}

func Test_lowerCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{""},
			want: "",
		},
		{
			args: args{"A"},
			want: "a",
		},
		{
			args: args{"Ab"},
			want: "ab",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowerCase(tt.args.s); got != tt.want {
				t.Errorf("lowerCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_underScore(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{""},
			want: "",
		},
		{
			args: args{"A"},
			want: "a",
		},
		{
			args: args{"ABC"},
			want: "a_b_c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := underScore(tt.args.s); got != tt.want {
				t.Errorf("underScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
