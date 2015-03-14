package gofidlib

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L .


#include <stdio.h>
#include "fidlib.h"
double fid_run(FidFunc *funcp, void *fbuf1, double in)
{
	return funcp(fbuf1, in);
}
*/
import "C"

import "unsafe"

type Filter struct {
	design *FilterDesign
	buf    unsafe.Pointer
}

//Run the filter on an input array of signal data.
//Fills the output array with filter signal data,
//in array is not altered.
func (f *Filter) Run(in float64) float64 {
	return float64(C.fid_run(f.design.funcp, f.buf, C.double(in)))
}

//Reset the filter buffer
func (f *Filter) Zap() {
	C.fid_run_zapbuf(f.buf)
}

//Deallocate the run and buffer objects
//For use when filtering complete
func (f *Filter) Close() {
	C.fid_run_free(f.design.fidRun)
	C.fid_run_freebuf(f.buf)
}

func NewFilter(design *FilterDesign) *Filter {
	buf := C.fid_run_newbuf(design.fidRun)
	return &Filter{
		design: design,
		buf:    buf,
	}
}
