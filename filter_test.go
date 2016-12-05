package bitFilter

import (
	"testing"
)

var bits *TBits
var maxMask uint32

func init() {
	bits = New(1024)
	maxMask = 1024<<3 - 1
}

func TestFilter(t *testing.T) {
	if len(bits.bytes) != 1024 {
		t.Fatalf("len bytes: %d != 2014", len(bits.bytes))
	}

	bits.AddCollisionKeys([]uint32{0, 1, 2, 3, 4, 8192, 8193, 8194, 8195, 8196})
	bits.AddAllHashKeys([]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	if len(bits.collision) != 10 {
		t.Fatalf("collision len != 10")
	}

	if bits.bytes[0] != 255 {
		t.Fatalf("bytes[0]:%d != 255", bits.bytes[0])
	}

	if bits.bytes[1] != 3 {
		t.Fatalf("bytes[1]:%d != 3", bits.bytes[1])
	}

	var i, hashKey uint32
	var exist bool
	for i = 0; i < 10000; i++ {
		hashKey = i & maxMask
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
