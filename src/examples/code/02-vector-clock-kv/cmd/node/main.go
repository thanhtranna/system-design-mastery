// Command node runs one KV node with HTTP API and async replication.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thanhtranna/system-design-mastery/examples/vclock-kv/store"
	"github.com/thanhtranna/system-design-mastery/examples/vclock-kv/vclock"
)

func main() {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "node-1"
	}
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8001"
	}
	peersEnv := os.Getenv("PEERS")
	var peers []string
	if peersEnv != "" {
		peers = strings.Split(peersEnv, ",")
	}

	st := store.New(nodeID)

	srv := &server{
		nodeID: nodeID,
		store:  st,
		peers:  peers,
		client: &http.Client{Timeout: 2 * time.Second},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/kv/", srv.handleKV)
	mux.HandleFunc("/replicate", srv.handleReplicate)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("[%s] starting on %s peers=%v", nodeID, listenAddr, peers)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	nodeID string
	store  *store.Store
	peers  []string
	client *http.Client
}

type putReq struct {
	Value         string      `json:"value"`
	BasedOnClocks []vclock.VC `json:"based_on_clocks,omitempty"`
}

type getResp struct {
	Key        string            `json:"key"`
	Values     []store.Versioned `json:"values"`
	Concurrent bool              `json:"concurrent"`
}

type replicateReq struct {
	Key     string          `json:"key"`
	Version store.Versioned `json:"version"`
}

func (s *server) handleKV(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/kv/")
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var req putReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}
		v := s.store.Put(key, req.Value, req.BasedOnClocks)
		go s.replicateToPeers(key, v)
		_ = json.NewEncoder(w).Encode(v)

	case http.MethodGet:
		entry, ok := s.store.Get(key)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		resp := getResp{
			Key:        key,
			Values:     entry.Versions,
			Concurrent: len(entry.Versions) > 1,
		}
		_ = json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) handleReplicate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req replicateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	s.store.ApplyReplicated(req.Key, req.Version)
	w.WriteHeader(http.StatusAccepted)
}

func (s *server) replicateToPeers(key string, v store.Versioned) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	body, _ := json.Marshal(replicateReq{Key: key, Version: v})
	for _, peer := range s.peers {
		url := fmt.Sprintf("http://%s/replicate", peer)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if _, err := s.client.Do(req); err != nil {
			log.Printf("[%s] replicate to %s failed: %v", s.nodeID, peer, err)
		}
	}
}
