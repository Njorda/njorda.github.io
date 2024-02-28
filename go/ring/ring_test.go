package ring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	r := New()
	err := r.Add("node2", 20)
	require.NoError(t, err)
	err = r.Add("node3", 20)
	require.NoError(t, err)
}

func TestGet(t *testing.T) {
	r := New()
	err := r.Add("node2", 20)
	require.NoError(t, err)
	err = r.Add("node3", 20)
	require.NoError(t, err)

	key, err := r.GetNode("hello")
	require.NoError(t, err)
	require.Equal(t, "node2", key)
	require.Len(t, r.nodes, 40)
}
