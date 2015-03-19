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

FidFilter*
fid_cat_(int freeme, FidFilter* filt1, FidFilter* filt2)
{
	FidFilter* filt;

	filt = fid_cat(freeme, filt1, filt2, 0);

	return filt;
}
*/
import "C"

import (
	"unsafe"
	"errors"
)

type FilterDesign struct {
	fidFilter *C.FidFilter
	fidRun    unsafe.Pointer
	funcp     *C.FidFunc
}

func NewFilterDesign(spec string, rate float64) (*FilterDesign, error) {
	var (
		rate_  = C.double(rate)
		spec_  = C.CString(spec)
		filt_  = new(C.FidFilter)
		funcp_ = new(C.FidFunc)
	)
	defer C.free(unsafe.Pointer(spec_))
	err := C.fid_parse(rate_, &spec_, &filt_)
	if err != nil {
		return nil, errors.New(C.GoString(err))
	}
	run_ := C.fid_run_new_(filt_, &funcp_)
	return &FilterDesign{
		fidFilter: filt_,
		fidRun:    run_,
		funcp:     funcp_,
	}, nil
}

func (fd *FilterDesign) Free() {
	C.free(unsafe.Pointer(fd.fidFilter))
	C.fid_run_free(fd.fidRun)
}

func FidCat(freeme int, filters []*FilterDesign) (*FilterDesign, error) {
	var filt_ *C.FidFilter

	cat := func(sumFilt *C.FidFilter, nextFilt *C.FidFilter) (*C.FidFilter, error) {
		filt_, err := C.fid_cat_(C.int(freeme), sumFilt, nextFilt)
		if err != nil {
			return nil, err
		}
		return filt_, nil
	}

	numFilts := len(filters)

	switch {
	case numFilts < 2:
		return nil, errors.New("Too few filters in argument slice. Must have atleast 2.")
	case numFilts >= 2:
		var err error
		if filt_, err = cat(filters[0].fidFilter, filters[1].fidFilter); err != nil {
			return nil, err
		}
		for i := 2; i < len(filters); i++ {
			if filt_, err = cat(filt_, filters[i].fidFilter); err != nil {
				return nil, err
			}
		}
		if freeme == 1 {
			for _, val := range filters {
				C.fid_run_free(val.fidRun)
			}
		}
	}

	funcp_ := new(C.FidFunc)
	run_ := C.fid_run_new_(filt_, &funcp_)

	return &FilterDesign{
		fidFilter: filt_,
		fidRun:    run_,
		funcp:     funcp_,
	}, nil
}
