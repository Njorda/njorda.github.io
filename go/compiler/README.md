# Copilie optimisation using pprof

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