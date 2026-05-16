package store

import (
	"testing"

	"github.com/thanhtranna/system-design-mastery/examples/vclock-kv/vclock"
)

func TestPutAndGet(t *testing.T) {
	s := New("node-1")
	s.Put("k", "v1", nil)

	entry, ok := s.Get("k")
	if !ok {
		t.Fatal("expected key found")
	}
	if len(entry.Versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(entry.Versions))
	}
	if entry.Versions[0].Value != "v1" {
		t.Errorf("got value %q, want v1", entry.Versions[0].Value)
	}
	if entry.Versions[0].Clock["node-1"] != 1 {
		t.Errorf("expected node-1 clock 1, got %d", entry.Versions[0].Clock["node-1"])
	}
}

func TestSequentialWritesAdvanceClock(t *testing.T) {
	s := New("node-1")
	s.Put("k", "v1", nil)
	s.Put("k", "v2", nil)
	s.Put("k", "v3", nil)

	entry, _ := s.Get("k")
	if len(entry.Versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(entry.Versions))
	}
	if entry.Versions[0].Value != "v3" {
		t.Errorf("got %q, want v3", entry.Versions[0].Value)
	}
	if entry.Versions[0].Clock["node-1"] != 3 {
		t.Errorf("expected node-1 clock=3, got %d", entry.Versions[0].Clock["node-1"])
	}
}

func TestConcurrentWritesKeepBoth(t *testing.T) {
	s1 := New("node-1")
	s2 := New("node-2")

	v1 := s1.Put("k", "from-node-1", nil)
	v2 := s2.Put("k", "from-node-2", nil)

	s2.ApplyReplicated("k", v1)
	s1.ApplyReplicated("k", v2)

	e1, _ := s1.Get("k")
	if len(e1.Versions) != 2 {
		t.Errorf("s1 expected 2 versions, got %d", len(e1.Versions))
	}
	e2, _ := s2.Get("k")
	if len(e2.Versions) != 2 {
		t.Errorf("s2 expected 2 versions, got %d", len(e2.Versions))
	}
}

func TestReconciliationCollapsesVersions(t *testing.T) {
	s1 := New("node-1")
	s2 := New("node-2")

	v1 := s1.Put("k", "from-node-1", nil)
	v2 := s2.Put("k", "from-node-2", nil)

	s2.ApplyReplicated("k", v1)
	s1.ApplyReplicated("k", v2)

	s3 := New("node-3")
	s3.ApplyReplicated("k", v1)
	s3.ApplyReplicated("k", v2)

	s3.Put("k", "merged", []vclock.VC{v1.Clock, v2.Clock})

	entry, _ := s3.Get("k")
	if len(entry.Versions) != 1 {
		t.Errorf("expected 1 version after reconciliation, got %d", len(entry.Versions))
	}
	if entry.Versions[0].Value != "merged" {
		t.Errorf("got %q, want merged", entry.Versions[0].Value)
	}
}
