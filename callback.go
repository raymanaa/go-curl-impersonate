package curl

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"unsafe"
)

//export goCallHeaderFunction
func goCallHeaderFunction(ptr *C.char, size C.size_t, ctx unsafe.Pointer) C.size_t {
	curl := context_map.Get(uintptr(ctx))
	if curl == nil || curl.headerFunction == nil {
		return 0
	}
	buf := C.GoBytes(unsafe.Pointer(ptr), C.int(size*1))

	if (*curl.headerFunction)(buf, curl.headerData) {
		return C.size_t(size * 1)
	}
	return GetCurlWritefuncPause()
}

//export goCallWriteFunction
func goCallWriteFunction(ptr *C.char, size C.size_t, ctx unsafe.Pointer) C.size_t {
	curl := context_map.Get(uintptr(ctx))
	if curl == nil || curl.writeFunction == nil {
		return 0
	}
	buf := C.GoBytes(unsafe.Pointer(ptr), C.int(size*1))

	if (*curl.writeFunction)(buf, curl.writeData) {
		return C.size_t(size * 1) // Return total bytes processed
	}
	return GetCurlWritefuncPause()
}

//export goCallProgressFunction
func goCallProgressFunction(dltotalC C.double, dlnowC C.double, ultotalC C.double, ulnowC C.double, ctx unsafe.Pointer) C.int {
	curl := context_map.Get(uintptr(ctx))
	if curl == nil || curl.progressFunction == nil {
		return 0
	}
	if (*curl.progressFunction)(float64(dltotalC), float64(dlnowC),
		float64(ultotalC), float64(ulnowC),
		curl.progressData) {
		return 0
	}
	return 1
}

//export goCallReadFunction
func goCallReadFunction(ptr *C.char, size C.size_t, numItems C.size_t, ctx unsafe.Pointer) C.size_t {
	curl := context_map.Get(uintptr(ctx))
	if curl == nil || curl.readFunction == nil {
		return GetCurlReadfuncAbort()
	}

	maxLenToRead := int(size * numItems)
	if maxLenToRead == 0 {
		return 0
	}

	goSliceForReading := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), maxLenToRead)
	bytesWrittenByGoFunc := (*curl.readFunction)(goSliceForReading, curl.readData)

	if bytesWrittenByGoFunc < 0 {
		return GetCurlReadfuncAbort()
	}

	if bytesWrittenByGoFunc > maxLenToRead {
		return GetCurlReadfuncAbort()
	}

	return C.size_t(bytesWrittenByGoFunc)
}
