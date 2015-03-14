package gofidlib

/*
#include <stdio.h>
#include "fidlib.h"
*/
import "C"

func Version() string {
	return C.GoString(C.fid_version())
}
