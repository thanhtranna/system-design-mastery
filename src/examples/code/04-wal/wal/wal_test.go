package wal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sd-book/04-wal/wal"
)

func TestAppendAndRecover(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.wal")

	w, err := wal.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, msg := range []string{"hello", "world", "wal"} {
		if _, err := w.Append([]byte(msg)); err != nil {
			t.Fatal(err)
		}
	}
	w.Close()

	entries, err := wal.Recover(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, want := range []string{"hello", "world", "wal"} {
		if string(entries[i].Data) != want {
			t.Errorf("entry %d: got %q, want %q", i, entries[i].Data, want)
		}
	}
}

func TestRecoverNonexistent(t *testing.T) {
	entries, err := wal.Recover("/tmp/nonexistent-sd-wal.log")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestCrashRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crash.wal")

	// write 3 entries cleanly
	w, _ := wal.Open(path)
	w.Append([]byte("entry-1"))
	w.Append([]byte("entry-2"))
	w.Append([]byte("entry-3"))
	w.Close()

	// simulate crash: append a partial (corrupt) record directly to the file
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	f.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}) // seq=4 header, truncated
	f.Close()

	entries, err := wal.Recover(path)
	if err != nil {
		t.Fatal(err)
	}
	// should recover 3 valid entries, stop at the corrupt one
	if len(entries) != 3 {
		t.Fatalf("expected 3 recovered entries, got %d", len(entries))
	}
}
