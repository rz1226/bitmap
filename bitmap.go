package bitmap

import (
	"sync/atomic"
	"time"
)

type BitMap struct {
	Name       string
	Key        string
	SetCount   int32
	data       *BitMap2
	LastUpTime int64
}

func NewBitMap(key string, name string) *BitMap {
	b := new(BitMap)
	b.Name = name
	b.Key = key
	b.data = NewBitMap2([]byte{})
	b.SetCount = 0
	b.LastUpTime = time.Now().Unix()
	return b
}

func (b *BitMap) Get(pos int) bool {
	return b.data.Get(pos)
}

func (b *BitMap) SetTrue(pos int) {
	atomic.AddInt32(&b.SetCount, 1)
	b.LastUpTime = time.Now().Unix()
	b.data.SetTrue(pos)
}

func (b *BitMap) SetFalse(pos int) {
	atomic.AddInt32(&b.SetCount, 1)
	b.LastUpTime = time.Now().Unix()
	b.data.SetFalse(pos)
}
