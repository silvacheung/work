package task

import (
	"errors"
	"runtime"
	"strings"
)

const (
	DefaultPool = "default"
	DefaultTask = "default"
)

var (
	DefaultWorkerNum = runtime.NumCPU()
	PoolNotFound     = errors.New("task pool not found")
	PoolIsClosing    = errors.New("task pool is closing")
	PoolIsClosed     = errors.New("task pool is closed")
)

type Option func(*option)

// 队列配置选项
type option struct {
	// 任务池名称
	poolName string
	// 执行队列大小
	goroutineNum int
	// worker数目
	workerNum int
}

func WithPoolName(name string) Option {
	return func(opt *option) {
		opt.poolName = name
	}
}

func WithGoroutineNum(num int) Option {
	return func(opt *option) {
		opt.goroutineNum = num
	}
}

func WithWorkerNum(num int) Option {
	return func(opt *option) {
		opt.workerNum = num
	}
}

func initOptions(opts ...Option) (opt *option) {

	opt = &option{}

	for _, optFun := range opts {
		optFun(opt)
	}

	if strings.Trim(opt.poolName, " ") == "" {
		opt.poolName = DefaultPool
	}

	if opt.workerNum <= 0 {
		opt.workerNum = DefaultWorkerNum
	}

	if opt.goroutineNum <= 0 {
		opt.goroutineNum = opt.workerNum
	}
	return
}
