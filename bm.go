package bitmap

import (
	"math/rand"
	"sync"
	"time"
)

type BitMap2 struct {
	m        *sync.Mutex
	byteLine []byte
}

func NewBitMap2(dataInit []byte) *BitMap2 {
	bm := &BitMap2{}
	bm.byteLine = dataInit
	bm.m = new(sync.Mutex)
	return bm
}

func (this *BitMap2) Bytes() []byte {
	return this.byteLine
}
func (this *BitMap2) String() string {
	return string(this.byteLine)
}

//加长，如果已经够长了，什么都不操作,否则全部补充零
//注意该长度是指byte的数量长度,不是位的长度
func (this *BitMap2) padWithZero(lenth int) {
	clen := len(this.byteLine)
	if clen < lenth {
		newBiggerByteLine := make([]byte, lenth)
		copy(newBiggerByteLine, this.byteLine)
		this.byteLine = newBiggerByteLine
	}
}

//对外接口
func (this *BitMap2) Get(position int) bool {
	this.m.Lock()
	defer this.m.Unlock()
	pos := position / 8
	if pos > len(this.byteLine)-1 {
		return false
	}
	value := this.byteLine[pos]
	pos2 := position % 8
	newValue := getSingleBytePositionValue(value, pos2)
	return newValue

}

//对外接口
func (this *BitMap2) SetTrue(position int) {
	this.setPostion(position, true)
}
func (this *BitMap2) SetFalse(position int) {
	this.setPostion(position, false)
}

func (this *BitMap2) setPostion(position int, val bool) {
	this.m.Lock()
	defer this.m.Unlock()
	whichByte := position / 8
	oriLen := len(this.byteLine)
	if whichByte > oriLen-1 {
		//多增加的部分为了减少pad的次数，每次pad要拷贝数据性能很低
		this.padWithZero(whichByte + 10 + oriLen/5)
	}
	value := this.byteLine[whichByte]
	mod := position % 8
	newValue := setSingleBytePositionValue(value, mod, val)
	this.update(whichByte, newValue)
}
func (this *BitMap2) Len() int {
	return len(this.byteLine) * 8
}

//update
func (this *BitMap2) update(pos int, val byte) {
	this.byteLine[pos] = val
}

func minLen(t, t2 *BitMap2) int {
	if len(t.byteLine) <= len(t2.byteLine) {
		return len(t.byteLine)
	}
	return len(t2.byteLine)
}

func maxLen(t, t2 *BitMap2) int {
	if len(t.byteLine) >= len(t2.byteLine) {
		return len(t.byteLine)
	}
	return len(t2.byteLine)
}

func (b *BitMap2) Or(b2 *BitMap2) *BitMap2 {
	len := maxLen(b, b2)
	b.padWithZero(len)
	b2.padWithZero(len)
	bm := &BitMap2{}
	bm.byteLine = make([]byte, len)
	for i := 0; i < len; i++ {
		bm.byteLine[i] = b.byteLine[i] | b2.byteLine[i]
	}
	return bm
}

func (b *BitMap2) And(b2 *BitMap2) *BitMap2 {
	len := minLen(b, b2)
	bm := &BitMap2{}
	bm.byteLine = make([]byte, len)
	for i := 0; i < len; i++ {
		bm.byteLine[i] = b.byteLine[i] & b2.byteLine[i]
	}
	return bm
}

//获取byte某个位置的值
func getSingleBytePositionValue(value byte, bitpos int) bool {
	if bitpos < 0 || bitpos > 7 {
		return false
	}
	reversedPosition := byte(7 - bitpos)
	factor := byte(1) << reversedPosition
	if factor == factor&value {
		return true
	}
	return false
}

//设置某个位置的值，返回新的byte
func setSingleBytePositionValue(haystack byte, bitpos int, val bool) byte {
	if bitpos < 0 || bitpos > 7 {
		return haystack
	}
	reversedPosition := byte(7 - bitpos)
	if val == true {
		factor := byte(1) << reversedPosition
		if factor == factor&haystack {
			return haystack
		} else {
			return haystack | factor
		}
	} else {
		factor := byte(1) << reversedPosition
		if byte(0) == factor&haystack {
			return haystack
		} else {
			return haystack & ^factor
		}
	}
}

//用于测试
func ByteToBinaryString(data byte) (str string) {
	var a byte
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1

		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

func getRand() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(255)
}
