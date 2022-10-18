package atomicshm

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
	"unsafe"
)

const (
	pageSize = 4096
)

type Uint64 struct {
	ptr *uint64
}

func OpenOrCreateUint64(name string) (*Uint64, error) {
	data, err := openOrCreateShm(name)
	if err != nil {
		return nil, err
	}
	u := &Uint64{
		ptr: (*uint64)(unsafe.Pointer(&data[0])),
	}
	return u, nil
}

func openOrCreateShm(name string) ([]byte, error) {
	file, err := shmOpen(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("open shm file error: %w", err)
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file error: %w", err)
	}
	if fi.Size() != pageSize {
		if err = file.Truncate(pageSize); err != nil {
			return nil, fmt.Errorf("truncate file error: %w", err)
		}
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, int(pageSize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap file error: %w", err)
	}

	if !isAligned(data, pageSize) {
		return nil, errors.New("shm is not aligned")
	}
	return data, nil
}

func isAligned(data []byte, alignSize int) bool {
	return (uintptr(unsafe.Pointer(&data[0])) & uintptr(alignSize-1)) == 0
}

func (u *Uint64) Store(v uint64) {
	atomic.StoreUint64(u.ptr, v)
}

func (u *Uint64) Load() (val uint64) {
	return atomic.LoadUint64(u.ptr)
}

func (u *Uint64) Add(delta uint64) (new uint64) {
	return atomic.AddUint64(u.ptr, delta)
}

func (u *Uint64) CompareAndSwap(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(u.ptr, old, new)
}

func (u *Uint64) Swap(new uint64) (old uint64) {
	return atomic.SwapUint64(u.ptr, new)
}
