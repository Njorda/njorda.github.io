---
layout: post
title: "Go memory arenas for apache arrow, Part 1"
subtitle: "Try to implement a new memory allocator for apache arrow (Go)"
date: 2023-03-29
author: "Niklas Hansson"
URL: "/2023/03/23/apache-arrow-memory-arena-go"
---

This blog post will try to dive down in to [apache arrow](https://arrow.apache.org/) and specifically the [Go memory](https://github.com/apache/arrow/tree/main/go) allocation for Apache arrow. Apache arrow state that they allow for the following types of memory allocations: 

- Go default allocations(standard go GC collected memory)
- CGo allocator(memory allocated through CG0)
- Checked Memory Allocator

Will deep dive in to these once in a follow up blog post. Today the goal is to extend with a new memory allocator, mostly becuase I read up on go memory [arena](https://go.dev/src/arena/arena.go) which where introduced in to 1.20 under the experimental flag( to access it you need to run go 1.20 or later and have `export GOEXPERIMENT=arenas`). 

The first question could be why do we do this? Mostly just due to the fact that memory arenas are a interesting concept for Go and I want to play around with it. However I would also like to see if it could give some performance benefits reducing the amount of GC and allocations. But mostly it is just to try it out and have fun. 


# What is Memory `arenas`?

>The arena package provides the ability to allocate memory for a collection
of Go values and free that space manually all at once, safely. The purpose
of this functionality is to improve efficiency: manually freeing memory
before a garbage collection delays that cycle. Less frequent cycles means
the CPU cost of the garbage collector is incurred less frequently.

In short it allows for allocating a chunk of memory and then utalise that chunk to store objects. The chunk will not be garbage collected but needs to be freed by the user. Referencing variables that lives in the memory arena after the arena has been freed will result in a panic, and thus care needs to be handled by the users. This also seems to be one of the[major concern](https://github.com/golang/go/issues/51317) from the go community about this feature. 

There is also a short blog post [Go memory arena]({{< ref "/post/2023-03-02-memory-arena" >}}) specifically about memory arenas. 

# Update apache arrow. 

So the first step is to fork the [apache arrow repo](https://github.com/apache/arrow) my fork live [here](https://github.com/NikeNano/arrow), currently the changes are in the `nikenano/MemoryArena` branch with the changes updates inside [arrow/go/arrow/memory
/go_allocator_memory_arena.go](https://github.com/NikeNano/arrow/blob/nikenano/MemoryArena/go/arrow/memory/go_allocator_memory_arena.go). The implementation is very similar to the Go default implementation(please let me know if you find any errors, most of this is new and I try to learn on the fly). 

# Test it out

Based upon this example from [Voltron data](https://voltrondata.com/resources/use-apache-arrow-and-go-for-your-data-workflows) we will try out our new allocator and see if we can dive down in to the memory profiling in go. We will take the first blog post and replace the memory allocator. 


First step is to reproduce the examples from the blog. 

according to the instructions the first step is to import the following package

```bash
go get -u github.com/apache/arrow/go/v11@latest
```

and the code is the following: 

```go
package examples_test

import (
    "fmt"

    "github.com/apache/arrow/go/v10/arrow"
    "github.com/apache/arrow/go/v10/arrow/array"
    "github.com/apache/arrow/go/v10/arrow/memory"
)

func Example_buildInt64() {
    bldr := array.NewInt64Builder(memory.DefaultAllocator)
    defer bldr.Release() // <-- Notice This!

    bldr.Append(25) // append single value
    bldr.AppendNull() // append a null value to the array
    // Append a slice of values with a slice of booleans
    // defining which ones are valid or not.
    bldr.AppendValues([]int64{1, 2}, []bool{true, false})
    // Or pass nil to assume all are valid
    bldr.AppendValues([]int64{3, 4, 5, 6}, nil)

    arr := bldr.NewInt64Array() // can be reused after this
    defer arr.Release() // <-- Notice!
    fmt.Println(arr)

    bldr.Append([]int64{7, 8, 9, 10}, nil)
    // get Array interface rather than typed pointer
    arr2 := bldr.NewArray()
    defer arr2.Release() // <-- Seeing the pattern?
    fmt.Println(arr2)

    // Output:
    // [25 (null) 1 (null) 3 4 5 6]
    // [7 8 9 10]
}
```

However this gave me the following issue: 


```bash
$ go run main.go
../../../../../go/pkg/mod/github.com/apache/arrow/go/v11@v11.0.0/parquet/internal/gen-go/parquet/parquet-consts.go:10:2: missing go.sum entry for module providing package github.com/apache/thrift/lib/go/thrift (imported by github.com/apache/arrow/go/v11/parquet/internal/gen-go/parquet); to add:
	go get github.com/apache/arrow/go/v11/parquet/internal/gen-go/parquet@v11.0.0
../../../../../go/pkg/mod/github.com/apache/arrow/go/v11@v11.0.0/arrow/array/binary.go:26:2: missing go.sum entry for module providing package github.com/goccy/go-json (imported by github.com/apache/arrow/go/v11/arrow/array); to add:
	go get github.com/apache/arrow/go/v11/arrow/array@v11.0.0
../../../../../go/pkg/mod/github.com/apache/arrow/go/v11@v11.0.0/parquet/compress/snappy.go:23:2: missing go.sum entry for module providing package github.com/golang/snappy (imported by github.com/apache/arrow/go/v11/parquet/compress); to add:
	go get github.com/apache/arrow/go/v11/parquet/compress@v11.0.0
../../../../../go/pkg/mod/github.com/apache/arrow/go/v11@v11.0.0/arrow/internal/flatbuf/Binary.go:22:2: missing go.sum entry for module providing package github.com/google/flatbuffers/go (imported by github.com/apache/arrow/go/v11/arrow/internal/flatbuf); to add:
	go get github.com/apache/arrow/go/v11/arrow/internal/flatbuf@v11.0.0
../../../../../go/pkg/mod/github.com/apache/arrow/go/v11@v11.0.0/internal/hashing/xxh3_memo_table.go:31:2: missing go.sum entry for module providing package github.com/zeebo/xxh3 (imported by github.com/apache/arrow/go/v11/internal/hashing); to add:
	go get github.com/apache/arrow/go/v11/internal/hashing@v11.0.0
../../../../../go/pkg/mod/github.com/apache/arrow/go/v11@v11.0.0/arrow/datatype_fixedwidth.go:25:2: missing go.sum entry for module providing package golang.org/x/xerrors (imported by github.com/apache/arrow/go/v11/arrow); to add:
	go get github.com/apache/arrow/go/v11/arrow@v11.0.0
```

this is due to the mix between the v10 and v11 reference in the go package. To solved this update the imports to `"github.com/apache/arrow/go/v11/PACKAGE_OF_INTEREST"` =>

```go
    "github.com/apache/arrow/go/v11/arrow"
    "github.com/apache/arrow/go/v11/arrow/array"
    "github.com/apache/arrow/go/v11/arrow/memory"
```

This make the code partialy run but we hit the next issue: 

```
$ go run main.go
# command-line-arguments
./main.go:26:36: too many arguments in call to bldr.Append
	have ([]int64, nil)
	want (int64)
```

The issue is probably due to that that `bldr.Append([]int64{7, 8, 9, 10}, nil)` should be `bldr.AppendValues([]int64{7, 8, 9, 10}, nil)`


# Go memory arena

The implementation for using memory arenas lives [here](https://github.com/NikeNano/arrow/blob/nikenano/MemoryArena/go/arrow/memory/go_allocator_memory_arena.go). Where the key is to generate a memory arena which is used to allocate all the memory and only free the arena when we want to free all the memory. There is no way to partly free memory from a memory arena. 

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

```

a couple of tests are also added. 

# Check memory profiles

In order to check the memory usage we add the following to the example above, we also add a loop and break it out to a function in order for the profile to have time to profile it: 


```go
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var mem runtime.MemStats

	log.Println("memory baseline...")

	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
```

Which gives the following code: 

```go 
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/apache/arrow/go/v11/arrow"
	"github.com/apache/arrow/go/v11/arrow/array"
	"github.com/apache/arrow/go/v11/arrow/memory"
)

func main() {
	f, err := os.Create("profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var mem runtime.MemStats

	log.Println("memory baseline...")

	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
	now := time.Now()
	alloc := memory.NewGoAllocator()
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

Looking at the heap allocations the memory arena is surprisingly high and something seems to be wrong. However this is for a second blog post. 