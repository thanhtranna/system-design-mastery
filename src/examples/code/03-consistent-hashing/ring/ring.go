package ring

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"sync"
)

// Ring is a consistent hash ring with virtual nodes.
type Ring struct {
	mu           sync.RWMutex
	virtualNodes int
	points       []point // sorted by hash
}

type point struct {
	hash uint64
	node string
}

// New creates a ring with the given number of virtual nodes per physical node.
// 150 is a reasonable default for even distribution.
func New(virtualNodes int) *Ring {
	return &Ring{virtualNodes: virtualNodes}
}

// AddNode adds a physical node to the ring.
func (r *Ring) AddNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.virtualNodes {
		h := hashKey(fmt.Sprintf("%s#%d", node, i))
		r.points = append(r.points, point{hash: h, node: node})
	}
	sort.Slice(r.points, func(i, j int) bool {
		return r.points[i].hash < r.points[j].hash
	})
}

// RemoveNode removes a physical node and all its virtual nodes from the ring.
func (r *Ring) RemoveNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := r.points[:0]
	for _, p := range r.points {
		if p.node != node {
			filtered = append(filtered, p)
		}
	}
	r.points = filtered
}

// Lookup returns the node responsible for the given key.
// Returns empty string if the ring is empty.
func (r *Ring) Lookup(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.points) == 0 {
		return ""
	}

	h := hashKey(key)
	idx := sort.Search(len(r.points), func(i int) bool {
		return r.points[i].hash >= h
	})
	// wrap around: if no point has hash >= h, use the first point
	if idx == len(r.points) {
		idx = 0
	}
	return r.points[idx].node
}

// Nodes returns the set of distinct physical nodes currently in the ring.
func (r *Ring) Nodes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	seen := make(map[string]struct{})
	for _, p := range r.points {
		seen[p.node] = struct{}{}
	}
	nodes := make([]string, 0, len(seen))
	for n := range seen {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)
	return nodes
}

func hashKey(key string) uint64 {
	h := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint64(h[:8])
}
