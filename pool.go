package work

import (
	"sync"
	"sync/atomic"
)

const (
	running int32 = 0
	closing int32 = 1
	closed  int32 = 2
)

var pools sync.Map

type pool struct {
	opt          *option
	queue        *queue
	numGoroutine chan struct{}
	wait         *sync.WaitGroup
	state        int32
}

func Start(opts ...Option) {
	opt := initOptions(opts...)
	pool := &pool{
		opt:          initOptions(opts...),
		queue:        initQueue(opt.workerNum, opt.goroutineNum),
		numGoroutine: make(chan struct{}, opt.goroutineNum),
		wait:         &sync.WaitGroup{},
	}
	pools.Store(opt.poolName, pool)
	pool.start()
}

func Stop(name string) {
	if pool, err := getPool(name); err == nil {
		pool.stop()
	}
}

func Put(name string, task ...*Task) error {
	if len(name) == 0 {
		name = DefaultPool
	}
	if pool, err := getPool(name); err != nil {
		return err
	} else {
		return pool.put(task...)
	}
}

// 先将状态改为关闭中
// 然后等待所有已经在队列中的任务执行完毕后,改为关闭
func (p *pool) stop() {
	if atomic.CompareAndSwapInt32(&p.state, running, closing) {
		p.wait.Wait()
	}
	if atomic.CompareAndSwapInt32(&p.state, closing, closed) {
		p.queue.close()
		pools.Delete(p.opt.poolName)
	}
}

func (p *pool) put(task ...*Task) error {
	state := atomic.LoadInt32(&p.state)
	if state == closing {
		return PoolIsClosing
	}
	if state == closed {
		return PoolIsClosed
	}
	taskNum := len(task)
	p.wait.Add(taskNum)
	var queue chan<- *Task
	for i := 0; i < taskNum; {
		queue = p.queue.roundRobinPick()
		select {
		case queue <- task[i]:
			i++
		default:
		}
	}
	return nil
}

// 使用wg保证worker的同步启动
func (p *pool) start() {
	wg := &sync.WaitGroup{}
	for i := 0; i < p.opt.workerNum; i++ {
		wg.Add(1)
		go p.worker(i, wg)
	}
	wg.Wait()
}

func (p *pool) worker(index int, wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case task, ok := <-p.queue.receive(index):
			if !ok {
				return
			}
			// 使用channel的阻塞来控制Goroutine的数量在设置范围
			p.numGoroutine <- struct{}{}
			go p.exec(index, task)
		}
	}
}

func (p *pool) exec(index int, t *Task) {
	t.Exec()
	p.wait.Done()
	<-p.numGoroutine
}

func getPool(name string) (*pool, error) {
	item, ok := pools.Load(name)
	if !ok {
		return nil, PoolNotFound
	}
	return item.(*pool), nil
}
