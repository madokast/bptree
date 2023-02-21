package memory

/*
bptree 工作需要用到内存管理工具，便于 mmap
*/

// Location 内存定位器，主要用于 mmap 后文件和内存的映射关系
type Location struct {
	BlockId     uint32 // 可以看作文件编号
	BlockOffset uint32 // 可以看作文件内偏移。blockId 和 blockOffset 一起完成寻址
}

// MemManager 内存管理器，一般用于管理 mmap 的内存
// 注意：内存只分配，不释放（简化简化）
type MemManager interface {
	// Allocate 向内存/文件系统请求 size 大小的内存，返回内存定位器 loc 和请求到的内存 pointer
	Allocate(size uint32) (loc Location, pointer uintptr)
	// PointerAt 由内存定位器获取指针
	PointerAt(loc Location) (pointer uintptr)
}
