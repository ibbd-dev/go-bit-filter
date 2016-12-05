package bitFilter

import ()

type TBits struct {
	bytes     []byte          // 字符数组
	bits      int             // 位空间，1个字节8个位
	collision map[uint32]bool // 冲突的key值（理论上冲突的key值应该很小才有效）
}

var (
	// 位的掩码
	bitMask = [8]uint8{1, 2, 4, 8, 16, 32, 64, 128}
)

func New(spaceBytes int) *TBits {
	return &TBits{
		bytes:     make([]byte, spaceBytes),
		bits:      spaceBytes << 3,
		collision: make(map[uint32]bool),
	}
}

func (b *TBits) AddAllHashKeys(hashKeys []uint32) {
	var bucket uint32
	for _, key := range hashKeys {
		bucket = key >> 3 // 确定在那个位上
		//println(bucket, key&7)
		b.bytes[bucket] |= bitMask[key&7]
	}
}

func (b *TBits) AddCollisionKeys(keys []uint32) {
	for _, key := range keys {
		b.collision[key] = true
	}
}

// Filter 过滤器，如果元素存在，则返回true，否则返回
func (b *TBits) Filter(key, hashKey uint32) (isExist bool) {
	if _, isExist = b.collision[key]; isExist {
		return isExist
	}

	bucket := hashKey >> 3
	if b.bytes[bucket]&bitMask[hashKey&7] == 0 {
		//println("==")
		return false
	}

	return true
}
