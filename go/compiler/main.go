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
	for i := 0; i < 100000; i++ {
		arrowFunc()
	}
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
