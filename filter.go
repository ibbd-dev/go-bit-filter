package bitFilter

import ()

const (
	// 2MB的存储空间. 1MB=2^20, 一个uint64长度为8个字节，就是2^3
	// 总共的位数：8*1024*1024*2 = 16 777 216
	Size2MB uint64 = 1 << (21 - 3)

	// 更大的空间
	Size8MB   uint64 = 1 << (23 - 3)
	Size32MB  uint64 = 1 << (25 - 3)
	Size128MB uint64 = 1 << (27 - 3)
	Size512MB uint64 = 1 << (29 - 3)
	Size2GB   uint64 = 1 << (31 - 3)

	// 8*1024*2 = 16 384
	Size2KB uint64 = 1 << (11 - 3)
)

type TBits struct {
	buckets []uint64 // 数组

	// 冲突的key值（理论上冲突的key值应该很小才有效）
	// 下标是key，值是hashKey
	collision map[uint64]uint32
}

var (
	// 位的掩码
	bitMask [64]uint64
)

func init() {
	var i uint
	for i = 0; i < 64; i++ {
		bitMask[i] = 1 << i
	}
}

func New(len uint64) *TBits {
	return &TBits{
		buckets:   make([]uint64, len),
		collision: make(map[uint64]uint32),
	}
}

// 复制一个结构
func (b *TBits) Copy() *TBits {
	new := &TBits{
		buckets:   b.buckets[:],
		collision: make(map[uint64]uint32),
	}
	for key, hashKey := range b.collision {
		new.collision[key] = hashKey
	}

	return new
}

// 初始化所有hash key值
func (b *TBits) InitAllHashKeys(hashKeys []uint32) {
	var bucket uint32
	for _, key := range hashKeys {
		bucket = key >> 6 // 确定在那个位上
		//println(bucket, key&7)
		b.buckets[bucket] |= bitMask[key&63] // 63 = 1<<6 - 1
	}
}

// 初始化冲突的key值
func (b *TBits) InitCollisionKeys(keys []uint64, hashKeys []uint32) {
	for index, key := range keys {
		b.collision[key] = hashKeys[index]
	}
}

// Filter 过滤器，如果元素存在，则返回true，否则返回
func (b *TBits) Filter(key uint64, hashKey uint32) (isExist bool) {
	if _, isExist = b.collision[key]; isExist {
		return isExist
	}

	bucket := hashKey >> 6
	if b.buckets[bucket]&bitMask[hashKey&63] == 0 {
		return false
	}

	return true
}

// Add 新增key
// 如果如果hashKey已经存在，但是又没有冲突的话，则返回false，这时需要先添加原来的冲突
func (b *TBits) Add(key uint64, hashKey uint32) (ok bool) {
	var (
		bucket  uint32 = hashKey >> 6        // 确定在那个位上
		keyMask        = bitMask[hashKey&63] // 63 = 1<<6 - 1
	)
	if b.buckets[bucket]&keyMask == keyMask {
		// 判断该haskkey是否已经冲突
		if isExist := b.hasHashKeyInCollision(hashKey); !isExist {
			// 需要添加原来的key到冲突key map才可以继续
			return false
		}

		// 原来已经冲突了，这是直接增加一个冲突即可
		b.collision[key] = hashKey
	} else {
		// 该key完全是新的
		b.buckets[bucket] |= keyMask
	}

	return true
}

// Remove 减少key
func (b *TBits) Remove(key uint64, hashKey uint32) {
	var oldKeys []uint64
	if _, isExist := b.collision[key]; isExist {
		for k, hk := range b.collision {
			if hk == hashKey {
				oldKeys = append(oldKeys, k)
			}
		}

		if len(oldKeys) > 2 {
			// 有三个或者三个以上冲突
			delete(b.collision, key)
		} else {
			// 只有两个冲突，则全部删除
			// 冲突的key，至少也会有两个
			for _, k := range oldKeys {
				delete(b.collision, k)
			}
		}
	} else {
		var (
			bucket  uint32 = hashKey >> 6        // 确定在那个位上
			keyMask        = bitMask[hashKey&63] // 63 = 1<<6 - 1
		)
		b.buckets[bucket] &= ^keyMask
	}
}

// AddCollisionKey 增加冲突key
func (b *TBits) AddCollisionKey(key uint64, hashKey uint32) {
	b.collision[key] = hashKey
}

func (b *TBits) hasHashKeyInCollision(hashKey uint32) (isExist bool) {
	for _, hk := range b.collision {
		if hk == hashKey {
			return true
		}
	}

	return isExist
}
