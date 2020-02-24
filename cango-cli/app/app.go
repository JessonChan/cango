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

import (
	"github.com/JessonChan/cango"
)

func Index(ps struct {
	cango.URI `value:"/index.html"`
}) interface{} {
	return cango.ModelView{
		Tpl: "index.html",
		Model: map[string]string{
			"Title":   "cango",
			"Content": "hello,cango",
		},
	}
}
func Ping(ps struct {
	cango.URI `value:"/ping/{year}/{car_age}/{Color}.json"`
	cango.PostMethod
	Year   int
	CarAge int
	Color  string
}) interface{} {
	return map[string]interface{}{
		"Pong":  "ok",
		"Year":  ps.Year,
		"Age":   ps.CarAge,
		"Color": ps.Color,
	}
}

func main() {
	cango.
		NewCan().
		RouteFunc(func(ps struct {
			cango.URI `value:"/;/hello"`
		}) interface{} {
			return cango.Content{String: "Hello,World!"}
		}).
		// 注册两个handle，一个是有函数名，一个没有
		RouteFunc(Index, func(ps struct {
			cango.URI `value:"/goto"`
		}) interface{} {
			return cango.Redirect{Url: "/index.html"}
		}).
		RouteFuncWithPrefix("/v2", Index).
		RouteFunc(Ping).
		Run(cango.Addr{Port: 8080})
}
