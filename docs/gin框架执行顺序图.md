# 1. Gin框架执行顺序
## 1.1 Gin 中间件本质
> Gin 中间件 = 一个“可前可后”的函数栈(AOP)

```go
func (c *gin.Context) {
	// before
	c.Next()
	// after
}
```

## 1.2 完整执行顺序
示例代码
```go
r := gin.New()

r.Use(M1, M2)

r.GET("/ping", M3, H)

func M1(c *gin.Context) {
fmt.Println("M1 before")
c.Next()
fmt.Println("M1 after")
}

func M2(c *gin.Context) {
fmt.Println("M2 before")
c.Next()
fmt.Println("M2 after")
}

func M3(c *gin.Context) {
fmt.Println("M3 before")
c.Next()
fmt.Println("M3 after")
}

func H(c *gin.Context) {
fmt.Println("Handler")
}
```
执行顺序:
```text
M1 before
M2 before
M3 before
Handler
M3 after
M2 after
M1 after
```

结论：
- 先注册 -> 先执行 `before`
- 后注册 -> 先执行 `after`
- 典型的栈结构(LIFO)

## 1.3 Gin内部是如何跑起来的
1️⃣ Handlers 链表
> gin 内部会把所有的 middleware + handler 合并为一个 slice

```go
c.handlers = []HandlerFunc{
	M1, M2, M3, H,
}
```
并维护一个索引:`c.index = -1`

2️⃣ Next() 的真正含义

简化源码逻辑:
```go
func (c *Context) Next() {
    c.index++
    for c.index < len(c.handlers) {
        c.handlers[c.index](c)
        c.index++
    }
}
```
- `Next()` 不等于执行下一个
- 是： 执行完剩下的所有handler
- `after` 逻辑是函数返回时自然执行的


## 1.4 `Abort` 的真正行为
示例:
```go
func Auth(c *gin.Context) {
    if !ok {
        c.AbortWithStatus(401)
        return
    }
    c.Next()
}
```
Abort 之后会发生什么呢？
- ❌ 后续 middleware / handler 不会执行
- ✅ 已经进入的 middleware 的`after`仍然会执行


Abort的本质是:
```go
c.index = abortIndex
```


## 1.5 几种特殊中间件行为
| 行为             | 是否执行后续 | 是否执行 after  |
| -------------- | ------ | ----------- |
| 不调用 Next       | ❌      | ❌           |
| Next 后 Abort   | ❌      | ✅           |
| Abort 后 return | ❌      | ✅           |
| panic          | ❌      | 交给 Recovery |


## 1.6 全局/路由/分组 中间件顺序
示例
```go
r.Use(G1)

v1 := r.Group("/v1")
v1.Use(G2)

v1.GET("/ping", G3, H)
```
实际顺序:
```text
G1 before
G2 before
G3 before
H
G3 after
G2 after
G1 after
```
规则:
> 越靠近 Handler 的中间件，越晚注册，越先返回

## 1.7 为什么 日志/Trace/Recovery 顺序不能乱
推荐，黄金顺序
```go
r.Use(
	RequestId(),
	Trace(),
	AccessLog(),
	Recovery()
	)
```
| 中间件       | 必须在前          |
| --------- | ------------- |
| RequestID | 后面都要用         |
| Trace     | log / error 要 |
| AccessLog | 要包住 handler   |
| Recovery  | 要兜底 panic     |





# Gin Recovery
> Recovery 是 Gin 中最后一道防线，兜底 panic,保证服务不崩

作用:
- 捕获handler / middleware 的panic
- 记录堆栈
- 返回500
- 不中断进程

## 官方 Recovery 使用
```go
r.Use(gin.Recovery())
```
等价于
```go
r.Use(gin.RecoveryWithWriter(os.Stderr))
```


## Recovery 核心源码(精简版)
```go
func RecoveryWithWriter(out io.Writer) HandlerFunc {
	logger := log.New(out, "\n\n\x1b[31m", log.LstdFlags)

	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// 1️⃣ 判断是否是断开连接类错误
				if brokenPipe := isBrokenPipe(err); brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// 2️⃣ 打印 panic 日志 + stack
				stack := stack(3)
				logger.Printf("[Recovery] panic recovered:\n%s\n%s\n",
					err, stack)

				// 3️⃣ 返回 500
				c.AbortWithStatus(500)
			}
		}()

		// 4️⃣ 继续执行后续 handler
		c.Next()
	}
}
```