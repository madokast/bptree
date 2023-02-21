package main

import (
	"fmt"
	"github.com/madokast/bptree/bptree"
	"github.com/madokast/bptree/memory"
	"reflect"
	"unsafe"
)

func main() {
	mem := memory.New(1024)
	tree := bptree.New(mem, keyComp)

	// 插入 123 -> 321
	{
		key := new(int64)
		val := new(int64)
		*key = 123
		*val = 321

		tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(val)), uint32(unsafe.Sizeof(*val)))
	}

	// 插入 3.14 -> "hello, world"
	{
		key := new(float64)
		val := []byte("hello, world")
		*key = 3.14

		valPointer := ((*reflect.SliceHeader)(unsafe.Pointer(&val))).Data

		tree.Insert(uintptr(unsafe.Pointer(key)), valPointer, uint32(len(val)))
	}

	// 查找 123
	{
		key := new(int64)
		*key = 123
		exist, valuePointer := tree.Find(uintptr(unsafe.Pointer(key)))
		if exist {
			value := *((*int64)(unsafe.Pointer(valuePointer)))
			fmt.Println(*key, "->", value)
		} else {
			fmt.Println("不存在", *key)
		}
	}

	// 查找 3.14，需要知道 value 的长度 12
	{
		key := new(float64)
		*key = 3.14
		exist, valuePointer := tree.Find(uintptr(unsafe.Pointer(key)))
		if exist {
			value := *((*string)(unsafe.Pointer(&reflect.StringHeader{
				Data: valuePointer,
				Len:  12,
			})))
			fmt.Println(*key, "->", value)
		} else {
			fmt.Println("不存在", *key)
		}
	}
}

func keyComp(k1, k2 uintptr) int {
	n1 := *((*int64)(unsafe.Pointer(k1)))
	n2 := *((*int64)(unsafe.Pointer(k2)))
	if n1 > n2 {
		return 1
	} else if n1 < n2 {
		return -1
	} else {
		return 0
	}
}
