---
layout: post
title: "Gorutines - deep dive into Go:s concurrency features"
subtitle: "What makes Gorutines so great"
date: 2023-08-15
author: "Niklas Hansson"
URL: "/2023/08/15/gorutines"
---


# Gorutine high level

Go:s native concurrency features is usually one of the first things developers bring up when describing the advantages of using go. However due to its eas of use it is also wildly miss used. 

> “You can have a second computer once you’ve
shown you know how to use the first one.”
> –Paul Barham

We will not deep dive in to when to use and not to use `Gorutines` specifically even though we touch upon some hints. However concurrency introduces a lot of complexity and makes code both harder to read and debug - so only use it when needed. 

`Gorutines` are a concept in Go for starting a concurrent(or possibly parallel function depending upon the available resources). From the definition in [Effective Go
](https://go.dev/doc/effective_go#concurrency) we can read: 

>  A goroutine has a simple model: it is a function executing concurrently with other goroutines in the same address space. 

But what does this mean? The first thing that is good to know is that the main function of a go program is a `Gorutine` and creating `Gorutines` using the keyword `go` is the same as the main `Gorutine`. `Gorutines` are multiplexed on to  OS threads. This means that multiple `Gorutines` are run on the same thread. If one function call blocks the OS thread the go runtime will automatically move `Gorutines` to other OS threads in order to unblock them. The goal with `Gorutines` is to make concurrency easier for to use while having a high performance for concurrent task avoiding the cost of creating a large amount of OS threads(which is expensive to create and delete). `Gorutines` have a low overhead beyond the overhead of the memory of the gorutine stack. In order to make the stack as small as possible Go uses a runtime resizable, bounded stacks. This means that the stack is allow grow. When the stack is reaching is max, it is copied to a larger memory section allowing it to grow further(the stack is also allowed to shrink). It should be noticed that Go copies the whole stack to a new location rather than allowing it to live in multiple sections of memory. `Gorutines` are described by [Dave Chenny](https://dave.cheney.net/tag/goroutines) as: 

> Go’s goroutines sit right in the middle, the same programmer interface as threads, a nice imperative coding model, but also efficient implementation model based on coroutines.

Gorutines and the runtime are created to handle a large amount(think 100k or more) of `Gorutines` but even though the memory footprint is low it is still something. Thus we can not have an infinite number of them, due to this spinning up `Gorutines` that don't finish will create a potential memory leak. This result in a very important consideration when creating `Gorutines` - when will it stop - in order to avoid leaking memory. When using `Gorutines` always be aware of when they start or finish.

# Runtime scheduling

The Go program will show up as a single process requesting and running multiple threads. The runtime scheduler is responsible for how the `Gorutines` are scheduled on to the OS threads. Go uses the [`GOMAXPROCS`](https://cs.opensource.google/go/go/+/go1.21.0:src/runtime/debug.go;l=16) env to set the maximum nbr of CPUs that can be used. However if your processor has [hyper-threading](https://en.wikipedia.org/wiki/Hyper-threading) each hardware thread will be considered one process(this is called P by the go runtime). The go scheduler is a [cooperative scheduler](https://en.wikipedia.org/wiki/Cooperative_multitasking) which in short means that each the runtime will not initialize a context switch but the `Gorutines` will yield to other `Gorutines`. `Gorutines` has three high level states: 

- Waiting - Refers to that a `Gorutine` is waiting for something, such as synchronization calls or OS calls. 
- Runnable - A `Gorutine` is ready to run and wants time on a thread in order to execute its instructions. 
- Executing - A `Gorutine` is running on a thread.

Go uses three abstractions for the scheduling of `Gorutines`, more info can be found [here](https://go.dev/src/runtime/HACKING): 

- G - stands for `Gorutine` 
- M - is a OS thread
- P - is resource required to execute Go code(for example memory allocated space). The number of P:s are equal to GOMAXPROCS(which is default to the nbr of threads tha can be run on the cores). 

The schedulers are responsible for matching the Gorutines, G to the P:s which then matches it on to the M. We stated that OS threads are expensive to create and destroy and thus we only want to create a limited nbr. The go runtime will only create OS threads to match P, however if a OS threads is blocked the runtime will look for free threads or create a new. 

In order to allow for a large amount of `Gorutines` each P has a queue of G:s ready to run. These are then handled by the P or stolen using `work stealing`(steels half of some other P:s queue) in order to balance the load evenly by other P:s. For long running G:s that don't have synchronization points the G:s are moved to a global queue that are given a lower priority in order to not block other G:s. Also if the distributed run queues for each P is full the G:s will be added to the global run queue. 

# Resources: 

- [Go docs](https://go.dev/doc/faq#goroutines)
- [Effective Go](https://go.dev/doc/effective_go#goroutines)
- [A deep dive in to Go's stack (Go Time Live!)](https://www.youtube.com/watch?v=FQNt-qp7FpQ)
- [Never start a goroutine without knowing how it will stop](https://dave.cheney.net/tag/goroutines)
- [Scheduling In Go : Part II - Go Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)
- [Go scheduler: Ms, Ps & Gs](https://povilasv.me/go-scheduler/)
- [Go runtime scheduler](https://speakerdeck.com/retervision/go-runtime-scheduler)
- [GopherCon 2018: Kavya Joshi - The Scheduler Saga](https://www.youtube.com/watch?v=YHRO5WQGh0k)

<!-- http://www.gotw.ca/publications/concurrency-ddj.htm -->