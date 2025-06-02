package curl

// No CGo block here - all curl.h dependencies are moved

import (
	"fmt"
	"unsafe" // Keep for unsafe.Pointer if CURLSH.handle remains so
)

// CurlShareError is now an alias for ShareCode (uint32) from others.go.
// The Error() method allows it to be used as an `error`.
func (e ShareCode) Error() string {
	return CurlShareStrerror(e) // Calls wrapper in others.go
}

// newCurlShareError converts ShareCode to error.
func newCurlShareError(errno ShareCode) error {
	if errno == GetCurlshOk() { // Use getter from others.go
		return nil
	}
	return errno // ShareCode itself implements error via the method above
}

type CURLSH struct {
	handle unsafe.Pointer // Opaque handle, result of CurlShareInit()
}

func ShareInit() *CURLSH {
	p := CurlShareInit() // Call wrapper in others.go
	if p == nil {
		return nil // Should not happen unless out of memory
	}
	return &CURLSH{handle: unsafe.Pointer(p)}
}

func (shcurl *CURLSH) Cleanup() error {
	if shcurl.handle == nil {
		return nil // Already cleaned up or never initialized
	}
	err := newCurlShareError(CurlShareCleanup(ShareHandle(shcurl.handle))) // Call wrapper
	shcurl.handle = nil                                                    // Mark as cleaned up
	return err
}

func (shcurl *CURLSH) Setopt(opt int, param any) error {
	if shcurl.handle == nil {
		return fmt.Errorf("curl: share handle is nil")
	}

	option := uint32(opt) // Convert to uint32 for wrapper

	// Share options (SHOPT_) are specific and don't have a broad type system like easy options.
	// You'll typically switch on the `opt` value directly using your defined constants (SHOPT_SHARE, etc.).
	if param == nil {
		// Check if the option is known to accept a nil pointer (e.g., for unsetting a callback)
		// For SHOPT_LOCKFUNC, SHOPT_UNLOCKFUNC, SHOPT_USERDATA, nil might be valid to unset.
		// For SHOPT_SHARE, SHOPT_UNSHARE, param is expected to be an int (curl_lock_data enum).
		// This path needs to be specific to the option.
		// Example:
		// if opt == SHOPT_USERDATA || opt == SHOPT_LOCKFUNC || opt == SHOPT_UNLOCKFUNC {
		return newCurlShareError(CurlShareSetoptPointer(ShareHandle(shcurl.handle), ShareOption(option), nil))
		// }
		// return fmt.Errorf("curl: nil parameter for share option %d which may not support it", opt)
	}

	// Use your defined constants for SHOPT_SHARE, SHOPT_UNSHARE
	// (These should be in your const.go or const_gen.go)
	switch opt {
	case SHOPT_SHARE, SHOPT_UNSHARE: // Assuming SHOPT_SHARE, SHOPT_UNSHARE are int constants
		if val, ok := param.(int); ok {
			return newCurlShareError(CurlShareSetoptLong(ShareHandle(shcurl.handle), ShareOption(option), int64(val)))
		} else if val, ok := param.(int32); ok {
			return newCurlShareError(CurlShareSetoptLong(ShareHandle(shcurl.handle), ShareOption(option), int64(val)))
		} else if val, ok := param.(int64); ok {
			return newCurlShareError(CurlShareSetoptLong(ShareHandle(shcurl.handle), ShareOption(option), val))
		} else {
			return fmt.Errorf("curl: SHOPT_SHARE/UNSHARE expects int parameter, got %T", param)
		}
	// case SHOPT_LOCKFUNC, SHOPT_UNLOCKFUNC:
	//   // These expect function pointers.
	//   // You'd need to get C function pointers for your Go callback wrappers.
	//   // Example: cCallbackPtr := GetShareLockCallbackFuncptr()
	//   // return newCurlShareError(CurlShareSetoptPointer(shcurl.handle, option, cCallbackPtr))
	//   return fmt.Errorf("curl: SHOPT_LOCKFUNC/UNLOCKFUNC not fully implemented in Setopt")
	// case SHOPT_USERDATA:
	//   // This expects a void* (unsafe.Pointer in Go).
	//   // Convert param to unsafe.Pointer. How depends on what `param` is.
	//   // If param is already a pointer: unsafe.Pointer(param_as_ptr_type)
	//   // If param is a Go interface holding a pointer: unsafe.Pointer(reflect.ValueOf(param).Pointer())
	//   // Or if it's intended to be a pointer *to* the param itself: unsafe.Pointer(¶m) (use with caution)
	//   if ptr, ok := param.(unsafe.Pointer); ok {
	//      return newCurlShareError(CurlShareSetoptPointer(shcurl.handle, option, ptr))
	//   }
	//   return fmt.Errorf("curl: SHOPT_USERDATA expects unsafe.Pointer or compatible, got %T", param)
	default:
		return fmt.Errorf("curl: unsupported share option %d or param type %T", opt, param)
	}
}
