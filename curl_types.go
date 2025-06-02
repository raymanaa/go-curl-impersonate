package curl

import (
	"unsafe"
)

type (
	CurlVersion uint32
	CurlInfo    uint32
	CurlError   uint32
	CurlCode    uint32
	Code        uint32
)

type (
	MultiHandle unsafe.Pointer
	ShareHandle unsafe.Pointer
	CurlMsg     unsafe.Pointer
)

type (
	MultiCode       uint32
	ShareCode       uint32
	MultiOption     uint32
	ShareOption     uint32
	CurlMultiMsgTag uint32
)

type CurlSlist unsafe.Pointer

type CurlHttpFormPost unsafe.Pointer

type CurlVersionInfoDataLayout struct {
	Age           uint32
	Version       uintptr
	VersionNum    uint32
	Host          uintptr
	Features      int32
	SslVersion    uintptr
	SslVersionNum int32
	LibzVersion   uintptr
	Protocols     uintptr

	Ares    uintptr
	AresNum int32

	Libidn uintptr

	IconvVerNum int32

	LibsshVersion uintptr

	BrotliVersionNum uint32
	BrotliVersion    uintptr

	Nghttp2VersionNum uint32
	Nghttp2Version    uintptr

	QuicVersion uintptr

	CAInfo uintptr
	CAPath uintptr

	ZstdVersionNum uint32
	ZstdVersion    uintptr

	HyperVersion uintptr

	GsKitVersion uintptr

	MesalinkVersion uintptr
}

type FdSetPlaceholder unsafe.Pointer

type TimeTPlaceholder int64
