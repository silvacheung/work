package task

type queue struct {
	queues []chan *Task
	pickCh chan chan *Task
}

func initQueue(workerNum, goroutineNum int) *queue {

	queue := &queue{
		queues: make([]chan *Task, workerNum),
		pickCh: make(chan chan *Task, workerNum),
	}

	// 保证能够被workerNum整除,这里差多少就补多少
	// 这里只是控制队列的大小,并不是实际的协程数,所以多点少点没关系
	if diff := goroutineNum % workerNum; diff != 0 {
		goroutineNum += diff
	}

	// 给每个worker分配队列
	for i := 0; i < workerNum; i++ {
		item := make(chan *Task, goroutineNum/workerNum)
		queue.queues[i] = item
		queue.pickCh <- item
	}

	return queue
}

func (q *queue) receive(index int) <-chan *Task {
	return q.queues[index]
}

func (q *queue) roundRobinPick() chan<- *Task {
	queue := <-q.pickCh
	q.pickCh <- queue
	return queue
}

func (q *queue) close() {
	for _, queue := range q.queues {
		close(queue)
	}
	close(q.pickCh)
}
