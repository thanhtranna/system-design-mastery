package vclock

import "testing"

func TestNewIsEmpty(t *testing.T) {
	c := New()
	if len(c) != 0 {
		t.Errorf("New should be empty, got %v", c)
	}
}

func TestIncrement(t *testing.T) {
	c := New().Increment("A").Increment("A").Increment("B")
	if c["A"] != 2 || c["B"] != 1 {
		t.Errorf("got %v, want A=2 B=1", c)
	}
}

func TestMerge(t *testing.T) {
	a := VC{"A": 3, "B": 1}
	b := VC{"A": 1, "B": 4, "C": 2}
	a.Merge(b)
	if a["A"] != 3 || a["B"] != 4 || a["C"] != 2 {
		t.Errorf("merge got %v, want A=3 B=4 C=2", a)
	}
}

func TestHappensBefore_StrictlyLess(t *testing.T) {
	a := VC{"A": 1, "B": 2}
	b := VC{"A": 2, "B": 2}
	if !a.HappensBefore(b) {
		t.Errorf("expected a happens-before b: a=%v b=%v", a, b)
	}
	if b.HappensBefore(a) {
		t.Errorf("did NOT expect b happens-before a")
	}
}

func TestHappensBefore_Equal(t *testing.T) {
	a := VC{"A": 1, "B": 2}
	b := VC{"A": 1, "B": 2}
	if a.HappensBefore(b) {
		t.Errorf("equal clocks should NOT happens-before each other")
	}
}

func TestConcurrent(t *testing.T) {
	a := VC{"A": 2, "B": 1}
	b := VC{"A": 1, "B": 2}
	if !a.Concurrent(b) {
		t.Errorf("a=%v b=%v should be concurrent", a, b)
	}
}

func TestDominates(t *testing.T) {
	a := VC{"A": 5, "B": 3, "C": 1}
	b := VC{"A": 3, "B": 2}
	if !a.Dominates(b) {
		t.Errorf("a=%v should dominate b=%v", a, b)
	}
}

// Real-world-shaped scenario: 3 nodes, concurrent writes.
func TestScenario_ConcurrentWrites(t *testing.T) {
	// node1 starts: empty
	c1 := New().Increment("node1") // {node1: 1}

	// node2 starts: empty (no replication yet)
	c2 := New().Increment("node2") // {node2: 1}

	if !c1.Concurrent(c2) {
		t.Errorf("c1=%v and c2=%v should be concurrent", c1, c2)
	}

	// node3 reads both, reconciles them, then writes
	merged := c1.Copy().Merge(c2)             // {node1:1, node2:1}
	c3 := merged.Increment("node3")           // {node1:1, node2:1, node3:1}

	if !c1.HappensBefore(c3) {
		t.Errorf("c1=%v should happens-before c3=%v", c1, c3)
	}
	if !c2.HappensBefore(c3) {
		t.Errorf("c2=%v should happens-before c3=%v", c2, c3)
	}
}
