package bitFilter

import (
	"fmt"
	//"hash"
	"encoding/binary"
	"hash/fnv"
	"testing"
)

var bits *TBits
var maxMask uint32

func TestFilter(t *testing.T) {
	bits = New(128)
	maxMask = 128<<6 - 1

	if len(bits.buckets) != 128 {
		t.Fatalf("len buckets: %d != 128", len(bits.buckets))
	}

	keys := []uint64{0, 1, 2, 3, 4, 8192, 8193, 8194, 8195, 8196, 16385}
	collHashKeys := []uint32{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1}
	hashKeys := []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	bits.InitCollisionKeys(keys, collHashKeys)
	bits.InitAllHashKeys(hashKeys)

	if len(bits.collision) != len(collHashKeys) {
		t.Fatalf("collision len != %d", len(collHashKeys))
	}

	if bits.buckets[0] != 1023 {
		t.Fatalf("buckets[0]:%d != 1023", bits.buckets[0])
	}

	if bits.buckets[1] != 0 {
		t.Fatalf("buckets[1]:%d != 0", bits.buckets[1])
	}

	var (
		i       uint64
		hashKey uint32
		exist   bool
	)
	for i = 0; i < 10000; i++ {
		hashKey = uint32(i) & maxMask
		//println(i, hashKey)
		exist = bits.Filter(i, hashKey)
		if hashKey < 10 {
			if !exist {
				t.Fatalf("%d<10 should be exist! %+v", i, bits.collision)
			}
		} else {
			if exist {
				t.Fatalf("%d>=10 should be not exist!", i)
			}
		}
	}
}

func TestNew(t *testing.T) {
	newBits := bits.Copy()

	if len(newBits.collision) != 11 {
		t.Fatalf("collision len != 11")
	}

	if newBits.buckets[0] != 1023 {
		t.Fatalf("buckets[0]:%d != 1023", newBits.buckets[0])
	}

	if newBits.buckets[1] != 0 {
		t.Fatalf("buckets[1]:%d != 0", newBits.buckets[1])
	}

	newBits.Remove(8192, 0)

	if len(newBits.collision) != 9 {
		t.Fatalf("collision len:9 != %d", len(newBits.collision))
	}

	if newBits.buckets[0] != 1023 {
		t.Fatalf("buckets[0]:%d != 1023", newBits.buckets[0])
	}

	if isExist := newBits.Filter(8192, 0); !isExist {
		// 8192这个值虽然remove了，但是0这个值还没remove
		// 所以这里还是存在的，这是误判
		t.Fatalf("error")
	}

	newBits.Remove(0, 0)
	if isExist := newBits.Filter(0, 0); isExist {
		t.Fatalf("error")
	}

	if newBits.buckets[0] != 1022 {
		t.Fatalf("buckets[0]:%d != 1022", newBits.buckets[0])
	}

	newBits.Remove(16385, 1)
	if len(newBits.collision) != 8 {
		t.Fatalf("collision len:8 != %d", len(newBits.collision))
	}

	if isExist := newBits.Filter(1, 1); !isExist {
		t.Fatalf("error")
	}

	if newBits.buckets[0] != 1022 {
		t.Fatalf("buckets[0]:%d != 1022", newBits.buckets[0])
	}

}

func TestRand(t *testing.T) {
	bits = New(Size2MB)
	maxMask = uint32(Size2MB<<6 - 1)

	var exist = make(map[uint32]bool)
	var errorCount int
	var projectId, posId, n1, n2 uint32
	var b1 = make([]byte, 4)
	var b2 = make([]byte, 4)
	alg := fnv.New32a()
	n1, n2 = 68000, 100
	for projectId = 1; projectId < n1; projectId++ {
		binary.BigEndian.PutUint32(b1, projectId)
		for posId = 1; posId < n2; posId++ {
			binary.BigEndian.PutUint32(b2, posId)
			alg.Reset()
			alg.Write(append(b1, b2...))
			hashKey := alg.Sum32()
			hashKey = hashKey<<6 - 1
			if _, ok := exist[hashKey]; ok {
				errorCount++
			} else {
				exist[hashKey] = true
			}
		}
	}

	fmt.Printf("Error Count: %d in %d*%d\nError Rate: %f\n", errorCount, n1, n2, float32(errorCount)/(float32(n1)*float32(n2)))

	projectId, posId = 32314, 23413
	binary.BigEndian.PutUint32(b1, projectId)
	binary.BigEndian.PutUint32(b2, posId)
	alg.Reset()
	alg.Write(append(b1, b2...))
	fmt.Printf("fnv.New32a: %d, %d: %d\n", projectId, posId, alg.Sum32()&maxMask)
}

func BenchmarkFilter(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bits.Filter(1234, 43)
		}
	})
}
