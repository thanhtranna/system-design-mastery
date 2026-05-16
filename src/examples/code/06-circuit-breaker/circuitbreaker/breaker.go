// Package circuitbreaker implements the three-state circuit breaker pattern
// (Closed / Open / Half-Open) over a rolling time window.
package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ErrOpen is returned when the breaker is in the OPEN state.
var ErrOpen = errors.New("circuit breaker is open")

// ErrTooManyProbes is returned when half-open state has admitted its
// configured max probe requests already.
var ErrTooManyProbes = errors.New("circuit breaker half-open: probe limit reached")

// State is the breaker state.
type State int32

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	}
	return "UNKNOWN"
}

// Config configures a Breaker.
type Config struct {
	Name                  string
	MinRequestThreshold   int           // min requests in window before evaluating
	FailureRatioThreshold float64       // 0..1 — trip if ratio >= this
	OpenStateDuration     time.Duration // how long to stay OPEN
	HalfOpenMaxProbes     int           // probe requests allowed in HALF_OPEN
	WindowSize            time.Duration // rolling window for stats
	BucketCount           int           // sub-buckets in the window (default 10)
	Now                   func() time.Time
}

func (c *Config) defaults() error {
	if c.Name == "" {
		c.Name = "default"
	}
	if c.MinRequestThreshold <= 0 {
		c.MinRequestThreshold = 20
	}
	if c.FailureRatioThreshold <= 0 || c.FailureRatioThreshold > 1 {
		c.FailureRatioThreshold = 0.5
	}
	if c.OpenStateDuration <= 0 {
		c.OpenStateDuration = 30 * time.Second
	}
	if c.HalfOpenMaxProbes <= 0 {
		c.HalfOpenMaxProbes = 3
	}
	if c.WindowSize <= 0 {
		c.WindowSize = 10 * time.Second
	}
	if c.BucketCount <= 0 {
		c.BucketCount = 10
	}
	if c.Now == nil {
		c.Now = time.Now
	}
	return nil
}

// Breaker is a circuit breaker.
type Breaker struct {
	cfg Config

	mu              sync.Mutex
	state           State
	openedAt        time.Time
	halfOpenProbes  int32 // atomic
	halfOpenInFlight int32 // atomic
	stats           *rollingStats
}

// New creates a new Breaker with the provided config.
func New(cfg Config) *Breaker {
	_ = cfg.defaults()
	return &Breaker{
		cfg:   cfg,
		state: StateClosed,
		stats: newRollingStats(cfg.WindowSize, cfg.BucketCount, cfg.Now),
	}
}

// State returns the current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Stats returns a snapshot of recent counts in the rolling window.
func (b *Breaker) Stats() (success, failure int) {
	return b.stats.snapshot()
}

// Do executes fn through the breaker. If the breaker is open, it returns
// ErrOpen without calling fn. In half-open state it admits up to
// HalfOpenMaxProbes concurrent probes.
func (b *Breaker) Do(ctx context.Context,
	fn func(ctx context.Context) (any, error)) (any, error) {

	if err := b.allow(); err != nil {
		return nil, err
	}

	result, err := fn(ctx)

	// Don't count context-cancellation errors against the downstream.
	// The downstream may have been fine; the caller gave up.
	if ctx.Err() != nil && errors.Is(err, ctx.Err()) {
		b.afterCall(true, true) // treat as "neutral" — but we record success
	} else {
		b.afterCall(err == nil, false)
	}

	return result, err
}

// allow decides whether to admit the next request and updates state.
func (b *Breaker) allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.cfg.Now()

	switch b.state {
	case StateClosed:
		return nil

	case StateOpen:
		if now.Sub(b.openedAt) >= b.cfg.OpenStateDuration {
			// Cooldown elapsed — transition to half-open
			b.state = StateHalfOpen
			atomic.StoreInt32(&b.halfOpenProbes, 0)
			atomic.StoreInt32(&b.halfOpenInFlight, 0)
			return b.admitHalfOpenLocked()
		}
		return ErrOpen

	case StateHalfOpen:
		return b.admitHalfOpenLocked()
	}
	return nil
}

func (b *Breaker) admitHalfOpenLocked() error {
	// Atomically count this probe; reject if limit reached.
	if atomic.AddInt32(&b.halfOpenProbes, 1) > int32(b.cfg.HalfOpenMaxProbes) {
		atomic.AddInt32(&b.halfOpenProbes, -1)
		return ErrTooManyProbes
	}
	atomic.AddInt32(&b.halfOpenInFlight, 1)
	return nil
}

// afterCall records the outcome and adjusts state.
func (b *Breaker) afterCall(success, neutral bool) {
	if !neutral {
		if success {
			b.stats.recordSuccess()
		} else {
			b.stats.recordFailure()
		}
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		b.maybeTripLocked()
	case StateHalfOpen:
		atomic.AddInt32(&b.halfOpenInFlight, -1)
		if !success {
			b.transitionToOpenLocked()
			return
		}
		// Successful probe. If all admitted probes have completed
		// successfully, close the breaker.
		if atomic.LoadInt32(&b.halfOpenInFlight) == 0 &&
			atomic.LoadInt32(&b.halfOpenProbes) >= int32(b.cfg.HalfOpenMaxProbes) {
			b.transitionToClosedLocked()
		}
	}
}

func (b *Breaker) maybeTripLocked() {
	s, f := b.stats.snapshot()
	total := s + f
	if total < b.cfg.MinRequestThreshold {
		return
	}
	ratio := float64(f) / float64(total)
	if ratio >= b.cfg.FailureRatioThreshold {
		b.transitionToOpenLocked()
	}
}

func (b *Breaker) transitionToOpenLocked() {
	b.state = StateOpen
	b.openedAt = b.cfg.Now()
}

func (b *Breaker) transitionToClosedLocked() {
	b.state = StateClosed
	b.stats.reset()
}

// --- String helpers ---

// String returns a human-readable description of the breaker.
func (b *Breaker) String() string {
	state := b.State()
	s, f := b.Stats()
	return fmt.Sprintf("Breaker(%s: state=%s success=%d failure=%d)",
		b.cfg.Name, state, s, f)
}
