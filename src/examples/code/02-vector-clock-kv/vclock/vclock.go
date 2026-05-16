// Package vclock implements vector clocks for tracking causality across
// distributed nodes.
package vclock

import "sort"

// VC is a vector clock — node ID → counter.
type VC map[string]uint64

// New returns an empty vector clock.
func New() VC {
	return make(VC)
}

// Copy returns a deep copy.
func (a VC) Copy() VC {
	out := make(VC, len(a))
	for k, v := range a {
		out[k] = v
	}
	return out
}

// Increment increments the counter for nodeID and returns the modified clock.
func (a VC) Increment(nodeID string) VC {
	a[nodeID]++
	return a
}

// Merge takes the per-node max of two clocks.
func (a VC) Merge(b VC) VC {
	for k, v := range b {
		if v > a[k] {
			a[k] = v
		}
	}
	return a
}

// Equal returns true iff every node entry is identical.
func (a VC) Equal(b VC) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// HappensBefore returns true iff a < b in vector clock partial order:
// every entry of a is <= corresponding entry of b, AND there is at least
// one entry strictly less.
func (a VC) HappensBefore(b VC) bool {
	if a.Equal(b) {
		return false
	}
	// Check every key in either map.
	keys := keySet(a, b)
	atLeastOneStrict := false
	for _, k := range keys {
		if a[k] > b[k] {
			return false
		}
		if a[k] < b[k] {
			atLeastOneStrict = true
		}
	}
	return atLeastOneStrict
}

// Concurrent returns true iff neither a HappensBefore b nor b HappensBefore a.
func (a VC) Concurrent(b VC) bool {
	return !a.Equal(b) && !a.HappensBefore(b) && !b.HappensBefore(a)
}

// Dominates returns true iff a >= b in every entry AND a != b.
// (Same as b.HappensBefore(a).)
func (a VC) Dominates(b VC) bool {
	return b.HappensBefore(a)
}

func keySet(a, b VC) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
