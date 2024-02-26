package ring

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sort"
)

// Way to slow and uninspiring
// Time to implement consistent hashing instead
// This is my ring :)

type virtualNode struct {
	node string
	// We should be able to call the node as well here ...
	start uint32
}

type Ring struct {
	nodes []virtualNode
	hash  hash.Hash32
	// Mapping from virtual node to the actual node
	// This is needed when we delete a node
	NodeMapping map[string][]string
}

func New() *Ring {
	return &Ring{
		nodes:       []virtualNode{},
		hash:        fnv.New32(),
		NodeMapping: map[string][]string{},
	}
}

func (r *Ring) Len() int           { return len(r.nodes) }
func (r *Ring) Swap(i, j int)      { r.nodes[i], r.nodes[j] = r.nodes[j], r.nodes[i] }
func (r *Ring) Less(i, j int) bool { return r.nodes[i].start < r.nodes[j].start }

// Add adds a new node -> number of virtual nodes to the ring
func (r *Ring) Add(node string, nbrVirtualNode int) error {
	if nbrVirtualNode < 1 {
		return fmt.Errorf("nbrVirtualNode needs to be larger than  0, currently %v", nbrVirtualNode)
	}
	for i := 0; i < nbrVirtualNode; i++ {
		// Never assume the buffer is clean, be defensive
		r.hash.Reset()
		_, err := r.hash.Write([]byte(node))
		if err != nil {
			return fmt.Errorf("write hash: %v", err)
		}
		r.nodes = append(r.nodes, virtualNode{
			node:  node,
			start: r.hash.Sum32(),
		})
		sort.Sort(r)
	}
	return nil
}

// GetNode gets the node
func (r *Ring) GetNode(key string) (string, error) {
	// Never assume the buffer is clean, be defensive
	r.hash.Reset()
	_, err := r.hash.Write([]byte(key))
	if err != nil {
		return "", fmt.Errorf("write hash: %v", err)
	}
	hash := r.hash.Sum32()

	// TODO: continue here to fix the issues with the tests
	// We want to use all virtual nodes ...
	// set the first node to 0?
	// The last node will not be used currently
	// We need to make it circular
	for idx, node := range r.nodes {
		if idx == 0 {
			return r.nodes[len(r.nodes)-1].node, nil
		}
		if node.start > hash {
			return r.nodes[idx-1].node, nil
		}
	}
	return r.nodes[0].node, nil
}

// RemoveNodes removes the actual node
// Data will not be moved, we just drop the nodes
func (r *Ring) RemoveNode(node string) error {
	return nil
}
