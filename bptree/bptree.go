package bptree

import (
	"github.com/madokast/bptree/memory"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

/**
B+树
1. 每个节点最多 degree 个 key。
2. key 的数目和子节点数目相同，即 key 和叶子节点一一对应。key 就是对应的叶子节点中最大的 key
3. key 可以为 null，非空时占用 8 bytes 空间，因此可以直接存储 int64、float64 数据。如果是 string 等变长类型需要先 hash 到 int64，在变长的 value 中链式存储
4. value 可以为 null，空间大小不限
*/

const (
	assert    = true
	printMode = true
	// 度。即节点 item 最大数目，决定节点大小
	degree = 3
	// key 长度，不要改
	keySize = 8
	// node 模式，叶子节点、根节点、中间节点
	modeLeaf = byte(1 << 2)
	modeRoot = byte(1 << 3)
	modeMid  = byte(1 << 4)
	// item 模式
	nullKeyFlag      = byte(1 << 2)
	notNullKeyFlag   = byte(1 << 3)
	nullBlockBidFlag = uint32(0x_FFFF_FFFF)
	nullStr          = "nil"
)

var nodeSz = uint32(unsafe.Sizeof(node{}))
var itemSz = uint32(unsafe.Sizeof(item{}))

// item 保存 key 值和指针信息
type item struct {
	null     byte // 第一个 byte 指定表示 key 为 null 与否
	padding  [7]byte
	key      [keySize]byte   // 可以直接保存 int64 或 float64 值，如果是 string 等变长类型需要先 hash 到 int64，在变长的 value 中链式存储
	valueLoc memory.Location // blockId = nullBlockBidFlag 表示 null
}

type node struct {
	itemNumber  uint32 // items 数目
	mode        byte   // 节点模式，叶子节点、根节点、中间节点
	padding     [3]byte
	items       [degree]item    // item 数据
	selfPoint   memory.Location // node 自己的地址信息
	fatherPoint memory.Location // father 指向父节点。fatherBlockId = nullBlockBidFlag 表示父节点为 null，说明自己就是根
	nextPoint   memory.Location // next 指向下一兄弟节点。nextBlockId = nullBlockBidFlag 表示下一兄弟节点为 null，说明自己就是最右边一个节点
}

type Tree struct {
	root    *node
	dir     memory.MemManager
	compare func(key1 uintptr, key2 *[keySize]byte, key2Null byte) int
}

func New(dir memory.MemManager, compareFunc func(k1, k2 uintptr) int) *Tree {
	return &Tree{
		root: nil,
		dir:  dir,
		compare: func(key1 uintptr, key2 *[keySize]byte, key2Null byte) int {
			// key1 == null 视为最小元素
			if key1 == 0 {
				if key2Null == nullKeyFlag {
					return 0
				} else {
					return -1
				}
			} else {
				if key2Null == nullKeyFlag {
					return 1
				}
				return compareFunc(key1, uintptr(unsafe.Pointer(key2)))
			}
		},
	}
}

// Insert 插入或者 update value
// key = 0 表示 key 为 null。value = 0 表示 value 为 null
func (t *Tree) Insert(key uintptr, value uintptr, valueLength uint32) {
	// 将 value、valueLength 转为定长的 blockId、blockOffset
	if value == 0 {
		t.insert0(key, memory.Location{BlockId: nullBlockBidFlag})
	} else {
		diskPtr, pointer := t.dir.Allocate(valueLength)
		memCopy(value, pointer, valueLength)
		t.insert0(key, diskPtr)
	}
}

// Find 查找 key 对应的 val。返回 exist 是否找到
// 因为可以存 null val，通过 value = 0 标识
// 没有存 value length 所以无法返回，看来写变长数据时需要自解析
func (t *Tree) Find(key uintptr) (exist bool, value uintptr) {
	leaf := t.findLeaf(key, false)
	for i := uint32(0); i < leaf.itemNumber; i++ {
		if t.compare(key, &leaf.items[i].key, leaf.items[i].null) == 0 {
			loc := leaf.items[i].valueLoc
			if loc.BlockId == nullBlockBidFlag {
				return true, 0
			} else {
				return true, t.dir.PointerAt(loc)
			}
		}
	}
	return false, 0
}

// insert0 实际插入逻辑
func (t *Tree) insert0(key uintptr, valLoc memory.Location) {
	if t.root == nil { // 懒初始化
		t.newRoot(key, valLoc)
		t.root.mode |= modeLeaf
	} else {
		leaf := t.findLeaf(key, true)
		ok := t.tryInsertNode(leaf, key, valLoc)
		if !ok { // 没有插入成功，说明满了，需要切开
			_, _ = t.splitAndInsert(leaf, key, valLoc)
		}
	}
}

// tryInsertNode 尝试插入 key 到 node 中，如果存在相同的 key 则更新 value，如果插不进去返回 false
func (t *Tree) tryInsertNode(n *node, key uintptr, valLoc memory.Location) bool {
	// 找到插入点，local 及其后面的都需要移动
	local := uint32(0)
	for local < n.itemNumber {
		if t.compare(key, &n.items[local].key, n.items[local].null) > 0 {
			local++
		} else {
			break
		}
	}

	// 可能 local 就是 key，写入即可
	if local < n.itemNumber && t.compare(key, &n.items[local].key, n.items[local].null) == 0 {
		n.items[local].valueLoc = valLoc
		return true
	}

	// 否则大于 local 的都需要移动，判断能否移动
	if n.itemNumber == degree {
		return false
	}

	// 移动
	if local < n.itemNumber {
		memCopy(uintptr(unsafe.Pointer(&n.items[local])), uintptr(unsafe.Pointer(&n.items[local+1])), (n.itemNumber-local)*itemSz)
	}

	// 写入
	i := item{
		valueLoc: valLoc,
	}
	i.setKey(key)
	n.items[local] = i
	n.itemNumber++
	return true
}

// updateNode 更新 n 中 key 对应的 val
func (t *Tree) updateNode(n *node, key uintptr, valLoc memory.Location) bool {
	for i := uint32(0); i < n.itemNumber; i++ {
		if t.compare(key, &n.items[i].key, n.items[i].null) == 0 {
			n.items[i].valueLoc = valLoc
			return true
		}
	}
	return false
}

// splitAndInsert leaf 中无法直接插入 key，应当先切分再插入，返回切分后新生成的 node，以及实际插入 key 的 node
func (t *Tree) splitAndInsert(leaf *node, key uintptr, valLoc memory.Location) (*node, *node) {
	newLeaf := t.newNode()
	newLeaf.mode = leaf.mode
	// 兄弟指针
	newLeaf.nextPoint = leaf.nextPoint
	leaf.nextPoint = newLeaf.selfPoint
	// 父指针
	newLeaf.fatherPoint = leaf.fatherPoint

	// leaf 的一半移过去
	mid := leaf.itemNumber / 2
	// 移动
	memCopy(uintptr(unsafe.Pointer(&leaf.items[mid])), uintptr(unsafe.Pointer(&newLeaf.items[0])), (leaf.itemNumber-mid)*itemSz)
	// 更新 itemNumber
	newLeaf.itemNumber = leaf.itemNumber - mid
	leaf.itemNumber = mid

	// newLeaf 被指需要修改
	if !newLeaf.isLeaf() {
		for i := uint32(0); i < newLeaf.itemNumber; i++ {
			child := t.readNode(newLeaf.items[i].valueLoc)
			child.fatherPoint = newLeaf.selfPoint
		}
	}

	// 插入新的，问题是插入哪个。如果 key 大于新叶子（右边的）newLeaf 第一个值，就插入新叶子，否则插入旧叶子 leaf
	var insertNode *node
	if t.compare(key, &newLeaf.items[0].key, newLeaf.items[0].null) > 0 {
		insertNode = newLeaf
	} else {
		insertNode = leaf
	}
	// 插入，必定成功
	ok := t.tryInsertNode(insertNode, key, valLoc)
	if !ok {
		panic("splitting cannot insert")
	}

	// 更新父节点
	t.insertFather(leaf, newLeaf)
	return newLeaf, insertNode
}

// insertFather 当节点分裂为 left 和 right 后，需要修改父节点一些信息
func (t *Tree) insertFather(left *node, right *node) {
	if left.isRoot() {
		// left 是根节点，说明没有父亲，自己 new 一个爸爸。把 left.maxKey 插入
		t.newRoot(left.maxKey(), left.selfPoint)
		// 把 right.maxKey 也插入
		ok := t.tryInsertNode(t.root, right.maxKey(), right.selfPoint)
		if !ok {
			panic("root cannot hold two")
		}
		// left 和 right 的模式删除 modeRoot
		left.mode ^= modeRoot
		right.mode ^= modeRoot
		//
		// left 和 right 都指向新爸爸
		left.fatherPoint = t.root.selfPoint
		right.fatherPoint = t.root.selfPoint
	} else {
		// 有父亲，那就读出来
		father := t.readNode(left.fatherPoint)
		// 把 left.maxKey 插到父结点中
		ok := t.tryInsertNode(father, left.maxKey(), left.selfPoint)
		if ok { // 父不分裂，更新父节点中 right.maxKey() 的值，指向现在的 right
			ok2 := t.tryInsertNode(father, right.maxKey(), right.selfPoint)
			if !ok2 {
				panic("updateDuplicateValue fail?")
			}
		} else { // 父需要分裂，拿到分裂后的 anotherFather，还有实际插入 left.maxKey 的 insertKeyFather
			anotherFather, insertKeyFather := t.splitAndInsert(father, left.maxKey(), left.selfPoint)
			// 父发生了分裂，说明父是中间节点
			father.mode, anotherFather.mode = modeMid, modeMid
			// left 指向 insertKeyFather
			left.fatherPoint = insertKeyFather.selfPoint
			// 从 father 和 anotherFather 找到 right 的父节点 rightFather
			rightFather := t.findFather(right, father, anotherFather)
			// right 指向 rightFather
			right.fatherPoint = rightFather.selfPoint
			// 更新 rightFather 的 right.maxKey 新值
			ok2 := t.updateNode(rightFather, right.maxKey(), right.selfPoint)
			if !ok2 {
				panic("replace fail?")
			}
		}
	}
}

func (i *item) setKey(key uintptr) {
	if key == 0 {
		i.null = nullKeyFlag
	} else {
		i.null = notNullKeyFlag
		memCopy(key, uintptr(unsafe.Pointer(&i.key)), keySize)
	}
}

func (t *Tree) PrintTree(keyString func(p uintptr) string, valString func(p uintptr) string) string {
	if keyString == nil {
		keyString = func(p uintptr) string {
			return strconv.Itoa(int(*((*int64)(unsafe.Pointer(p)))))
		}
	}
	if valString == nil {
		valString = keyString
	}

	if t.root == nil {
		return "empty"
	}

	sb := strings.Builder{}

	nodes := []*node{t.root}
	for len(nodes) > 0 {
		cur := nodes[0]
		nodes = nodes[1:] // pop

		sb.WriteString("[")
		if printMode {
			sb.WriteString(cur.modeStr())
		}

		for i := uint32(0); i < cur.itemNumber; i++ {
			it := &cur.items[i]
			if it.isNullKey() {
				sb.WriteString(nullStr)
			} else {
				sb.WriteString(keyString(uintptr(unsafe.Pointer(&it.key))))
			}
			if cur.isLeaf() {
				if it.isNullValue() {
					sb.WriteString(":" + nullStr)
				} else {
					sb.WriteString(":" + valString(t.dir.PointerAt(it.valueLoc)))
				}

			} else {
				nodes = append(nodes, t.readNode(it.valueLoc))
			}
			if i < cur.itemNumber-1 {
				sb.WriteString(",")
			}
		}

		sb.WriteString("]")
		if cur.hasNext() {
			sb.WriteString("->")
		} else {
			sb.WriteString("\n")
		}
	}
	s := sb.String()
	return s[:len(s)-1]
}

func (t *Tree) AllKeys(keyFun func(p uintptr) interface{}) []interface{} {
	keys := make([]interface{}, 0)
	// null 是最小的 key
	leaf := t.findLeaf(0, false)
	for leaf != nil {
		for i := uint32(0); i < leaf.itemNumber; i++ {
			it := &leaf.items[i]
			if it.isNullKey() {
				keys = append(keys, nullStr)
			} else {
				keys = append(keys, keyFun(uintptr(unsafe.Pointer(&it.key))))
			}

		}
		if leaf.hasNext() {
			leaf = t.readNode(leaf.nextPoint)
		} else {
			leaf = nil
		}
	}
	return keys
}

/*========== new =============*/

// newRoot 新建一个 root，并插入 key
func (t *Tree) newRoot(key uintptr, valLoc memory.Location) {
	t.root = t.newNode()
	t.root.mode = modeRoot
	t.root.itemNumber = 1
	// 没有父亲，没有兄弟
	t.root.fatherPoint.BlockId = nullBlockBidFlag
	t.root.nextPoint.BlockId = nullBlockBidFlag

	i := item{
		valueLoc: valLoc,
	}
	i.setKey(key)

	t.root.items[0] = i
}

func (t *Tree) newNode() *node {
	diskPtr, pointer := t.dir.Allocate(nodeSz)
	n := (*node)(unsafe.Pointer(pointer))
	n.selfPoint = diskPtr
	return n
}

/*========== finder =============*/

func (t *Tree) findLeaf(key uintptr, updateMaxKey bool) *node {
	leaf := t.root
	for !leaf.isLeaf() {
		it := uint32(0)
		for it < leaf.itemNumber {
			// key > leaf.items[it].key
			if t.compare(key, &leaf.items[it].key, leaf.items[it].null) > 0 {
				it++
			} else {
				break
			}
		}

		if it == leaf.itemNumber {
			it--
			if updateMaxKey {
				(&leaf.items[it]).setKey(key)
			}
		}

		leaf = t.readNode((&leaf.items[it]).valueLoc)
	}

	return leaf
}

func (t *Tree) findFather(n *node, f1 *node, f2 *node) *node {
	maxKey := n.maxKey()
	for i := uint32(0); i < f1.itemNumber; i++ {
		if t.compare(maxKey, &f1.items[i].key, f1.items[i].null) == 0 {
			return f1
		}
	}
	for i := uint32(0); i < f2.itemNumber; i++ {
		if t.compare(maxKey, &f2.items[i].key, f2.items[i].null) == 0 {
			return f2
		}
	}
	panic("no father in them")
}

/*========== reader =============*/

func (i *item) isNullKey() bool {
	if assert {
		if i.null != nullKeyFlag && i.null != notNullKeyFlag {
			panic(i.null)
		}
	}

	return i.null == nullKeyFlag
}

func (i *item) isNullValue() bool {
	if assert {
		if i.null != nullKeyFlag && i.null != notNullKeyFlag {
			panic(i.null)
		}
	}

	return i.valueLoc.BlockId == nullBlockBidFlag
}

func (t *Tree) readNode(valLoc memory.Location) *node {
	if assert {
		if valLoc.BlockId == nullBlockBidFlag {
			panic("read null")
		}
	}

	p := t.dir.PointerAt(valLoc)
	return (*node)(unsafe.Pointer(p))
}

/*========== node method =============*/

func (n *node) isLeaf() bool {
	if assert {
		allMode := modeLeaf | modeRoot | modeMid
		if (n.mode | allMode) != allMode {
			panic(n.mode)
		}
	}

	return n.mode&modeLeaf == modeLeaf
}

func (n *node) isRoot() bool {
	if assert {
		allMode := modeLeaf | modeRoot | modeMid
		if (n.mode | allMode) != allMode {
			panic(n.mode)
		}
	}

	return n.mode&modeRoot == modeRoot
}

func (n *node) isMid() bool {
	if assert {
		allMode := modeLeaf | modeRoot | modeMid
		if (n.mode | allMode) != allMode {
			panic(n.mode)
		}
	}

	return n.mode&modeMid == modeMid
}

func (n *node) maxKey() uintptr {
	if assert && n.itemNumber == 0 {
		panic("no key")
	}
	// bug 狗屁 go 语言，这里必须取地址
	maxItem := &n.items[n.itemNumber-1]
	if maxItem.isNullKey() {
		return 0
	} else {
		return uintptr(unsafe.Pointer(&(maxItem.key[0])))
	}
}

func (n *node) minKey() uintptr {
	if assert && n.itemNumber == 0 {
		panic("no key")
	}
	maxItem := &n.items[0]
	if maxItem.isNullKey() {
		return 0
	} else {
		return uintptr(unsafe.Pointer(&(maxItem.key[0])))
	}
}

func (n *node) hasNext() bool {
	return n.nextPoint.BlockId != nullBlockBidFlag
}

func (n *node) modeStr() string {
	if assert {
		allMode := modeLeaf | modeRoot | modeMid
		if (n.mode | allMode) != allMode {
			panic(n.mode)
		}
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	if n.isRoot() {
		sb.WriteString("R")
	}
	if n.isMid() {
		sb.WriteString("M")
	}
	if n.isLeaf() {
		sb.WriteString("E")
	}
	sb.WriteString(")")
	return sb.String()
}

/*========== utils =============*/

func memCopy(src, des uintptr, length uint32) {
	Len := int(length)
	if length == 0 {
		return
	}
	srcHelper := *((*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: src, Len: Len, Cap: Len,
	})))
	desHelper := *((*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: des, Len: Len, Cap: Len,
	})))
	copy(desHelper, srcHelper)
}
