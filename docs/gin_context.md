# gin.Context详解
`gin.Context` = HTTP 请求的一次生命周期上下文 + 共享数据容器

封装了:
- `http.Request`
- `http.ResponseWriter`
- 路由参数
- 中间件控制
- key-value存储(夸middleware)
- 请求/响应操作工具

# API

## 最重要的四个维度
可以把`gin.Context`方法分为四类:
- 请求读取：Query / Param / Body / Header
- 响应写入: JSON /String / Status
- 中间件控制: Next / Abort
- 上下文数据: Set / Get

## 请求相关方法
1️⃣ 路由参数
```go
GET /users/:id
```
```go
id := c.Param("id")
```
- 来自路由匹配结果
- 只和 path 有关
- 比 Query 优先级高

2️⃣ Query参数(URL ?)
```go
name := c.Query("name")
name := c.DefaultQuery("name","default")
name, ok := c.GetQuery("name")
```
- Query: 不存在返回 `""`
- GetQuery: 返回`(value, exists)`
- DefaultQuery: 支持默认值

3️⃣ Header
```go
ua := c.GetHeader("User-Agent")

// 等价

c.Request.Header.Get("User-Agent")
```
- 推荐使用`gin`框架封装的

4️⃣ Body 绑定(重要)

JSON绑定
```go
var req CreateUserReq
if err:= c.ShouldBindJSON(&req);err != nil {
	c.JSON(400,gin.H{
		"error": err.Error()
   })
	return
}
```
- BindJSON: 失败自动写400
- ShouldBindJSON: 手动处理error (企业级推荐)
- ShouldBind: 自动根据Content-Type


5️⃣原始Body(谨慎使用)
```go
body , _ := io.ReadAll(c.Request.Body)
```
- ⚠️：Body只能读取一次，日志/签名校验时要注意(通常使用中间件缓存)


6️⃣ 响应相关方法


JSON响应
```go
c.JSON(200,gin.H{
	"code": 0,
	"msg": "oK"
})
```
- 本质：Content-Type: application/json


String/HTML/Data
```go
c.String(200,"hello")
c.HTML(200,"index.tmpl",data)
c.Data(200,"application/octet-stream",bytes)
```

设置状态码
```go
c.Status(204)
```
- 不会写body


重定向
```go
c.Redirect(302,"/login")
```

7️⃣ 中间件控制(核心)

Next - 执行后续中间件
```go
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before
		c.Next()
		// after
    }
}
```
- 这是 gin 能做 AOP的核心


Abort - 中断执行链
```go
c.Abort()
```
- 后续中间件 & handler 不会执行
- 已执行的不会回滚


AbortWithStatus / JSON
```go
c.AbortWithStatus(401)
```
```go
c.AbortWithStatusJSON(403,gin.H{"msg":"forbidden"})
```
- 鉴权/限流 必用


判断是否中断
```go
if c.IsAborted() {
	return
}
```

8️⃣ 上下文读取(夸中间件共享)

Set / Get
```go
c.Set("user_id", 123)

uid,ok := c.Get("user_id")
```
- 常用于：request_id, user信息, trace_id, auth结果等

MustGet
```go
uid := c.MustGet("user_id").(int)
```
- 不存在直接panic


9️⃣ 错误管理

1
```go
c.Error(err)
```
- 把错误挂到context
- 不会自动返回
- 中间件可统一处理
```go
c.Errors.Last()
```


2
```go
for _,e := range c.Errors {
	log.Error(e.Err)
}
```
- 非常适合做统一的 error middleware



🔟 Request/Response 底层对象

原始Request
```go
req := c.Request
```
- `*http.Request`
- header / Body / Context 都在这里


Writer(响应控制)
```go
c.Writer.Status()
c.Writer.Size()
```
- AccessLog / metrics 必用


# 生命周期总结
```text
请求进入
  ↓
创建 gin.Context
  ↓
middleware before
  ↓
handler
  ↓
middleware after
  ↓
写响应
```