package circuitbreaker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// mockClock controls the time source for deterministic tests.
type mockClock struct {
	now atomic.Int64 // nanoseconds since epoch
}

func newMockClock(start time.Time) *mockClock {
	c := &mockClock{}
	c.now.Store(start.UnixNano())
	return c
}

func (c *mockClock) Now() time.Time {
	return time.Unix(0, c.now.Load())
}

func (c *mockClock) Advance(d time.Duration) {
	c.now.Add(int64(d))
}

func TestBreaker_StartsClosed(t *testing.T) {
	b := New(Config{Name: "test"})
	if got := b.State(); got != StateClosed {
		t.Errorf("initial state = %v, want CLOSED", got)
	}
}

func TestBreaker_TripsAfterFailureRatio(t *testing.T) {
	clock := newMockClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	b := New(Config{
		Name:                  "test",
		MinRequestThreshold:   10,
		FailureRatioThreshold: 0.5,
		OpenStateDuration:     time.Second,
		HalfOpenMaxProbes:     2,
		WindowSize:            time.Second,
		Now:                   clock.Now,
	})
	ctx := context.Background()

	// 10 failures should trip the breaker.
	for i := 0; i < 10; i++ {
		_, _ = b.Do(ctx, func(_ context.Context) (any, error) {
			return nil, errors.New("downstream failed")
		})
	}
	if got := b.State(); got != StateOpen {
		t.Errorf("state after 10 failures = %v, want OPEN", got)
	}
}

func TestBreaker_StaysClosedBelowMinThreshold(t *testing.T) {
	clock := newMockClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	b := New(Config{
		MinRequestThreshold:   20,
		FailureRatioThreshold: 0.5,
		WindowSize:            time.Second,
		Now:                   clock.Now,
	})
	ctx := context.Background()

	// Only 10 failures — below min threshold of 20.
	for i := 0; i < 10; i++ {
		_, _ = b.Do(ctx, func(_ context.Context) (any, error) {
			return nil, errors.New("fail")
		})
	}
	if got := b.State(); got != StateClosed {
		t.Errorf("state below min threshold = %v, want CLOSED", got)
	}
}

func TestBreaker_TransitionsThroughHalfOpenOnSuccess(t *testing.T) {
	clock := newMockClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	b := New(Config{
		MinRequestThreshold:   5,
		FailureRatioThreshold: 0.5,
		OpenStateDuration:     5 * time.Second,
		HalfOpenMaxProbes:     2,
		WindowSize:            time.Second,
		Now:                   clock.Now,
	})
	ctx := context.Background()

	// Trip breaker.
	for i := 0; i < 5; i++ {
		_, _ = b.Do(ctx, func(_ context.Context) (any, error) {
			return nil, errors.New("fail")
		})
	}
	if got := b.State(); got != StateOpen {
		t.Fatalf("expected OPEN, got %v", got)
	}

	// While open, calls should fast-fail.
	_, err := b.Do(ctx, func(_ context.Context) (any, error) {
		t.Error("function should not be called while OPEN")
		return nil, nil
	})
	if !errors.Is(err, ErrOpen) {
		t.Errorf("expected ErrOpen, got %v", err)
	}

	// Advance clock past cooldown.
	clock.Advance(6 * time.Second)

	// Probe 1: success → still HALF_OPEN (need 2 probes).
	_, err = b.Do(ctx, func(_ context.Context) (any, error) {
		return nil, nil
	})
	if err != nil {
		t.Errorf("probe 1 returned err: %v", err)
	}
	if got := b.State(); got != StateHalfOpen {
		t.Errorf("after 1 probe = %v, want HALF_OPEN", got)
	}

	// Probe 2: success → CLOSED.
	_, err = b.Do(ctx, func(_ context.Context) (any, error) {
		return nil, nil
	})
	if err != nil {
		t.Errorf("probe 2 returned err: %v", err)
	}
	if got := b.State(); got != StateClosed {
		t.Errorf("after 2 probes = %v, want CLOSED", got)
	}
}

func TestBreaker_HalfOpenFailureReopens(t *testing.T) {
	clock := newMockClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	b := New(Config{
		MinRequestThreshold:   5,
		FailureRatioThreshold: 0.5,
		OpenStateDuration:     time.Second,
		HalfOpenMaxProbes:     3,
		WindowSize:            time.Second,
		Now:                   clock.Now,
	})
	ctx := context.Background()

	// Trip.
	for i := 0; i < 5; i++ {
		_, _ = b.Do(ctx, func(_ context.Context) (any, error) {
			return nil, errors.New("fail")
		})
	}
	clock.Advance(2 * time.Second)

	// First probe fails → straight back to OPEN.
	_, _ = b.Do(ctx, func(_ context.Context) (any, error) {
		return nil, errors.New("still failing")
	})
	if got := b.State(); got != StateOpen {
		t.Errorf("after failed probe = %v, want OPEN", got)
	}
}

func TestBreaker_RespectsContext(t *testing.T) {
	b := New(Config{
		MinRequestThreshold:   1,
		FailureRatioThreshold: 0.5,
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := b.Do(ctx, func(c context.Context) (any, error) {
		return nil, c.Err()
	})
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}
