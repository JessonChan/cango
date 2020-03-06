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

type FormValue interface {
	constructor
	FormPlaceHolder()
}

var formValueType = reflect.TypeOf((*FormValue)(nil)).Elem()
var formValueTypeName = formValueType.Name()
var valueOfEmptyForm = reflect.ValueOf(&emptyFormValueConstructor{})

type emptyFormValueConstructor struct {
}

func (e *emptyFormValueConstructor) Construct(r *http.Request) {
}
func (e *emptyFormValueConstructor) FormPlaceHolder() {
}
