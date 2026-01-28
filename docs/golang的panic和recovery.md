# panic
当`panic`后会发生什么
```go
panic("boom")
```
运行时行为
```text
当前 goroutine
  ↓
开始回溯调用栈
  ↓
执行 defer（倒序）
  ↓
如果没人 recover
  ↓
打印 stack trace
  ↓
当前 goroutine 终止
```
- `panic` 只影响当前 goroutine
- main goroutine panic -> 进程退出
- 子 goroutine panic -> 仅该 goroutine死亡(往往是bug)


panic 本质不是“异常”，而是控制流跳转
- 只字节中断正常返回路径
- 不再执行后续语句
- 只逆序执行 defer


# recover 的使用规则
1. recover 只在 defer 中生效
2. recover 只对正在panic 的goroutine 生效


# panic 常见来源
- 空指针, `var p *int; *p`
- slice越界, `a[10`
- map并发写, `fatal error`
- 类型断言,`v.(int)`
- 手动panic, `panic(err)`


# panic vs error (设计哲学)
Go 官方共识：
> error 是业务可预期错误
> 
> panic 是程序不可继续的bug

| 场景            | 用什么   |
| ------------- | ----- |
| 参数校验          | error |
| DB 返回空        | error |
| 网络超时          | error |
| nil pointer   | panic |
| invariant 被破坏 | panic |


# 一个标准的 panic - recover 莫办
```go
func SafeRun(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    fn()
    return nil
}
```
常见用途:
- goroutine 池
- 任务调度器
- worker 模型