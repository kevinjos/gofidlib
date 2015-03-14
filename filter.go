package gofidlib

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L .


#include <stdio.h>
#include <stdlib.h>
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

func (f *Filter) Response(freq float64) float64 {
	resp := C.fid_response(f.design.fidFilter, C.double(freq))
	return float64(resp)
}

func (f *Filter) PhaseResponse(freq float64) (rsp float64, phase float64) {
	var phase_ C.double
	resp_ := C.fid_response_pha(f.design.fidFilter, C.double(freq), &phase_)
	return float64(resp_), float64(phase_)
}

func (f *Filter) Delay() int {
	delay_ := C.fid_calc_delay(f.design.fidFilter)
	return int(delay_)
}

//Reset the filter buffer
func (f *Filter) Zap() {
	C.fid_run_zapbuf(f.buf)
}

//Deallocate the run and buffer objects
//For use when filtering complete
func (f *Filter) Close() {
	C.free(unsafe.Pointer(f.design.fidFilter))
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
