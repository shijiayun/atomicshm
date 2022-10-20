// +build linux

package atomicshm

import (
	"os"
	"path/filepath"
)

const shmDir = "/dev/shm"

func shmOpen(name string, flag int, perm os.FileMode) (*os.File, error) {
	name = filepath.Join(shmDir, name)
	return os.OpenFile(name, flag, perm)
}

func shmUnlink(name string) error {
	if !filepath.IsAbs(name) {
		name = filepath.Join(shmDir, name)
	}
	return os.Remove(name)
}
