# 简介

**TODO**  
a) 动态tag

**cango** **web 开发框架** 通过tag进行URI自发现，不需要显式的定义路由；同时支持对于请求参数自动构造，减少编码和转化，更高效简洁的进行开发。

## 安装

cango需要Go1.18及以上，同时本教程需要依赖go mod。   
添加cango包

```bash
go get -u github.com/JessonChan/cango
```

为了更好更快的了解cango，建议通过以下命令安装`cango-cli`工具

```bash
GO111MODULE=off go get -u  github.com/JessonChan/cango-cli
GO111MODULE=off go install github.com/JessonChan/cango-cli
```
## 编程代码片段
在编辑器中添加代码片段，可以方便、高效的生成代码。
### vscode
添加的方法是在设置中选择配置代码片段
```json
	"Gen Cango Func": {
		"prefix": ".cf",
		"body": [
			"func (p *$1) $2(ps struct {",
			"\tcango.URI `value:\"/$3\"`",
			"}) interface{} {",
			"\treturn nil",
			"}"
		],
		"description": "Log output to console"
	},
	"Gen Cango Post": {
		"prefix": ".post",
		"body": [
			"func (p *$1) $2(ps struct {",
			"\tcango.URI `value:\"/$3\"`",
			"\tcango.PostMethod",
			"}) interface{} {",
			"\treturn nil",
			"}"
		],
		"description": "cango Post method"
	},
	"Gen Cango Struct": {
		"prefix": ".ctrl",
		"body": [
			"type $1 struct{",
			"\tcango.URI `value:\"$2\"`",
			"}\n",
			"var _ = cango.RegisterURI(&$1{})",
		],
		"description": "生成cango的struct"
	},
	"Gen Cango Filter": {
		"prefix": ".filter",
		"body": [
			"type $1 struct{",
			"\tcango.Filter `value:\"$2\"`",
			"}\n",
			"var _ = cango.RegisterFilter(&$1{})\n",
			"func (p *$1) PreHandle(request *cango.WebRequest) interface{} {",
			"\treturn true",
			"}",
			"func (p *$1) PostHandle(request *cango.WebRequest) interface{} {",
			"\treturn true",
			"}",
		],
		"description": "生成cango的struct"
	},
```

### Goland
Goland的中添加 `Live Template` 
```go
func (p *$NAME$) $FUNC$(ps struct {
	cango.URI `value:"/$URL$"`
	cango.PostMethod
}) interface{} {
	return nil
}

```

## Hello, World!

开箱，看看cango最简单的样子。在终端进入自己熟悉的目录（比如/tmp）执行以下命令：

```bash
cango-cli start
cd start
```  

本命令完成的工作：创建一个文件，route.go，写入下面的代码（也可以手动输入以下代码）。

```go start.route.go
package main

import "github.com/JessonChan/cango"

func main() {
	cango.
		NewCan().
		RouteFunc(func(cango.URI) interface{} {
			return cango.Content{String: "Hello,World!"}
		}).
		Run()
}
```

```bash
go run route_func.go
```

打开 `http://127.0.0.1:8080`[链接](http://127.0.0.1:8080),就会看到`Hello,World!` 。 接下来，通过修改过上面的代码来展示如何通过tag来定义路由。

```go start.func.go
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
		Run()
}
```     

这段代码在`start`文件夹的`route_func.go`中，通过对func入参`ps`来定义两个等效的路由`/`和`/hello`，
意味着，我们既可以通过`http://127.0.0.1:8080`[链接](http://127.0.0.1:8080)来访问，也可以通过
`http://127.0.0.1:8080/hello`[链接](http://127.0.0.1:8080/hello)来访问。

上述两段的代码初步展示cango的使用。在项目中，更为推荐的写法是将函数定义在特定的struct上，创建新的文件 route_struct.go，并写入代码如下(start目录下的route_struct.go)

```go start.route_struct.go
package main

import "github.com/JessonChan/cango"

type Ctrl struct {
	cango.URI `value:"/hello"`
}

// 路由定义是方法的 receiver 中定义的
// 这个方法对应的 /hello
func (c *Ctrl) Hello(cango.URI) interface{} {
	return cango.Content{String: "Hello,World!"}
}

// 路由定义为 receiver中的URI tag和 参数列表中的 URI和tag的相加，
// 这个示例为 /hello/world.html
func (c *Ctrl) World(ps struct {
	cango.URI `value:"/world.html"`
}) interface{} {
	return cango.Content{String: "Hello,Cango!"}
}

func main() {
	cango.
		NewCan().
		Route(&Ctrl{}).
		Run()
}

```

```bash
go run route_struct.go
```

打开 `http://127.0.0.1:8080/hello`[链接](http://127.0.0.1:8080/hello),就会看到`Hello,World!`;打开 `http://127.0.0.1:8080/hello/world.html`[链接](http://127.0.0.1:8080/hello/world.html),就会看到`Hello,Cango!`；

## 路由简介

cango中，路由都定义在控制器结构体和控制器方法参数上。在`Hello, World!`这一节中已经对此进行初步的说明，可以知道，只需要定义将cango.URI做为成员变量引入结构体（控制器或者控制器方法入参），
就会自动写入路由，路由规则是由tag中的value值来定义的。 需要特殊说明的是value值可以使用`;`来定义多个同时生效的平行规则，也可以为空。
定义的函数入参为cango.URI类型，出参为interface{}，当前版本必须要返回值的， 返回值做为请求的返回依据。

```go
// RouteFunc 方法路由，可以传入多个方法
can.RouteFunc(...func (cango.URI,...cango.Constructor)interface{})
// RouteFuncWithPrefix 带有前缀的方法路由，可以传入多个方法（便于版本、分组等管理）
can.RouteFuncWithPrefix(prefix, ...func(cango.URI, ...cango.Constructor)interface{})
// 路由结构体上所有的方法
can.Route(cango.URI)
// 路由结构体上所有的方法，并使用前缀
can.RouteWithPrefix(cango.URI)
// 在定义struct的时候引入，也这是非常推荐的方法
var _ = cango.RegisterURI(cango.URI)
// RegisterURIWithPrefix 在定义struct的时候引入，同时使用prefix做为路由前缀，也这是非常推荐的方法
var _ = cango.RegisterURIWithPrefix(prefix string, uri URI) 
```

* `can` 是cango.NewCan()的实例
* `cango.URI`是用来保存路由及请求相关数据的。
* `prefix`是定义在路由上的前缀，使用prefix参数后，路由地址为prefix+URI-Tag-Value
* `interface{}` 映射请求的返回值，当前版本支持的类型如下

```go
// 返回模板和数据
cango.ModelView
// 文件类型，会调用http.ServeFile,一般用不到
cango.StaticFile
// 重定向
cango.Redirect
cango.RedirectWithCode
// 上面的例子中使用，会直接返回Content.String
cango.Content
cango.ContentWithCode
```

如果返回值不是以上类型，当前版本的处理逻辑是返回JSON。

下面的代码为我们展示本小节的所有内容。 在终端执行

```bash
cango-cli app
```

本命令完成以下工作：创建app.go和view和static两个文件夹，在view上创建index.html模板，在static创建index.css文件。形式如下：

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
		Run()
}
```

同时在index.html中写入

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
.blue {
    color: blue;
}
```

在终端执行以下命令

```bash
go run app.go
```

在浏览器打开 `http://127.0.0.1:8080/index.html`[链接](http://127.0.0.1:8080/index.html)、`http://127.0.0.1:8080/v2/index.html`[链接](http://127.0.0.1:8080/v2/index.html)和`http://127.0.0.1:8080/goto`[链接](http://127.0.0.1:8080/goto)

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
		Run()
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

### 请求参数

获取请求中的变量和路由变量的定义方式相同，只需要在入参加入相应的变量名称。


### Header变量

获取header中的值，只需要结构体成员使用header这个tag名就可以，如

```go 
struct{
    Name string `header:"~"`
}
```

需要指出的是，在header后使用`~`这个符号表示，对应的header名称是自动从成员变量名读取。


### Cookie变量

获取cookie中的值，只需要结构体成员使用cookie这个tag名就可以，如

```go 
struct{
    Name string `cookie:"~"`
}
```

需要指出的是，在cookie后使用`~`这个符号表示，对应的cookie名称是自动从成员变量名读取。

### Session变量

需要在配置文件中开启session，方法为在cango.ini中添加`cookie_session_key=my_session_key`，加入后会 自动开启session，当前使用的session是gorilla/session包。
获取session中的值，只需要结构体成员使用session这个tag名就可以，如

```go 
struct{
    Name string `session:"~"`
}
```

后使用`~`这个符号表示，对应的session名称是自动从成员变量名读取

### 通配符路由

当前版本支持在路由中最多包含一个 `*` 通配符的路径，并且通配符路由不能再包含路径变量，定义方法如下

```go
cango.URI `value:"/goto/*"`
```

### 过滤器

当前版本支持前置过滤器和后置过滤器，目前不支持对执行结果的修改。定义如下

```go
type VisitFilter struct {
cango.Filter `value:"/static/*.css;/static/*.js"`
}
func (v *VisitFilter) PreHandle(req *cango.WebRequest) interface{} {
return true
}
func (v *VisitFilter) PostHandle(req *cango.WebRequest) interface{} {
return true
}
```

如上所写，使用tag定义filter路径是最推荐的。  
上面的代码完成会在执行controller的函数前先执行PreHandle，再执行函数，最后执行PostHandle。它的作用范围是所有在`static`目录下的的`css`和`js`静态文件。  
但是你也可以根据自己的需要，只注册某些cango.URI，如下所示。

```go
can.Filter(f cango.Filter, uris ...cango.URI)
cango.RegisterFilter(cango.Filter)
```

上在的接口支持对某些接口单独的Filter。

### 配置文件

当前版本的cango支持配置参数从文件读取，配置文件位置必须是程序运行目录下的`conf/cango.ini`文件，支持的配置参数如下表  
| 参数名 |参数类型 |参数说明 | | ---- |---- |---- | | host/Host |string |监听的主机 | | port/Port |int |监听的端口 | | tpl_dir/tplDir/TplDir
|string |模板文件文件夹，相对程序运行路径 | | static_dir/staticDir/StaticDir |string |静态文件文件夹，相对程序运行路径 | |
tpl_suffix/tplSuffix/TplSuffix |[]string |模板文件后续名，默认为 .tpl 和 .html | | debug_tpl/debugTpl/DebugTpl |bool |是否调试页面，true
表示每次都重新加载模板 | | canlog_path/canlogPath/CanlogPath |string |日志文件位置，绝对路径 | |
cookie_session_key/cookieSessionKey/CookieSessionKey |string |gorilla cookie store 的key | |
cookie_session_secure/cookieSessionSecure/CookieSessionSecure |string |gorilla cookie store 加密使用的key |

### 和其它框架结合（Gin为例）

以下为如何和gin框架结合

```go
package main

import "github.com/gin-gonic/gin"
import "github.com/JessonChan/cango"

func main() {

	can := cango.
		NewCan().
		RouteFunc(func(ps struct {
			cango.URI `value:"/;/hello"`
		}) interface{} {
			return map[string]string{"message": "hello,world"}
		})
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	can.FallbackHandler(r)
	can.Run()
}
```
如下，引入了cango，可以很好的兼容原来的gin项目。
```go
package main

import (
	"github.com/JessonChan/cango"
	"github.com/gin-gonic/gin"
)

type Api struct {
	cango.URI `value:"/api"`
}

// 通过 http://127.0.0.1:8080/api/msg/bob.json 来访问
func (a *Api) Get(ps struct {
	cango.URI `value:"/msg/{name}.json"`
	Name      string
}) any {
	return map[string]string{"msg": "Very Good! " + ps.Name}
}

func main() {

	r := gin.Default()
	cango.NewCan().Route(&Api{}).GinRoute(r)
	r.Run()
}

```

## 更多例子

为了更好的理解和使用cango，`cango-cli`中还包含`can`、`short_url`和`demo`三个示例，请自己执行查看。
另外，可以查看[can_blog](http://www.github.com/JessonChan/can_blog)这个简单的博客项目。