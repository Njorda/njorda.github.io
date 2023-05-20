package filter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
