# cango
cango 是一个实验性的restful web开发框架，也有可能会支持模板，高效的使用golang的tag特性，将路由定义、变量赋值和更多的操作通过在tag中定义。


## 路由设计
Route方法注册所有的Controller  
Controller结构体上定义路由路径和路径上对应的变量，例如  
```go 
type Controller struct {    
    URI `value:"/blog/{blogName}/article/{articleId}"`  
    BlogName string  
    ArticleId int  
}
```
在Controller的方法传参上定义具体的执行方法，例如
```go 
func (c *Controller)Comment(param struct{
    URI `value:"/commnet/{commentId}.json"`
    CommentId int
}){
    //do sth
}
```
我们约定：  
1、所有的路由变量都使用{}包围  
2、所有的对应变量都只能使用首字母小写的驼峰来命名  
3、所有的路由方法必须只有一个参数，且实现了URI接口
4、所有的路由方法有且只有第一个返回值做为是restful的返回值，如果没有返回值则返回{}

## 路由设计具体实现漫谈  
以GET方法为例说明。  
从Controller中抽离出路由urlC和变量列表varsC，从Controller中抽离出方法路由urlM和变量列表varsM，则urlC+urlM为路由方法最终的路由，varsC+varsM为路由方法最终的变量列表。
当某个已知的路由来时，先使用gorilla/mux包Match出指定的控制器和方法，通过反射将控制器进行初始化、通过反射对路由方法的唯一参数进行初始化，然后使用反射进行调用。

定义一个map来存放路由方法，由于路由方法第一个参数就是controller本身，所以可以很好的还原现场，实现调用.  

## 过滤器设计
过滤器要支持两种URL匹配风格，一种就是在工程已经实现的，基于controller去匹配，另外一种就是通过tag来定义，如实现最通用的AntPath
