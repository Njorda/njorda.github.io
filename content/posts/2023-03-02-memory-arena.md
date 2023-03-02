# ---
layout:     post 
title:      "Go memory arena"
subtitle:   "Optimize the GC by not using it"
date:       2023-03-01
author:     "Niklas Hansson"
URL: "/2023/03/01/go_memory_arena/"
iframe: "https://nikenano.github.io/shinylive/"
---



As part of the go 1.20 release [memory areas](https://go.dev/src/arena/arena.go) where introduced to the standard lib but not mentioned in the [release notes](https://go.dev/doc/go1.20) but is still being discussed as of 2023-03-02 [here](https://github.com/golang/go/issues/51317). Memory arenas allow users to allocate memory and are described by the docs as: 

> The arena package provides the ability to allocate memory for a collection
of Go values and free that space manually all at once, safely. The purpose
of this functionality is to improve efficiency: manually freeing memory
before a garbage collection delays that cycle. Less frequent cycles means
the CPU cost of the garbage collector is incurred less frequently.

But why would one do this? Go is garbage collected so why is the need to get rid of that functionality for certain objects of your code? Garbage Collection(GC) comes with a certain over head. Depending upon your program this could be a substantial cost in terms of CPU which memory arenas could help reduce at the cost or manually managing the memory. The `memory arenas` as en example of [region based memory](https://en.wikipedia.org/wiki/Region-based_memory_management). Where a region( aka area) can be reallocated or deallocated all at once. Often region and areas are implemented such that all objects in a area are allocated in a [single contiguous range](https://en.wikipedia.org/wiki/Region-based_memory_management) of memory(same as stack frames). It should be noted that areas should NEVER be accessed by multiple `goroutines`. 


Areas are only avilable in go 1.20 when `experimental` is used. 

```
GOEXPERIMENT=arenas go run main.go
```

```
import "arena"


func process() {
	// Create an arena
	mem := arena.NewArena()
	// Free the arena in the end.
	defer mem.Free()

	// Or a slice with length and capacity.
	slice := arena.MakeSlice[T](mem, 100, 200)
}

```


if a slice out grows the capacity it will be moved to the `heap` if not reallocated. 

I will play around and try to see if I can validate the performance boost that [google reported of 15%](https://github.com/golang/go/issues/51317)
