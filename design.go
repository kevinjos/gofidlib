/*  gofidlib provides go bindings to Jim Peters' fidlib library.

    Copyright (C) 2015  Kevin Schiesser

    fidlib was originally designed as a backend for the 'Fiview'
    application, and to provide performance filtering services to
    EEG applications, such as those in the OpenEEG project:

    http://uazu.net/fiview/
    http://openeeg.sf.net/

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package gofidlib

/*
#cgo CFLAGS: -DT_LINUX -O6 -DRF_CMDLIST
#cgo LDFLAGS: -lm

#include <stdio.h>
#include <stdlib.h>
#include "fidlib.h"

FidRun*
fid_run_new_(FidFilter* filt, FidFunc** funcp)
{
	FidRun *run;

	run = fid_run_new(filt, funcp);

	return run;
}


fid_run(FidFunc *funcp, void *fbuf1, double* in, int size_in, double* out)
{
	int i;

	for (i = 0; i < size_in; i++)
	{
		out[i] = funcp(fbuf1, in[i]);
	}
}
*/
import "C"

import (
	"log"
	"unsafe"
)

type Filter struct {
	design *FilterDesign
	buf    unsafe.Pointer
}

func (f *Filter) Run(in []float64, out []float64) {
	C.fid_run(f.design.Funcp, f.buf, (*C.double)(&in[0]), C.int(len(in)), (*C.double)(&out[0]))
}

func NewFilter(design *FilterDesign) *Filter {
	buf := C.fid_run_newbuf(design.FidRun)
	return &Filter{
		design: design,
		buf:    buf,
	}
}

type FilterDesign struct {
	FidFilter *C.FidFilter
	FidRun    unsafe.Pointer
	Funcp     *C.FidFunc
}

func NewFilterDesign(spec string, rate float64) *FilterDesign {
	var (
		rate_  = C.double(rate)
		spec_  = C.CString(spec)
		filt_  = new(C.FidFilter)
		funcp_ = new(C.FidFunc)
	)
	defer C.free(unsafe.Pointer(spec_))
	err := C.fid_parse(rate_, &spec_, &filt_)
	if err != nil {
		log.Fatal("Error in filter design:", C.GoString(err))
	}
	run := C.fid_run_new_(filt_, &funcp_)
	return &FilterDesign{
		FidFilter: filt_,
		FidRun:    run,
		Funcp:     funcp_,
	}
}

func Version() string {
	return C.GoString(C.fid_version())
}
