//go:build darwin || freebsd

package atomicshm

import (
	"os"
	"syscall"
	"unsafe"
)

func shmOpen(name string, flag int, perm os.FileMode) (*os.File, error) {
	pname, err := syscall.BytePtrFromString(name)
	if err != nil {
		return nil, err
	}
	fd, _, errno := syscall.Syscall(syscall.SYS_SHM_OPEN, uintptr(unsafe.Pointer(pname)), uintptr(flag), uintptr(perm))
	if errno != 0 {
		return nil, errno
	}
	return os.NewFile(fd, name), nil
}

func shmUnlink(name string) error {
	pname, err := syscall.BytePtrFromString(name)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall(syscall.SYS_SHM_UNLINK, uintptr(unsafe.Pointer(pname)), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
