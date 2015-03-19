package gofidlib

import "testing"

type TestDesign struct {
	inSpec string
	inRate float64
	out bool
}

var testDesigns = []TestDesign{
	{"HpBe1/1", 250, true},
	{"HpBe1/10", 5, false},
}

func TestFilterDesign(t *testing.T) {
	for i, test := range testDesigns {
		res := true
		_, err := NewFilterDesign(test.inSpec, test.inRate)
		if err != nil {
			res = false
		}
		if res != test.out {
			t.Errorf("#%d: Response(%v)=(%v, %v); want %v", i, test.inSpec, test.inRate, res, test.out)	
		}
	}
}

type TestR struct {
	in float64
	out float64
}

var testRs = []TestR{
	{2.5, 1},
	{1.5, 1},
}

func TestFilterResponse(t *testing.T) {
	design, _ := NewFilterDesign("HpBeZ1/2", 8)
	filt := NewFilter(design)
	defer func() {
		design.Free()
		filt.Free()
	}()
	for i, test := range testRs {
		resp := filt.Response(test.in)
		if resp != test.out {
			t.Errorf("#%d: Response(%v)=%v; want %v", i, test.in, resp, test.out)	
		}
	}
}

type TestPR struct {
	in float64
	out []float64
}

var testPRs = []TestPR{
	{2.5, []float64{1.0, 1.0}},
	{1.5, []float64{1.0, 1.0}},
}

func TestFilterPhaseResponse(t *testing.T) {
	design, _ := NewFilterDesign("HpBeZ1/2", 8)
	filt := NewFilter(design)
	defer func() {
		design.Free()
		filt.Free()
	}()

	for i, test := range testPRs {
		resp, phase := filt.PhaseResponse(test.in)
		if resp != test.out[0] || phase != test.out[1] {
			t.Errorf("#%d: PhaseResponse(%v)=%v; want %v", i, test.in, resp, phase, test.out)	
		}
	}
}

type TestD struct {
	in string
	out int
}

var testDs = []TestD{
	{"HpBeZ1/2", 1},
	{"HpBuZ1/2", 1},
	{"LpBeZ1/2", 0},
	{"LpBuZ1/2", 0},
	{"HpBe1/2", 0},
	{"HpBu1/2", 0},
}


func TestFilterDelay(t *testing.T) {
	for i, test := range testDs {
		design, _ := NewFilterDesign(test.in, 8)
		filt := NewFilter(design)
		delay := filt.Delay()
		if delay != test.out {
			t.Errorf("#%d: Delay(%s)=%v; want %v", i, test.in, delay, test.out)	
		}
		filt.Free()
		design.Free()
	}
}
