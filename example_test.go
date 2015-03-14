package gofidlib

import (
	"fmt"
)

func ExampleFilter_Run() {
	design, _ := NewFilterDesign("HpBeZ1/2", 4)
	filt := NewFilter(design)

	in := make([]float64, 4)
	out := make([]float64, 4)
	for i := range in {
		in[i] = float64(i)
		out[i] = filt.Run(in[i])
	}

	fmt.Println("IN:", in)
	fmt.Println("OUT:", out)

	filt.Zap()

	for i := range in {
		in[i] = float64(i)
		out[i] = filt.Run(in[i])
	}
	fmt.Println("IN:", in)
	fmt.Println("OUT:", out)
	// Output:
	// IN: [0 1 2 3]
	// OUT: [0 0.5 0.5 0.5]
	// IN: [0 1 2 3]
	// OUT: [0 0.5 0.5 0.5]
}

func Example_FidCat() {
	design1, _ := NewFilterDesign("HpBe1/1", 250)
	design2, _ := NewFilterDesign("HpBu1/1", 250)
	design3, _ := NewFilterDesign("LpBe1/50", 250)
	design4, _ := NewFilterDesign("LpBu1/50", 250)
	designs := []*FilterDesign{design1, design2, design3, design4}
	fmt.Println("IN:", design1.fidFilter.typ, design2.fidFilter.typ, design3.fidFilter.typ, design4.fidFilter.typ)
	design, _ := FidCat(1, designs) //concatinate the filters into one filter design, and free prev
	fmt.Println("OUT:", design.fidFilter.typ)
	// Output:
	// IN: 70 70 70 70
	// OUT: 70
}


