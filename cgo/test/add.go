package test

/*
#cgo CXXFLAGS: -std=c++11
#cgo CFLAGS: -I/usr/include/c++/11
#include "add.h"
*/
import "C"

func Add(a, b int) int {
	return int(C.add(C.int(a), C.int(b)))
}
