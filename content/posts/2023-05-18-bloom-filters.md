---
layout: post
title: "Bloom filters"
subtitle: "Push based vs pull based query engine for OLAP"
date: 2023-04-14
author: "Niklas Hansson"
URL: "/2023/05/17push_based_query_engine"
---


Bloom filters is a probabilistic data structure, which is space efficient. Bloom filters can be used to quickly check if a value don't exists or might exists in a set, false positives are possible(with a low likelihood) but false negatives are not possible. The time to check if an element exsist or add an element is also constant O(k), where k is the number of hash functions(we will covert this later). 


## How does it works?
Bloom filters build upon `k` hash functions and an bit array of length `m`. Each value is hashed `k` times, where each hash is converted to a position in the bit array, and these bits are set to positive(if they are already positive they are not swapped). In order to check if a value exists it is hashed using the `k` hashing functions, if all the bits are set to 1 the value exists(or we have a false positive). 

Example: 

$$
x = "Thomas" \\
f^1(x) = 79 \\
f^2(x) = 1 \\
f^k(x) = 1000 \\
$$

The bits are set to 1 for each of the hash values. 


## Implementation

The implementation is done in GO in this case, in this case I implemented it as a struct with three methods: 

- hash(Not exported) - for hashing the input values
- Add(exported) - adding a value to the bit map
- Check(exported) - check if a value exists in the filter

The struct is defined as: 

``` go 
type bloom struct {
	hashes               int // nbr of hash functions
	maxVal               uint64 // the length of the bit array(bool array in this case)
	hashDifferentiations []string // trick to create n hash functions
	h                    hash.Hash // the base hash functions
	bloomFilter          []bool // the "bit array" which stores the hash values
}
```

to initialize a bloom filter a `New` function is created also to avoid having to export the struct it self. 

```go
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
```

The key for a bloom filter is the hash function in this case I emulate `k` hash function by making changes to the input in a deterministic way using the `hashDifferentiations`, in this case the upper limit of `k` is `len(hashDifferentiations)`

```go
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
```

the `hash` function is never exported byt used internally inside `Add` and `check`. The function to add values to the bloom filters is: 

```go
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
```

The bit map is emulated as a bool slice(since we don't know the size are compile time). The final function of interest is the function to check if a value exists(false positives are positive). 

```go 
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
```

The more values we add to the bloom filter the higher is the chance that we get a false positive. Since more and more of the bits in the bit array are positive. 

Some tests: 

```go
func TestNew(t *testing.T) {
	assert.NotNil(t, New(10, 10))
}

func TestHash(t *testing.T) {
	t.Run("Single hash", func(t *testing.T) {
		bloom := New(10, 10000)
		assert.NoError(t, bloom.Add([]byte("Test1")))
	})
	t.Run("Multiple hashes hash", func(t *testing.T) {
		bloom := New(10, 10000)
		for i := 0; i < 10; i++ {
			assert.NoError(t, bloom.Add([]byte(fmt.Sprintf("Test%v", i))))
		}
	})
}

func TestCheck(t *testing.T) {
	t.Run("Single hash", func(t *testing.T) {
		bloom := New(10, 10000)
		val := []byte("Test1")
		assert.NoError(t, bloom.Add(val))
		exsist, err := bloom.Check(val)
		assert.NoError(t, err)
		assert.True(t, exsist)

	})
	t.Run("Multiple hashes hash", func(t *testing.T) {
		bloom := New(10, 10000)
		for i := 0; i < 10; i++ {
			val := []byte(fmt.Sprintf("Test%v", i))
			assert.NoError(t, bloom.Add(val))
			exsist, err := bloom.Check(val)
			assert.NoError(t, err)
			assert.True(t, exsist)
			exsist, err = bloom.Check(append(val, []byte("n")...))
			assert.NoError(t, err)
			assert.False(t, exsist)
		}
	})
}

```