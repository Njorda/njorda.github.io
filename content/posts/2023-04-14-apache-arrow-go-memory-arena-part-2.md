---
layout: post
title: "Go memory arenas for apache arrow, Part 2"
subtitle: "Try to implement a new memory allocator for apache arrow (Go)"
date: 2023-04-14
author: "Niklas Hansson"
URL: "/2023/04/14/apache-arrow-memory-arena-go-part-2"
---

This blog post will continue to try to dive down in to [apache arrow](https://arrow.apache.org/) and specifically the [Go memory](https://github.com/apache/arrow/tree/main/go) allocation for Apache arrow. This is a follow up to [Go memory arenas for apache arrow, Part 1]({{< ref "/posts/2023-03-29-apache-arrow-go-memory-arena.md" >}}). 

First of all why do we want to manage memory manually instead of using the GC? One of arrows key features is it support to share memory with out copy between programs however for a GC collected language this will not work that great. What happens if Go believes the memory is not used any longer(which is correct for the go part) but used by someone else outside of go.  In these situations it would be beneficial to be able to control when the memory is released. 

The implementation is: 

```go
package memory

import (
	"arena" // requires export GOEXPERIMENT=arenas to be set
	"sync"
)

type GoArenaAllocator struct {
	mem *arena.Arena
	// Keep track on all the allocations, when all use then we can call free.
	// map with the allocations which we need, I think this would be awesome.
	addrs map[int]bool
	sync.Mutex
}

func NewGoArenaAllocator() *GoArenaAllocator {
	return &GoArenaAllocator{arena.NewArena(), map[int]bool{}, sync.Mutex{}}
}

func (a *GoArenaAllocator) Allocate(size int) []byte {
	buf := arena.MakeSlice[byte](a.mem, size+alignment, size+alignment) // padding for 64-byte alignment, I dont think this is needed in the arena since we make all 64 bit aligned
	addr := int(addressOf(buf))
	// So data will be loaded based upon division with 64, here we check the address pointer.
	// If the data is even division with 64 we can load it to the cache way more efficient and gain speed ups
	// What we do here is move ths buffer around so the address is has a start that is even with 64 so we can load it faster.
	//
	next := roundUpToMultipleOf64(addr)
	a.Lock()
	defer a.Unlock()
	if addr != next {
		shift := next - addr
		out := buf[shift : size+shift : size+shift]
		addr := int(addressOf(out))
		a.addrs[addr] = true
		return out
	}
	a.addrs[addr] = true
	return buf
}

func (a *GoArenaAllocator) CheckSize() int {
	return len(a.addrs)
}

func (a *GoArenaAllocator) Reallocate(size int, b []byte) []byte {
	if size == len(b) {
		return b
	}
	newBuf := a.Allocate(size)
	copy(newBuf, b)
	return newBuf
}

func (a *GoArenaAllocator) Free(b []byte) {
	addr := int(addressOf(b))
	a.Lock()
	delete(a.addrs, addr)
	a.Unlock()
	if len(a.addrs) > 0 {
		return
	}
	a.mem.Free()
}

var _ Allocator = &GoArenaAllocator{}
```


Memory arenas will only release all the memory in one go, but allow for allocating new memory to the same arena. Thus in order to make sure that we all memory can be released a map is used with where the pointer address is stored for book keeping. When a specific piece of memory is released this is removed from the map. When the last piece of memory is released in the arena the arena is deleted, an arena can not be reuased when it is deleted. However the `GoArenaAllocator` is thread safe to use(this is supported through a `sync.Mutex`. 



# Byte alignment

In order for a CPU to work more effective with the memory if the data is aligned with the bits in the architecture of the CPU. Thus a 64 bit CPU will be more performant if the data is stored in 8 consecutive bytes and the first byte lis on a 8 byte boundary.  

>A memory address a is said to be n-byte aligned when a is a multiple of n (where n is a power of 2). In this context, a byte is the smallest unit of memory access, i.e. each memory address specifies a different byte

But why does this matter? Every fetch from memory will fetch a cache line usually 64 bits(not related to 64 bit processors). Thus if the data is aligned with this size in mind we will reduce spill over between cache lines and thus reduce the number of cache lines that needs to be fetched. Memory modifications are also effecting cache lines since it might result in that other CPUs have to fetch the same line again and. Accidentally placing unrelated data on the same cache line is known as "false sharding"

For the interested reader [this is a great paper.](https://www.akkadia.org/drepper/cpumemory.pdf). 

an example program using the new arena allocator: 

```go 
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
)

func main() {
	f, err := os.Create("profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	now := time.Now()
	alloc := memory.NewGoArenaAllocator()
	var mem runtime.MemStats

	log.Println("memory baseline...")

	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
	for i := 0; i < 1000000; i++ {
		run(alloc)
	}
	fmt.Println(time.Since(now))

	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
}

func run(alloc memory.Allocator) {
	lb := array.NewFixedSizeListBuilder(alloc, 3, arrow.PrimitiveTypes.Int64)

	defer lb.Release()

	vb := lb.ValueBuilder().(*array.Int64Builder)
	vb.Reserve(10)

	lb.Append(true)
	vb.Append(0)
	vb.Append(1)
	vb.Append(2)

	lb.AppendNull()
	vb.AppendValues([]int64{-1, -1, -1}, nil)

	lb.Append(true)
	vb.Append(3)
	vb.Append(4)
	vb.Append(5)

	lb.Append(true)
	vb.Append(6)
	vb.Append(7)
	vb.Append(8)

	lb.AppendNull()

	arr := lb.NewArray().(*array.FixedSizeList)
	defer arr.Release()

	// Output:
	// [25 (null) 1 (null) 3 4 5 6]
	// [7 8 9 10]
}

```




```bash
$ go run main.go
2023/04/14 23:59:58 memory baseline...
2023/04/14 23:59:58 10190192
2023/04/14 23:59:58 10190192
2023/04/14 23:59:58 10190192
2023/04/14 23:59:58 12156928
2.218185631s
2023/04/15 00:00:00 892640320
2023/04/15 00:00:00 1975439680
2023/04/15 00:00:00 892640320
2023/04/15 00:00:00 1119322112
```


Replacing `	alloc := memory.NewGoArenaAllocator()` -> `	alloc := memory.NewGoAllocator()` gives the following response. 

```bash
$ go run main.go
2023/04/15 00:04:37 memory baseline...
2023/04/15 00:04:37 1788424
2023/04/15 00:04:37 1788424
2023/04/15 00:04:37 1788424
2023/04/15 00:04:37 3932160
1.285498898s
2023/04/15 00:04:38 1788424
2023/04/15 00:04:38 1788424
2023/04/15 00:04:38 1788424
2023/04/15 00:04:38 3932160
1.285526056s
```

Shows that in when using the `GoAllocator` we put everything on the stack while the `GoArenaAllocator` makes a large amount of heap allocations. 

# Tests

Which the following tests, 

```go 
package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requires export GOEXPERIMENT=arenas to be set

func TestMemory(t *testing.T) {
	a := NewGoArenaAllocator()
	require.Equal(t, 0, a.CheckSize())
	s1 := a.Allocate(10)
	require.Equal(t, 1, a.CheckSize())
	s2 := a.Allocate(11)
	require.Equal(t, 2, a.CheckSize())
	a.Free(s1)
	require.Equal(t, 1, a.CheckSize())
	a.Free(s2)
	require.Equal(t, 0, a.CheckSize())
}

func TestNewGoArenaAllocator_Allocate(t *testing.T) {
	tests := []struct {
		name string
		sz   int
	}{
		{"lt alignment", 33},
		{"gt alignment unaligned", 65},
		{"eq alignment", 64},
		{"large unaligned", 4097},
		{"large aligned", 8192},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alloc := NewGoArenaAllocator()
			buf := alloc.Allocate(test.sz)
			assert.NotNil(t, buf)
			assert.Len(t, buf, test.sz)
			defer alloc.Free(buf)
		})
	}
}

func TestGoArenaAllocator_Reallocate(t *testing.T) {
	tests := []struct {
		name     string
		sz1, sz2 int
	}{
		{"smaller", 200, 100},
		{"same", 200, 200},
		{"larger", 200, 300},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alloc := NewGoArenaAllocator()
			buf := alloc.Allocate(test.sz1)
			for i := range buf {
				buf[i] = byte(i & 0xFF)
			}

			exp := make([]byte, test.sz2)
			copy(exp, buf)

			newBuf := alloc.Reallocate(test.sz2, buf)
			assert.Equal(t, exp, newBuf)
			defer alloc.Free(newBuf)
		})
	}
}

```

# Links

- https://stackoverflow.com/questions/34860366/why-buffers-should-be-aligned-on-64-byte-boundary-for-best-performance
- https://en.wikipedia.org/wiki/Data_structure_alignment#:~:text=A%20memory%20address%20a%20is,address%20specifies%20a%20different%20byte.