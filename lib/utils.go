package main

import "C"
import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"unsafe"
)

func ServerFromPtr(srv uintptr) *server.Server {
	return (*server.Server)(unsafe.Pointer(srv))
}

func PlayerFromPtr(pl uintptr) *player.Player {
	return (*player.Player)(unsafe.Pointer(pl))
}

func GoArrayToCArray[T any](t []*T) uintptr {
	cArray := C.malloc(C.size_t(C.int(len(t))) * C.size_t(unsafe.Sizeof(uintptr(0))))

	a := (*[1<<30 - 1]uintptr)(cArray)
	for index, value := range t {
		a[index] = uintptr(unsafe.Pointer(value))
	}

	return (uintptr)(unsafe.Pointer(cArray))
}

type (
	CArray  = uintptr
	CServer = uintptr
	CPlayer = uintptr
	CString = *C.char
)
