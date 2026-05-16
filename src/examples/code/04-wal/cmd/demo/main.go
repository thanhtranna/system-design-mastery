package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/sd-book/04-wal/store"
)

const walPath = "/tmp/sd-wal-demo.log"

func main() {
	// clean start
	os.Remove(walPath)

	fmt.Println("=== Write-Ahead Log Demo ===")
	fmt.Println()

	fmt.Println("[1] Writing keys to store...")
	s, err := store.Open(walPath)
	if err != nil {
		panic(err)
	}
	entries := map[string]string{
		"name": "Alice",
		"city": "Singapore",
		"role": "engineer",
	}
	for k, v := range entries {
		fmt.Printf("  SET %s=%s\n", k, v)
		s.Set(k, v)
	}
	// simulate crash: close the store (WAL is synced; in a real crash the process dies here)
	s.Close()

	fmt.Println()
	fmt.Println("[2] Simulating crash (close without cleanup)...")
	fmt.Println()

	fmt.Println("[3] Recovering from WAL...")
	s2, err := store.Open(walPath)
	if err != nil {
		panic(err)
	}
	defer s2.Close()
	fmt.Printf("  Replayed %d entries\n", len(s2.Keys()))

	fmt.Println()
	fmt.Println("[4] State after recovery:")
	keys := s2.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		v, _ := s2.Get(k)
		fmt.Printf("  %-6s = %s\n", k, v)
	}

	fmt.Println()
	fmt.Println("Recovery complete. No data lost.")
}
