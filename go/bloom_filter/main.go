package main

import (
	"fmt"

	"github.com/Njorda/go/bloom_filter/filter"
)

func main() {
	println("Start of bloom filter")
	bloom := filter.New(2, 1000)
	if err := bloom.Add([]byte("Thomas")); err != nil {
		panic(err)
	}
	if err := bloom.Add([]byte("filip")); err != nil {
		panic(err)
	}
	if err := bloom.Add([]byte("daniel")); err != nil {
		panic(err)
	}
	exsist, err := bloom.Check([]byte("Tommy"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Tommy exists: %v\n", exsist)
}
