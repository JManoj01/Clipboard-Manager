//go:build windows
package clipboard

import (
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	openClipboard    = user32.NewProc("OpenClipboard")
	closeClipboard   = user32.NewProc("CloseClipboard")
	getClipboardData = user32.NewProc("GetClipboardData")
	globalLock       = kernel32.NewProc("GlobalLock")
	globalUnlock     = kernel32.NewProc("GlobalUnlock")
)

const CF_UNICODETEXT = 13

func Read() (string, error) {
	r, _, _ := openClipboard.Call(0)
	if r == 0 {
		return "", nil
	}
	defer closeClipboard.Call()

	h, _, _ := getClipboardData.Call(CF_UNICODETEXT)
	if h == 0 {
		return "", nil
	}

	l, _, _ := globalLock.Call(h)
	if l == 0 {
		return "", nil
	}
	defer globalUnlock.Call(h)

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:])
	return text, nil
}
