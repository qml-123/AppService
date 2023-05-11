package ffmpeg

/*
   #cgo pkg-config: libavcodec libavformat libavutil
   #include <stdlib.h>
   #include "ffmpeg_wrapper.h"
*/
import "C"
import (
	"errors"
	"log"
	"unsafe"
)

// 初始化ffmpeg
func init() {
	C.init_ffmpeg()
}

// AV1编码函数
func EncodeAV1(input []byte) ([]byte, error) {
	var output *C.uint8_t
	var outputSize C.int

	ret := C.encode_av1((*C.uint8_t)(&input[0]), C.int(len(input)), &output, &outputSize)
	if ret != 0 {
		return nil, errors.New("failed to encode AV1")
	}

	// 注意: 你需要释放C分配的内存
	defer C.free(unsafe.Pointer(output))

	return C.GoBytes(unsafe.Pointer(output), outputSize), nil
}

// AV1解码函数
func DecodeAV1(input []byte) ([]byte, error) {
	var output *C.uint8_t
	var outputSize C.int

	ret := C.decode_av1((*C.uint8_t)(&input[0]), C.int(len(input)), &output, &outputSize)
	if ret != 0 {
		return nil, errors.New("failed to decode AV1")
	}

	// 注意: 你需要释放C分配的内存
	defer C.free(unsafe.Pointer(output))

	return C.GoBytes(unsafe.Pointer(output), outputSize), nil
}

func test() {
	// 测试
	input := []byte("test data")

	encoded, err := EncodeAV1(input)
	if err != nil {
		panic(err)
	}

	// encoded现在是AV1编码后的视频数据

	decoded, err := DecodeAV1(encoded)
	if err != nil {
		panic(err)
	}
	log.Printf("%v", decoded)
	// decoded现在是解码后的YUV420P格式的视频数据
}
