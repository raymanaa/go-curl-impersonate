// multi.go

package curl

import "C" // Keep if using C.int, C.long for variables passed to wrappers.
import (
	"fmt"
	// "syscall" // No longer needed for FdSet
	"unsafe"
)

// CurlMultiError Error method and newCurlMultiError remain the same
func (e MultiCode) Error() string {
	return CurlMultiStrerror(e)
}
func newCurlMultiError(errno MultiCode) error {
	if errno == GetCurlmOk() {
		return nil
	}
	return errno
}

// Add String method for CurlMultiMsgTag for better logging
func (t CurlMultiMsgTag) String() string {
	if t == GetCurlmsgDone() { // Assumes GetCurlmsgDone() returns the correct underlying value
		return "CURLMSG_DONE"
	}
	// You can add more cases here if libcurl exposes other CURLMSG types you use
	// e.g. CURLMSG_NONE if that's a possible value from CurlMsgGetMsg
	return fmt.Sprintf("CURLMSG type %d", t)
}

type CURLMessage struct {
	Msg         CurlMultiMsgTag
	Easy_handle *CURL
	DoneResult  Code
	PointerVal  unsafe.Pointer
}

// newCURLMessage remains the same
func newCURLMessage(opaqueCM CurlMsg) *CURLMessage {
	if opaqueCM == nil {
		return nil
	}

	goMsg := new(CURLMessage)
	goMsg.Msg = CurlMsgGetMsg(opaqueCM) // Use accessor

	easyHandlePtr := CurlMsgGetEasyHandle(opaqueCM) // Use accessor (returns unsafe.Pointer)
	if easyHandlePtr != nil {
		// Assuming context_map is properly defined and accessible for mapping C handles to Go *CURL structs
		// goEasyHandle := context_map.Get(uintptr(easyHandlePtr))
		// if goEasyHandle != nil {
		// 	goMsg.Easy_handle = goEasyHandle
		// } else {
		// For simplicity if context_map isn't shown or fully implemented for multi:
		goMsg.Easy_handle = &CURL{handle: easyHandlePtr}
		// }
	}

	if goMsg.Msg == GetCurlmsgDone() {
		goMsg.DoneResult = CurlMsgGetResult(opaqueCM) // Use accessor
	} else {
		goMsg.PointerVal = CurlMsgGetWhatever(opaqueCM) // Use accessor
	}
	return goMsg
}

type CURLM struct {
	handle unsafe.Pointer
}

// MultiInit, Cleanup, Perform, AddHandle, RemoveHandle, Timeout, Setopt
// remain as in the previous *correct* version where they call wrappers
// from others.go and use newCurlMultiError.

func MultiInit() *CURLM {
	p := CurlMultiInit()
	if p == nil {
		return nil
	}
	return &CURLM{handle: unsafe.Pointer(p)}
}

func (mcurl *CURLM) Cleanup() error {
	if mcurl.handle == nil {
		return nil
	}
	err := newCurlMultiError(CurlMultiCleanup(MultiHandle(mcurl.handle)))
	mcurl.handle = nil
	return err
}

func (mcurl *CURLM) Perform() (int, error) {
	if mcurl.handle == nil {
		return 0, fmt.Errorf("curl: multi handle is nil")
	}
	var runningHandles C.int = -1 // C.int might be int32
	err := newCurlMultiError(CurlMultiPerform(MultiHandle(mcurl.handle), unsafe.Pointer(&runningHandles)))
	return int(runningHandles), err
}

func (mcurl *CURLM) AddHandle(easy *CURL) error {
	if mcurl.handle == nil {
		return fmt.Errorf("curl: multi handle is nil")
	}
	if easy == nil || easy.handle == nil {
		return fmt.Errorf("curl: easy handle is nil")
	}
	return newCurlMultiError(CurlMultiAddHandle(MultiHandle(mcurl.handle), easy.handle))
}

func (mcurl *CURLM) RemoveHandle(easy *CURL) error {
	if mcurl.handle == nil {
		return fmt.Errorf("curl: multi handle is nil")
	}
	if easy == nil || easy.handle == nil {
		return fmt.Errorf("curl: easy handle is nil to remove")
	}
	return newCurlMultiError(CurlMultiRemoveHandle(MultiHandle(mcurl.handle), easy.handle))
}

func (mcurl *CURLM) Timeout() (int, error) {
	if mcurl.handle == nil {
		return -1, fmt.Errorf("curl: multi handle is nil")
	}
	var timeoutMs C.long = -1 // C.long can be int32 or int64 depending on platform
	err := newCurlMultiError(CurlMultiTimeout(MultiHandle(mcurl.handle), unsafe.Pointer(&timeoutMs)))
	return int(timeoutMs), err
}

func (mcurl *CURLM) Setopt(opt int, param any) error {
	if mcurl.handle == nil {
		return fmt.Errorf("curl: multi handle is nil")
	}
	option := uint32(opt)
	if param == nil {
		return newCurlMultiError(CurlMultiSetoptPointer(MultiHandle(mcurl.handle), MultiOption(option), nil))
	}
	switch v := param.(type) {
	case int:
		return newCurlMultiError(CurlMultiSetoptLong(MultiHandle(mcurl.handle), MultiOption(option), int64(v)))
	case int32:
		return newCurlMultiError(CurlMultiSetoptLong(MultiHandle(mcurl.handle), MultiOption(option), int64(v)))
	case int64:
		return newCurlMultiError(CurlMultiSetoptLong(MultiHandle(mcurl.handle), MultiOption(option), v))
	case bool:
		var val int64 = 0
		if v {
			val = 1
		}
		return newCurlMultiError(CurlMultiSetoptLong(MultiHandle(mcurl.handle), MultiOption(option), val))
	default:
		if p, ok := param.(unsafe.Pointer); ok {
			return newCurlMultiError(CurlMultiSetoptPointer(MultiHandle(mcurl.handle), MultiOption(option), p))
		}
		return fmt.Errorf("curl: unsupported Setopt param type %T for multi option %d", param, opt)
	}
}

// REMOVE Fdset method
// func (mcurl *CURLM) Fdset(rset, wset, eset *syscall.FdSet) (int, error) {
// 	if mcurl.handle == nil {
// 		return -1, fmt.Errorf("curl: multi handle is nil")
// 	}
// 	var maxFd C.int = -1
// 	err := newCurlMultiError(CurlMultiFdset(MultiHandle(mcurl.handle),
// 		FdSetPlaceholder(rset), FdSetPlaceholder(wset), FdSetPlaceholder(eset), &maxFd))
// 	return int(maxFd), err
// }

// Wait calls curl_multi_wait.
// For simplicity, extraFds (should be *C.struct_curl_waitfd) and extraNumFds are currently not used from Go, pass nil and 0.
// numFdsReady must be a pointer to C.int, and will be populated with the number of file descriptors with activity.
func (mcurl *CURLM) Wait(extraFds unsafe.Pointer, extraNumFds int, timeoutMs int, numFdsReady *C.int) error {
	if mcurl.handle == nil {
		return fmt.Errorf("curl: multi handle is nil")
	}
	// Pass numFdsReady directly as it's already a pointer to C.int
	return newCurlMultiError(CurlMultiWait(MultiHandle(mcurl.handle), extraFds, extraNumFds, timeoutMs, unsafe.Pointer(numFdsReady)))
}

// Info_read uses the wrapper that returns CurlMsg (unsafe.Pointer)
func (mcurl *CURLM) Info_read() (*CURLMessage, int) {
	if mcurl.handle == nil {
		return nil, 0
	}
	var msgsInQueue C.int = 0
	opaqueCM := CurlMultiInfoRead(MultiHandle(mcurl.handle), unsafe.Pointer(&msgsInQueue))
	return newCURLMessage(opaqueCM), int(msgsInQueue)
}
