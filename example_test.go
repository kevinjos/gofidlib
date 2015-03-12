package gofidlib

import (
	"fmt"
)

func ExampleFilter_Run() {
	in := make([]float64, 4)
	out := make([]float64, 4)
	for i := range in {
		in[i] = float64(i)
	}

	design := NewFilterDesign("HpBeZ1/2", 4)
	filt := NewFilter(design)

	filt.Run(in, out)

	fmt.Println("IN:", in)
	fmt.Println("OUT:", out)
	// Output:
	// IN: [0 1 2 3]
	// OUT: [0 0.5 0.5 0.5]
}
