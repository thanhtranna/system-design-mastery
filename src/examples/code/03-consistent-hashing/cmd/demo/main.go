package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/sd-book/03-consistent-hashing/ring"
)

func main() {
	interactive := flag.Bool("interactive", false, "run interactive mode")
	flag.Parse()

	if *interactive {
		runInteractive()
	} else {
		runDemo()
	}
}

func runDemo() {
	const totalKeys = 10_000
	const virtualNodes = 150

	r := ring.New(virtualNodes)
	r.AddNode("node-A")
	r.AddNode("node-B")
	r.AddNode("node-C")

	fmt.Printf("Ring with 3 nodes, %d virtual nodes each.\n", virtualNodes)
	fmt.Printf("Distributing %d keys...\n", totalKeys)

	keys := make([]string, totalKeys)
	for i := range totalKeys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	initial := assignAll(r, keys)
	printDistribution(initial, totalKeys)

	fmt.Println("\nAdding node-D...")
	r.AddNode("node-D")
	after := assignAll(r, keys)
	printRemapping(initial, after, totalKeys)

	fmt.Println("\nRemoving node-B...")
	r.RemoveNode("node-B")
	after2 := assignAll(r, keys)
	printRemapping(after, after2, totalKeys)
}

func assignAll(r *ring.Ring, keys []string) map[string]string {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = r.Lookup(k)
	}
	return m
}

func printDistribution(assignments map[string]string, total int) {
	counts := make(map[string]int)
	for _, node := range assignments {
		counts[node]++
	}
	nodes := sortedKeys(counts)
	for _, node := range nodes {
		fmt.Printf("  %s: %d keys (%.1f%%)\n", node, counts[node], float64(counts[node])/float64(total)*100)
	}
}

func printRemapping(before, after map[string]string, total int) {
	counts := make(map[string]int)
	remapped := 0
	for k, newNode := range after {
		counts[newNode]++
		if before[k] != newNode {
			remapped++
		}
	}
	nodes := sortedKeys(counts)
	for _, node := range nodes {
		delta := counts[node] - countNode(before, node)
		sign := "+"
		if delta < 0 {
			sign = ""
		}
		fmt.Printf("  %s: %d keys  (%s%d)\n", node, counts[node], sign, delta)
	}
	fmt.Printf("  Keys remapped: %d / %d (%.1f%%)\n", remapped, total, float64(remapped)/float64(total)*100)
}

func countNode(m map[string]string, node string) int {
	n := 0
	for _, v := range m {
		if v == node {
			n++
		}
	}
	return n
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func runInteractive() {
	r := ring.New(150)
	fmt.Println("Consistent Hash Ring — interactive mode")
	fmt.Println("Commands: add <node>, remove <node>, lookup <key>, nodes, quit")
	fmt.Println()

	var cmd, arg string
	for {
		fmt.Print("> ")
		_, err := fmt.Scan(&cmd)
		if err != nil || cmd == "quit" {
			break
		}
		switch cmd {
		case "add":
			fmt.Scan(&arg)
			r.AddNode(arg)
			fmt.Printf("Added %s\n", arg)
		case "remove":
			fmt.Scan(&arg)
			r.RemoveNode(arg)
			fmt.Printf("Removed %s\n", arg)
		case "lookup":
			fmt.Scan(&arg)
			fmt.Printf("%s → %s\n", arg, r.Lookup(arg))
		case "nodes":
			fmt.Println(r.Nodes())
		default:
			fmt.Println("Unknown command")
		}
	}
}
