package circuitbreaker

import (
	"sync"
	"time"
)

// rollingStats maintains success/failure counts over a sliding window
// using a ring of sub-buckets.
type rollingStats struct {
	mu          sync.Mutex
	buckets     []bucket
	bucketDur   time.Duration
	windowDur   time.Duration
	now         func() time.Time
	currentIdx  int
	currentTime time.Time
}

type bucket struct {
	success int
	failure int
}

func newRollingStats(window time.Duration, bucketCount int, now func() time.Time) *rollingStats {
	bd := window / time.Duration(bucketCount)
	return &rollingStats{
		buckets:     make([]bucket, bucketCount),
		bucketDur:   bd,
		windowDur:   window,
		now:         now,
		currentTime: now(),
	}
}

func (r *rollingStats) rotateLocked() {
	now := r.now()
	elapsed := now.Sub(r.currentTime)
	if elapsed < r.bucketDur {
		return
	}
	steps := int(elapsed / r.bucketDur)
	if steps >= len(r.buckets) {
		// Whole window passed; reset everything.
		for i := range r.buckets {
			r.buckets[i] = bucket{}
		}
		r.currentIdx = 0
	} else {
		for i := 0; i < steps; i++ {
			r.currentIdx = (r.currentIdx + 1) % len(r.buckets)
			r.buckets[r.currentIdx] = bucket{}
		}
	}
	r.currentTime = now
}

func (r *rollingStats) recordSuccess() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rotateLocked()
	r.buckets[r.currentIdx].success++
}

func (r *rollingStats) recordFailure() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rotateLocked()
	r.buckets[r.currentIdx].failure++
}

func (r *rollingStats) snapshot() (success, failure int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rotateLocked()
	for _, b := range r.buckets {
		success += b.success
		failure += b.failure
	}
	return
}

func (r *rollingStats) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.buckets {
		r.buckets[i] = bucket{}
	}
	r.currentIdx = 0
	r.currentTime = r.now()
}
