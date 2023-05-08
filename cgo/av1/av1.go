package av1

/*
#cgo CXXFLAGS: -std=c++11
#cgo CFLAGS: -I/root/cpp/cpp-include/FFmpeg-PlusPlus/libavcodec
#cgo CFLAGS: -I/root/cpp/cpp-include/FFmpeg-PlusPlus/libavutil
#cgo LDFLAGS: -L/root/cpp/cpp-include/FFmpeg-PlusPlus/libavcodec -L/root/cpp/cpp-include/FFmpeg-PlusPlus/libavutil -lavcodec -lavutil -lm -lavformat -lswresample -lz -llzma
#include <stdlib.h>
#include "av1_wrapper.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func AV1Encode(input []byte) ([]byte, error) {
	var output *C.uint8_t
	var outputSize C.size_t

	ret := C.av1_encode((*C.uint8_t)(unsafe.Pointer(&input[0])), C.size_t(len(input)), &output, &outputSize)
	if ret != 0 {
		return nil, fmt.Errorf("AV1 encoding error: %d", ret)
	}

	encoded := C.GoBytes(unsafe.Pointer(output), C.int(outputSize))
	C.free(unsafe.Pointer(output))

	return encoded, nil
}

func AV1Decode(input []byte) ([]byte, error) {
	var output *C.uint8_t
	var outputSize C.size_t

	ret := C.av1_decode((*C.uint8_t)(unsafe.Pointer(&input[0])), C.size_t(len(input)), &output, &outputSize)
	if ret != 0 {
		return nil, fmt.Errorf("AV1 decoding error: %d", ret)
	}

	decoded := C.GoBytes(unsafe.Pointer(output), C.int(outputSize))
	C.free(unsafe.Pointer(output))

	return decoded, nil
}
