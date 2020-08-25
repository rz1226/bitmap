package bitmap

import "errors"

type BitMapSet struct {
	data map[string]*BitMap
}

func (bs *BitMapSet) Set(b *BitMap) {
	key := b.Key
	bs.data[key] = b
}

func (bs *BitMapSet) Get(key string) (*BitMap, error) {
	if b, ok := bs.data[key]; ok {
		return b, nil
	}
	return nil, errors.New("found no map")
}
