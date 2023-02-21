package memory

import (
	"testing"
	"unsafe"
)

func TestNew(t *testing.T) {
	directory := New(1024)
	t.Log(directory)
	_, p := directory.Allocate(100)
	*((*byte)(unsafe.Pointer(p))) = 1
	*((*byte)(unsafe.Pointer(p + 1))) = 2
	*((*byte)(unsafe.Pointer(p + 2))) = 3
	t.Log(directory.blocks[0].data[:16])
	t.Log(directory)
	_, _ = directory.Allocate(200)
	t.Log(directory)
}
