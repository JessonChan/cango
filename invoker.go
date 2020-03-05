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
)

const (
	invokeBySelf     = 0
	invokeByReceiver = 1
)

// invoker实际执行请求的函数
// kind用来表示是通过struct来注册的还是只是通过函数来注册的
// 0 invokeBySelf --- 通过函数
// 1 invokeByReceiver --- 通过struct
type invoker struct {
	kind int
	*reflect.Method
}
