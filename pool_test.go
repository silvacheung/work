package task

import (
	"runtime"
	"testing"
)

var task = New("empty-task", func() {
	// empty task for benchmark
})

func init() {
	StartPool(WithPoolName(DefaultPool), WithGoroutineNum(1000), WithWorkerNum(runtime.NumCPU()))
	//StopPool(DefaultPool)
}

// BenchmarkPut-4           3000000               412 ns/op               0 B/op          0 allocs/op
func BenchmarkPut(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := Put(DefaultPool, task); err != nil {
			b.Errorf("Task error : %v \n", err)
			b.FailNow()
		}
	}
}
