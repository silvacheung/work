# work
一个golang的任务协程池简单实现，用于限制系统无休止开辟协程执行任务的场景
# 简单使用
```
// 开启一个任务池,指定最大的协程数为1000,worker数为cpu的核数
work.Start(WithName("default"), WithGoroutineNum(1000), WithWorkerNum(runtime.NumCPU()))

// 初始化一个任务
task = work.New("empty-task", func() {
    // empty task
})

// 添加任务到池
if err := work.Put("default", task); err != nil {
    // handle error``
}

// 根据需要关闭任务池
work.Stop("default")
```

# 简单测试
[简单测试](./pool_test.go)
