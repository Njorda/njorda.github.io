---
layout:     post 
title:      "Go compiler optimizations"
subtitle:   "Profile-guided inlining optimization "
date:       2023-02-08
author:     "Niklas Hansson"
URL: "/2023/02/08/Profile-guided_inlining_optimization/"
iframe: "https://nikenano.github.io/shinylive/"
---


This is based upon the new feature released in [go v1.20](https://tip.golang.org/doc/go1.20) where the compiler can optimize using a pprof file. 

In order to run the pprof we will use flags: 

```go
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
```

more info can be found [here](https://go.dev/blog/pprof). 

In order to run the profiling use the following command: 

```bash
go run main.go -cpuprofile=prof.prof
```

The specific information related to how the complier will optimize the code can be found [here](https://tip.golang.org/doc/go1.20#compiler). But why is this relevant? It has shown that this can allow for inlining to be optimized: 

> Benchmarks for a representative set of Go programs show enabling profile-guided inlining optimization improves performance about 3–4%.
> 

```bash
go build -pgo=prof.prof
```

to run the code use: 

```
./compiler
```

The complete code: 



```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/memory"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	arrowFunc()
}

func arrowFunc() {
	pool := memory.NewGoAllocator()

	lb := array.NewFixedSizeListBuilder(pool, 3, arrow.PrimitiveTypes.Int64)
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

	fmt.Printf("NullN()   = %d\n", arr.NullN())
	fmt.Printf("Len()     = %d\n", arr.Len())
	fmt.Printf("Type()    = %v\n", arr.DataType())
	fmt.Printf("List      = %v\n", arr)
}
```