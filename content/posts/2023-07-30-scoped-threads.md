---
layout: post
title: "Rust: Scoped threads - easier multithreading"
subtitle: "How, why and when you should use scoped threads"
date: 2023-07-30
author: "Niklas Hansson"
URL: "/2023/07/30/rust-scoped-threads"
---

The current development in CPU design is going towards large amount of cores rather than faster cores and thus writing parallel code becomes more important in order to utilize the full potential ([Concurrency is not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg)). In this blog post we will dive into scoped threads, what it is and what is the difference between threads in rust in general. First of all only use threads if you need the speed up, introducing threads to a program adds complexity which both makes the program harder to maintain but if not done correct also slower to run(due to communications between threads and scheduling). 

Assuming threads is the way to go in your rust program the next step is to understand what threads are in rust. Compare to `go` which has the concept of `goroutines` which is a layer on top of operatins system threads(where multiple go routines are multiplexed on to a single os thread, super cool concept that we will dive deeper in to in a separate blog post) rust kicks of a os thread per thread created with the `use std::thread;` create. Is should thus be noted that for smaller tasks this will not be worth to do, only for larger chunks of work will this result in a speed up of the the program. 

To start a thread in rust(example from [the book](https://doc.rust-lang.org/book/ch16-01-threads.html)): 

```rust
use std::thread;
use std::time::Duration;

fn main() {
    thread::spawn(|| {
        for i in 1..10 {
            println!("hi number {} from the spawned thread!", i);
            thread::sleep(Duration::from_millis(1));
        }
    });

    for i in 1..5 {
        println!("hi number {} from the main thread!", i);
        thread::sleep(Duration::from_millis(1));
    }
}
```

This will output: 

```bash
hi number 1 from the main thread!
hi number 1 from the spawned thread!
hi number 2 from the main thread!
hi number 2 from the spawned thread!
hi number 3 from the main thread!
hi number 3 from the spawned thread!
hi number 4 from the main thread!
hi number 4 from the spawned thread!
hi number 5 from the spawned thread!
```

However it should be noted that there is nothing in here that guarantees that the main thread will wait for the spawned once. This could be handled a `handle` could be used, more info [here](https://doc.rust-lang.org/book/ch16-01-threads.html). 

So to the point of scoped threads what is the difference and why does it matter? But maybe first of all what is scope in rust? Shortly variable scope can be described as the part of the code where a variable can be accessed. Rust implements what is known as [Resource acquisition is initialization](https://en.wikipedia.org/wiki/Resource_acquisition_is_initialization) or just RAII for short, which means that variable in Rust not only hold the data, but also owns the resource. The main advantage of RAII is that it encapsulates of resources by tying the resource lifetime to a stack variable. When a variable goes out of scope(not accessable any longer in the program) the resources are freed and thus as long as we avoid leaking variable we avoid leaking resources. 

