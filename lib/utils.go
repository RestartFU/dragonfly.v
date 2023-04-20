package main

import "C"
import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"unsafe"
)

func ServerFromPtr(srv CServer) *server.Server {
	return (*server.Server)(srv.Ptr())
}

func PlayerFromPtr(pl CPlayer) *player.Player {
	return (*player.Player)(pl.Ptr())
}

func GoArrayToCArray[T any](t []*T) CArray {
	cArray := C.malloc(C.size_t(C.int(len(t))) * C.size_t(unsafe.Sizeof(uintptr(0))))

	a := (*[1<<30 - 1]uintptr)(cArray)
	for index, value := range t {
		a[index] = uintptr(unsafe.Pointer(value))
	}

	return (CArray)(unsafe.Pointer(cArray))
}

type (
	UINTPTR uintptr

	CArray  = UINTPTR
	CServer = UINTPTR
	CPlayer = UINTPTR
	CString = *C.char
)

func (u UINTPTR) Ptr() unsafe.Pointer {
	return unsafe.Pointer(u)
}
