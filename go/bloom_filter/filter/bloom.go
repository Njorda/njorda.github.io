package filter

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
)

type bloom struct {
	hashes               int
	maxVal               uint64
	hashDifferentiations []string
	h                    hash.Hash
	bloomFilter          []bool
}

func New(hashes int, cardinality uint64) *bloom {
	bloomFilter := []bool{}
	for i := 0; i < int(cardinality); i++ {
		bloomFilter = append(bloomFilter, false)
	}
	return &bloom{
		hashes:               hashes,
		maxVal:               cardinality,
		hashDifferentiations: []string{"a", "b", "c ", "d", "e", "f", "g", "h", "i", "j"},
		h:                    sha1.New(),
		bloomFilter:          bloomFilter,
	}
}

func (b *bloom) hash(input []byte) ([]uint64, error) {
	defer b.h.Reset()
	out := []uint64{}
	if b.hashes > len(b.hashDifferentiations) {
		return nil, fmt.Errorf("max hashes is: %v", len(b.hashDifferentiations))
	}
	for i := 0; i < b.hashes; i++ {
		_, err := b.h.Write(append(input, []byte(b.hashDifferentiations[i])...))
		if err != nil {
			return nil, err
		}
		out = append(out, binary.BigEndian.Uint64(b.h.Sum(nil))%b.maxVal)
	}
	return out, nil
}

func (b *bloom) Add(input []byte) error {
	hashes, err := b.hash(input)
	if err != nil {
		return fmt.Errorf("hash: %w", err)
	}
	for _, idx := range hashes {
		b.bloomFilter[idx] = true
	}
	return nil
}

func (b *bloom) Check(input []byte) (bool, error) {
	hashes, err := b.hash(input)
	if err != nil {
		return false, fmt.Errorf("hash: %w", err)
	}
	for _, idx := range hashes {
		if !b.bloomFilter[idx] {
			return false, nil
		}
	}
	return true, nil
}
