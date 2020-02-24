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
package main

import "github.com/JessonChan/cango"

type CanApp struct {
	cango.URI `value:"/hello"`
}

// 路由定义是方法的 receiver 中定义的
// 这个方法对应的 /hello
func (c *CanApp) Hello(cango.URI) interface{} {
	return cango.Content{String: "Hello,World!"}
}

// 路由定义为 receiver中的URI tag和 参数列表中的 URI和tag的相加，
// 这个示例为 /hello/world.html
func (c *CanApp) World(ps struct {
	cango.URI `value:"/world.html"`
}) interface{} {
	return cango.Content{String: "Hello,Cango!"}
}

func main() {
	cango.
		NewCan().
		Route(&CanApp{}).
		Run(cango.Addr{Port: 8080})
}
