package memory

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

/*
内存控制，模拟 mmap
内存由一个一个定长的 block 组成，典型大小为 128M
一组 block 称为 Directory 目录
注意：内存只分配，不释放
*/

type Directory struct {
	blocks    []*block
	blockSize uint32 // 一个 block 的大小，注意不是 len(blocks)
}

type block struct {
	data       []byte  // 物理数据，引用防止 gc，无意义
	header     uintptr // 头指针
	freeOffset uint32  // 未分配位置
	remaining  uint32  // 剩余空间。blockSize = freeOffset + remaining
}

func New(blockSize uint32) *Directory {
	return &Directory{
		blocks:    []*block{newBlock(blockSize)},
		blockSize: blockSize,
	}
}

func (d *Directory) Allocate(size uint32) (ptr Location, pointer uintptr) {
	ptr.BlockId = uint32(len(d.blocks) - 1)
	last := d.blocks[ptr.BlockId]
	if last.remaining >= size {
		ptr.BlockOffset = last.freeOffset
		last.freeOffset += size
		last.remaining -= size
		return ptr, d.PointerAt(ptr)
	} else if size <= d.blockSize {
		d.blocks = append(d.blocks, newBlock(d.blockSize))
		return d.Allocate(size)
	} else {
		panic(strconv.Itoa(int(size)) + " is too large")
	}
}

func (d *Directory) PointerAt(ptr Location) uintptr {
	// 头指针 + 偏移
	return d.blocks[ptr.BlockId].header + uintptr(ptr.BlockOffset)
}

func newBlock(blockSize uint32) *block {
	data := make([]byte, blockSize, blockSize)
	return &block{
		data:       data,
		header:     sliceHeader(data).Data,
		freeOffset: 0,
		remaining:  blockSize,
	}
}

func (d *Directory) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("blockSize=%d, len(blocks)=%d\n", d.blockSize, len(d.blocks)))
	for _, b := range d.blocks {
		sb.WriteString(b.String() + "\n")
	}
	return sb.String()
}

func (b *block) String() string {
	return fmt.Sprintf("header=%d, freeOffset=%d, remaining=%d", b.header, b.freeOffset, b.remaining)
}

func sliceHeader(bytes []byte) reflect.SliceHeader {
	return *((*reflect.SliceHeader)(unsafe.Pointer(&bytes)))
}
