// Package wal implements a minimal write-ahead log.
// Every entry is length-prefixed and checksummed; corrupt or truncated
// entries are detected during recovery.
package wal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"sync"
)

// Entry is a single WAL record.
type Entry struct {
	Seq  uint64
	Data []byte
}

// WAL is an append-only log backed by a single file.
type WAL struct {
	mu   sync.Mutex
	f    *os.File
	seq  uint64
	path string
}

// Open opens or creates a WAL at the given path.
// Existing entries are NOT replayed here — call Recover for that.
func Open(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open wal: %w", err)
	}
	return &WAL{f: f, path: path}, nil
}

// Append writes an entry to the log and syncs to disk before returning.
// The caller's state machine should only be updated after Append returns nil.
func (w *WAL) Append(data []byte) (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.seq++
	if err := writeRecord(w.f, w.seq, data); err != nil {
		return 0, err
	}
	// fdatasync: the durability guarantee
	if err := w.f.Sync(); err != nil {
		return 0, fmt.Errorf("sync: %w", err)
	}
	return w.seq, nil
}

// Recover reads all valid entries from the log file.
// Truncated or corrupt entries (e.g. from a crash mid-write) are silently stopped at.
func Recover(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	for {
		entry, err := readRecord(f)
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			break // clean end or truncation — stop here
		}
		if err != nil {
			return nil, fmt.Errorf("wal corrupt: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Close closes the underlying file.
func (w *WAL) Close() error {
	return w.f.Close()
}

// record format:
//   [8 bytes seq] [4 bytes data_len] [4 bytes crc32] [data_len bytes data]
const headerSize = 8 + 4 + 4

func writeRecord(w io.Writer, seq uint64, data []byte) error {
	checksum := crc32.ChecksumIEEE(data)
	header := make([]byte, headerSize)
	binary.BigEndian.PutUint64(header[0:8], seq)
	binary.BigEndian.PutUint32(header[8:12], uint32(len(data)))
	binary.BigEndian.PutUint32(header[12:16], checksum)

	if _, err := w.Write(header); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}

func readRecord(r io.Reader) (Entry, error) {
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return Entry{}, err
	}
	seq := binary.BigEndian.Uint64(header[0:8])
	dataLen := binary.BigEndian.Uint32(header[8:12])
	checksum := binary.BigEndian.Uint32(header[12:16])

	data := make([]byte, dataLen)
	if _, err := io.ReadFull(r, data); err != nil {
		return Entry{}, err
	}
	if crc32.ChecksumIEEE(data) != checksum {
		return Entry{}, fmt.Errorf("checksum mismatch for seq %d", seq)
	}
	return Entry{Seq: seq, Data: data}, nil
}
