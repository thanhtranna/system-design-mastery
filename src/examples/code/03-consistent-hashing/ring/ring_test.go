package ring_test

import (
	"fmt"
	"testing"

	"github.com/sd-book/03-consistent-hashing/ring"
)

func TestLookupReturnsConsistentNode(t *testing.T) {
	r := ring.New(150)
	r.AddNode("node-A")
	r.AddNode("node-B")
	r.AddNode("node-C")

	node := r.Lookup("my-key")
	if node == "" {
		t.Fatal("expected non-empty node")
	}
	// same key always maps to the same node
	for range 100 {
		if got := r.Lookup("my-key"); got != node {
			t.Fatalf("inconsistent: got %s, want %s", got, node)
		}
	}
}

func TestEmptyRing(t *testing.T) {
	r := ring.New(150)
	if got := r.Lookup("any-key"); got != "" {
		t.Fatalf("expected empty string for empty ring, got %s", got)
	}
}

func TestRemapping(t *testing.T) {
	const keys = 10_000
	const virtualNodes = 150

	r := ring.New(virtualNodes)
	r.AddNode("node-A")
	r.AddNode("node-B")
	r.AddNode("node-C")

	// record initial assignments
	initial := make(map[string]string, keys)
	for i := range keys {
		k := fmt.Sprintf("key-%d", i)
		initial[k] = r.Lookup(k)
	}

	// add a fourth node
	r.AddNode("node-D")

	remapped := 0
	for i := range keys {
		k := fmt.Sprintf("key-%d", i)
		if r.Lookup(k) != initial[k] {
			remapped++
		}
	}

	pct := float64(remapped) / float64(keys) * 100
	t.Logf("Remapped %d/%d keys (%.1f%%) after adding node-D (ideal: 25%%)", remapped, keys, pct)

	// allow ±10% around the ideal 25%
	if pct < 15 || pct > 35 {
		t.Errorf("remapping percentage %.1f%% is outside expected range [15%%, 35%%]", pct)
	}
}

func TestDistribution(t *testing.T) {
	const keys = 10_000
	r := ring.New(150)
	r.AddNode("node-A")
	r.AddNode("node-B")
	r.AddNode("node-C")

	counts := make(map[string]int)
	for i := range keys {
		node := r.Lookup(fmt.Sprintf("key-%d", i))
		counts[node]++
	}

	for node, count := range counts {
		pct := float64(count) / float64(keys) * 100
		t.Logf("%s: %d keys (%.1f%%)", node, count, pct)
		if pct < 25 || pct > 45 {
			t.Errorf("node %s has %.1f%% of keys — distribution too uneven", node, pct)
		}
	}
}
