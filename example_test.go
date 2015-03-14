package gofidlib

import (
	"fmt"
)

func ExampleFilter_Run() {
	design := NewFilterDesign("HpBeZ1/2", 4)
	filt := NewFilter(design)

	in := make([]float64, 4)
	out := make([]float64, 4)
	for i := range in {
		in[i] = float64(i)
		out[i] = filt.Run(in[i])
	}

	filt.Close()

	fmt.Println("IN:", in)
	fmt.Println("OUT:", out)
	// Output:
	// IN: [0 1 2 3]
	// OUT: [0 0.5 0.5 0.5]
}
