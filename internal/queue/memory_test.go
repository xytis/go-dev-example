package queue

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

func TestMemoryGCEatsUnderlyingArrayWhileAppending(t *testing.T) {
	q := NewArrayQueue()

	msg := Message("o")
	cork := 100
	limit := 100_000_000

	var mem runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&mem)
	before := mem.HeapInuse

	for i := 0; i < cork; i++ {
		_ = q.Push(msg)
	}

	for i := 0; i < limit; i++ {
		_ = q.Push(msg)
		_, _ = q.Pull()
	}

	runtime.GC()
	runtime.ReadMemStats(&mem)
	after := mem.HeapInuse

	// Observed behaviour: before ~= after
	fmt.Println(before, after)
	fmt.Println(math.Abs(float64(after) - float64(before)))
}

func TestMemoryGCEatsUnderlyingArrayEvenWithoutAppend(t *testing.T) {
	q := NewArrayQueue()

	msg := Message("o")
	limit := 100_000_000
	remainder := 100

	var mem runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&mem)
	before := mem.HeapInuse

	for i := 0; i < limit; i++ {
		_ = q.Push(msg)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem)
	during := mem.HeapInuse

	for i := 0; i < limit-remainder; i++ {
		_, _ = q.Pull()
	}

	runtime.GC()
	runtime.ReadMemStats(&mem)
	after := mem.HeapInuse

	// Observed behaviour: before ~= after
	fmt.Println(before, during, after)
	fmt.Println(math.Abs(float64(after) - float64(before)))
}
