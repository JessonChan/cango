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
	"net/http"
	"reflect"
)

type Filter interface {
	PreHandle(req *http.Request) interface{}
}

func (can *Can) filter(f Filter, uri URI) {
	rp := reflect.ValueOf(uri)
	if rp.Kind() != reflect.Ptr {
		panic("route controller must be prt")
	}
	fs := can.filterMap[rp.String()]
	if len(fs) == 0 {
		fs = make([]Filter, 1)
	}
	fs = append(fs, f)
	can.filterMap[rp.String()] = fs
}

func (can *Can) Filter(f Filter, uris ...URI) *Can {
	for _, uri := range uris {
		can.filter(f, uri)
	}
	return can
}
