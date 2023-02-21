package bptree

import (
	"fmt"
	"github.com/madokast/bptree/memory"
	"math/rand"
	"strconv"
	"testing"
	"unsafe"
)

var key = new(int64)
var key2 = new(int64)
var key3 = new(int64)
var key4 = new(int64)

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

func keyString(p uintptr) string {
	return strconv.Itoa(int(*((*int64)(unsafe.Pointer(p)))))
}

func keyFunc(p uintptr) interface{} {
	return *((*int64)(unsafe.Pointer(p)))
}

func TestEmpty(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOneNilKey(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 1
	tree.Insert(0, 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOneNilKey2(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 123
	tree.Insert(0, uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOneNilValue(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 1
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOne(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOne2(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOneTwice(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 15
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestOneTwice2Val(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	*key2 = 120
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestTwo(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestTwoReverse(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 14
	*key2 = 13
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	point := tree.root.items[0].valueLoc
	p := directory.PointerAt(point)
	t.Log(*((*int64)(unsafe.Pointer(p))))
}

func TestThree(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key3)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key3)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestThree2(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key3)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestThree3nil(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(0, uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key3)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestThree4(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	tree.Insert(uintptr(unsafe.Pointer(key3)), uintptr(unsafe.Pointer(key3)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestFour(t *testing.T) {
	directory := memory.New(2048)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	*key4 = 14
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key3)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key3)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key4)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestFourReverse(t *testing.T) {
	directory := memory.New(2048)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	*key4 = 14
	tree.Insert(uintptr(unsafe.Pointer(key4)), uintptr(unsafe.Pointer(key3)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key3)), uintptr(unsafe.Pointer(key)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key2)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
	tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)
	t.Log(tree.PrintTree(keyString, keyString))
}

func Test10(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	for i := 0; i < 10; i++ {
		*key = int64(i)
		tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
		t.Log(tree.PrintTree(keyString, keyString))
	}
}

func Test10_2(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	for i := 0; i < 10; i++ {
		*key = int64(i)
		tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key)), 8)
		t.Log(tree.PrintTree(keyString, keyString))
	}
}

func Test10_null(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	tree.Insert(0, 0, 0)
	for i := 0; i < 10; i++ {
		*key = int64(i)
		tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key)), 8)
		t.Log(tree.PrintTree(keyString, keyString))
	}
}

func Test10_null2(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	for i := 0; i < 10; i++ {
		*key = int64(i)
		tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key)), 8)
		t.Log(tree.PrintTree(keyString, keyString))
	}
	tree.Insert(0, 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestRandom10(t *testing.T) {
	for temp := 0; temp < 1000; temp++ {
		directory := memory.New(1024)
		tree := New(directory, keyComp)
		set := map[int64]struct{}{}
		for i := 0; i < 10; i++ {
			*key = int64(rand.Int31n(100)) - 50
			set[*key] = struct{}{}
			//t.Log("add", *key)
			tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
			//t.Log(tree.PrintTree(keyString, keyString))
		}
		keys := tree.AllKeys(keyFunc)
		if len(keys) != len(set) {
			panic(fmt.Sprintf("%v\n%v", keys, set))
		}
		for _, key := range keys {
			_, ok := set[key.(int64)]
			if !ok {
				panic(fmt.Sprintf("%v\n%v", keys, set))
			}
		}
	}
}

func TestRandom100(t *testing.T) {
	for temp := 0; temp < 200; temp++ {
		directory := memory.New(1024)
		tree := New(directory, keyComp)
		set := map[int64]struct{}{}
		for i := 0; i < 100; i++ {
			*key = int64(rand.Int31n(100)) - 50
			set[*key] = struct{}{}
			//t.Log("add", *key)
			tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
			//t.Log(tree.PrintTree(keyString, keyString))
		}
		keys := tree.AllKeys(keyFunc)
		if len(keys) != len(set) {
			panic(fmt.Sprintf("%v\n%v", keys, set))
		}
		for _, key := range keys {
			_, ok := set[key.(int64)]
			if !ok {
				panic(fmt.Sprintf("%v\n%v", keys, set))
			}
		}
	}
}

func TestRandom1000(t *testing.T) {
	for temp := 0; temp < 200; temp++ {
		directory := memory.New(1024)
		tree := New(directory, keyComp)
		set := map[int64]struct{}{}
		for i := 0; i < 1000; i++ {
			*key = int64(rand.Int31n(100)) - 50
			set[*key] = struct{}{}
			//t.Log("add", *key)
			tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
			//t.Log(tree.PrintTree(keyString, keyString))
		}
		keys := tree.AllKeys(keyFunc)
		if len(keys) != len(set) {
			panic(fmt.Sprintf("%v\n%v", keys, set))
		}
		for _, key := range keys {
			_, ok := set[key.(int64)]
			if !ok {
				panic(fmt.Sprintf("%v\n%v", keys, set))
			}
		}
	}
}

func TestRandom10000(t *testing.T) {
	for temp := 0; temp < 100; temp++ {
		directory := memory.New(1024)
		tree := New(directory, keyComp)
		set := map[int64]struct{}{}
		for i := 0; i < 10000; i++ {
			*key = int64(rand.Uint64())
			set[*key] = struct{}{}
			//t.Log("add", *key)
			tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
			//t.Log(tree.PrintTree(keyString, keyString))
		}
		keys := tree.AllKeys(keyFunc)
		if len(keys) != len(set) {
			panic(fmt.Sprintf("%v\n%v", keys, set))
		}
		for _, key := range keys {
			_, ok := set[key.(int64)]
			if !ok {
				panic(fmt.Sprintf("%v\n%v", keys, set))
			}
		}
	}
}

func Test10_3(t *testing.T) {
	// 改了错误
	//  // newLeaf 被指需要修改
	//	if !newLeaf.isLeaf() {
	//		for i := uint32(0); i < newLeaf.itemNumber; i++ {
	//			child := t.readNode(newLeaf.items[i].valueLoc)
	//			child.fatherPoint = newLeaf.selfPoint
	//		}
	//	}
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 7
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = -29
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = 39
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = 49
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = -50
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = -45
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = 38
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = -12
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = -47
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
	*key = 5
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 0)
	t.Log(tree.PrintTree(keyString, keyString))
}

func TestFindNull(t *testing.T) {
	directory := memory.New(1024)
	tree := New(directory, keyComp)
	*key = 11
	*key2 = 12
	*key3 = 13
	tree.Insert(uintptr(unsafe.Pointer(key)), 0, 8)
	tree.Insert(uintptr(unsafe.Pointer(key2)), 0, 8)
	tree.Insert(uintptr(unsafe.Pointer(key3)), 0, 8)

	exist, _ := tree.Find(uintptr(unsafe.Pointer(key)))
	if !exist {
		panic(exist)
	}
	exist, _ = tree.Find(uintptr(unsafe.Pointer(key2)))
	if !exist {
		panic(exist)
	}
	exist, _ = tree.Find(uintptr(unsafe.Pointer(key3)))
	if !exist {
		panic(exist)
	}
	*key = 10
	*key2 = 14
	*key3 = 9
	exist, _ = tree.Find(uintptr(unsafe.Pointer(key)))
	if exist {
		panic(exist)
	}
	exist, _ = tree.Find(uintptr(unsafe.Pointer(key2)))
	if exist {
		panic(exist)
	}
	exist, _ = tree.Find(uintptr(unsafe.Pointer(key3)))
	if exist {
		panic(exist)
	}
}

func TestFind(t *testing.T) {
	for temp := 0; temp < 200; temp++ {
		directory := memory.New(1024)
		tree := New(directory, keyComp)
		set := map[int64]struct{}{}
		for i := 0; i < 1000; i++ {
			*key = int64(rand.Int31n(100)) - 50
			*key2 = *key + 1000
			set[*key] = struct{}{}
			tree.Insert(uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(key2)), 8)

			exist, value := tree.Find(uintptr(unsafe.Pointer(key)))
			if !exist {
				panic(exist)
			}
			if *key2 != readInt64(value) {
				panic(readInt64(value))
			}

		}
		for k := range set {
			*key = k
			exist, value := tree.Find(uintptr(unsafe.Pointer(key)))
			if !exist {
				panic(exist)
			}
			if k+1000 != readInt64(value) {
				panic(readInt64(value))
			}
		}
	}
}

func readInt64(p uintptr) int64 {
	return *((*int64)(unsafe.Pointer(p)))
}
