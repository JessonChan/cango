# cango
cango 是一个实验性的web开发框架，支持restful和mvc编程，高效的使用golang的tag特性，将路由定义、路由变量赋值和更多的操作通过在tag中定义。


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

## 过滤器设计
