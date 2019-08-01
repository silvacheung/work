# taskpool
一个golang的任务协程池简单实现，用于限制系统无休止开辟协程执行任务的场景
# 简单使用
```
// 开启一个任务池,指定最大的协程数为1000,worker数为cpu的核数
taskpool.StartPool(WithPoolName("default"), WithGoroutineNum(1000), WithWorkerNum(runtime.NumCPU()))

// 初始化一个任务
task = taskpool.New("empty-task", func() {
    // empty task
})

// 添加任务到池
if err := taskpool.Put("default", task); err != nil {
    // handle error
}

// 根据需要关闭任务池
taskpool.StopPool("default")
```

# 简单测试
[简单测试](./pool_test.go)
