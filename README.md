
#   简介

**cango** **web 开发框架** 区别于其它框架的重要的功能就是基于tag的URI自发现机制，不需要显示的定义路由，让开发人员专注于代码而不需要过分关注路由定义。
## 安装

cango需要Go1.12及以上。  
安装cango

```bash
go get -u github.com/JessonChan/cango
```

## Hello, World!

开箱，看看cango最简单的样子。  
创建一个文件，如app.go，写入下面的代码。
```go
package main

import "github.com/JessonChan/cango"

func main() {
	cango.
		NewCan().
		RouteFunc(func(ps struct {
			cango.URI `value:"/;/hello"`
		}) interface{} {
			return cango.Content{String: "Hello,World!"}
		}).
		Run(cango.Addr{Port: 8080})
}
```

```bash
go run app.go
```
打开 `http://127.0.0.1:8080`,就会看到`Hello,World!` ，也可以打开 `http://127.0.0.1:8080/hello`，会看到同样的内容 
上面的例子只是为了展示初步的使用，还可以有另外的写法，就是将函数定义在特定的struct上，代码如下
```go
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
```

```bash
go run app.go
```
打开 `http://127.0.0.1:8080/hello`,就会看到`Hello,World!`;打开 `http://127.0.0.1:8080/hello`,就会看到`Hello,Cango!`；
上面的例子只是为了展示初步的使用，还可以有另外的写法，就是将函数定义在特定的struct上，代码如下

## 路由简介 

cango中，路由都定义在方法参数或者结构体上，以上面的例子，只需要定义将cango.URI做为成员引入实现URI接口的函数，就会自动写入路由，路由模式是由tag中的value值来定义的，需要特殊说明的是value值可以使用`;`来定义多个。  
定义的函数入参为cango.URI类型，出参为interface{}，当前版本必须要返回值的，返回值做为请求的返回依据。

```go
// 路由一个方法
can.RouteFunc(...func(cango.URI)interface{})
// 定义方法并且带有路由前缀（便于版本、分组等管理）
can.RouteFuncWithPrefix(prefix, ...func(cango.URI)interface{})
// 路由结构体上所有的方法
can.Route(cango.URI)
// 路由结构体上所有的方法，并使用前缀
can.RouteWithPrefix(cango.URI)
```
* `can` 是cango.NewCan()的实例
* `cango.URI`是用来保存路由及请求相关数据的。
* `prefix`是定义在路由上的前缀，使用prefix参数后，路由地址为prefix+URI-Tag-Value，看以下的例子
* `interface{}` 映射请求的返回值，当前版本支持的类型如下

```go
cango.ModelView //返回模板和数据
cango.StaticFile //文件类型，会调用http.ServeFile,一般用不到
cango.Redirect //重定向
cango.Content //上面的例子中使用，会直接返回Content.String
```
如果返回值不是以上类型，当前版本的处理逻辑是返回JSON。   

下面的代码为我们展示本小节的所有内容。
创建app.go和view和static两个文件夹，在view上创建index.html模板，在static创建index.css文件。形式如下：
```bash
.
├── app.go
├── static
│   └── index.css
└── view
    └── index.html
```
app.go中写入
```go
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
		Run(cango.Addr{Port: 8080})
}
```
在index.html中写入
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/index.css">
</head>
<body>
<span class="blue">{{.Content}}</span>
</body>
</html>
```
在css中写入
```css
.blue{
    color: blue;
}
```
在终端执行以下命令
```bash
go run app.go
```
在浏览器打开 `http://127.0.0.1:8080/index.html`、`http://127.0.0.1:8080/v2/index.html`和`http://127.0.0.1:8080/goto`

### 请求方法

在以上两个例子中都只有GET方法，如果相要定义其它访方法只要在请求中加入对应的方法即可。

```go
type GetMethod httpMethod
type HeadMethod httpMethod
type PostMethod httpMethod
type PutMethod httpMethod
type PatchMethod httpMethod
type DeleteMethod httpMethod
type OptionsMethod httpMethod
type TraceMethod httpMethod
```
还以上面的例子说明。在app.go中加入新的方法，更新后如下
```go
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
// 演示POST方法
func Ping(ps struct {
	cango.URI `value:"/ping.json"`
	cango.PostMethod
}) interface{} {
	return map[string]string{
		"Pong": "ok",
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

```
重新执行 `go run app.go`，在终端执行
```bash
curl  -d ""  http://127.0.0.1:8080/ping.json
```
可以看到
```bash
{"Pong":"ok"}
```

### 路由变量
当前版本支持的路由变量定义方式是使用`{}`，将上面例子中Ping函数进行修改如下
```go
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
		"Age":  ps.CarAge,
		"Color": ps.Color,
	}
}
```
重新执行 `go run app.go`，在终端执行
```bash
curl  -d ""  http://127.0.0.1:8080/ping/2020/15/white.json
```
可以看到
```bash
{"Age":1,"Color":"white","Pong":"ok","Year":2020}
```
也就是定义在路由中的变量

```go
// GET http://localhost:8080/hello%20world

app.Get("/:value", func(c *fiber.Ctx) {
  c.Send("Get request with value: " + c.Params("value"))
  // => Get request with value: hello world
})
```

### 通配符路由
当前版本支持在路由中最多包含一个 `*` 通配符的路径，并且通配符路由不能再包含路径变量，定义方法如下
```go
cango.URI `value:"/goto/*"`
```
