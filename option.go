package work

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
	DefaultWorkerNum = uint(runtime.NumCPU())
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
	queueSize uint
	// goroutine数目
	goroutineNum uint
	// worker数目
	workerNum uint
}

func WithPoolName(name string) Option {
	return func(opt *option) {
		opt.poolName = name
	}
}

func WithQueueSize(size uint) Option {
	return func(opt *option) {
		opt.queueSize = size
	}
}

func WithGoroutineNum(num uint) Option {
	return func(opt *option) {
		opt.goroutineNum = num
	}
}

func WithWorkerNum(num uint) Option {
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

	if opt.workerNum == 0 {
		opt.workerNum = DefaultWorkerNum
	}

	if opt.goroutineNum == 0 {
		opt.goroutineNum = opt.workerNum
	}

	if opt.queueSize == 0 {
		opt.queueSize = opt.goroutineNum
	}

	return
}
