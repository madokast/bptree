# bptree
Golang 实现的 B+树，mmap 友好

## 特点
1. 不依赖于特定的 mmap 库，只要实现一个简单的内存管理器 MemManager 就可以使用。（为什么造轮子理由1）（没有内存释放逻辑，简单避免出错）
2. 可以自定义 key 的排序算法。（为什么造轮子理由2）
3. key 和 value 都支持 null，也就是说插入 \<null, null\> 是有语义的
4. value 大小任意。

## 限制
1. 只能插入/修改，无法删除。（可以插入 value=null 实现删除）
2. key 的大小固定，非空时占用 8 bytes 空间，因此可以直接存储 int64、float64 数据。如果是 string 等变长类型需要先 hash 到 int64，遇到冲突在 value 中链式存储

## 使用方法
使用上比较原始，需要进一步封装

```go
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
```