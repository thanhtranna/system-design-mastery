// Command chaos runs a chaos test of the circuit breaker against a
// flaky simulated downstream. It demonstrates the breaker's behavior
// across phases: healthy → fully broken → recovered → flaky.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	cb "github.com/thanhtranna/system-design-mastery/examples/circuitbreaker/circuitbreaker"
)

type downstream struct {
	mu          sync.Mutex
	failureRate float64 // 0..1
}

func (d *downstream) call(ctx context.Context) (any, error) {
	d.mu.Lock()
	rate := d.failureRate
	d.mu.Unlock()

	// Simulate some latency.
	select {
	case <-time.After(time.Duration(5+rand.Intn(15)) * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if rand.Float64() < rate {
		return nil, errors.New("downstream error")
	}
	return "ok", nil
}

func (d *downstream) setFailureRate(r float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.failureRate = r
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ds := &downstream{failureRate: 0.0}
	breaker := cb.New(cb.Config{
		Name:                  "downstream",
		MinRequestThreshold:   20,
		FailureRatioThreshold: 0.5,
		OpenStateDuration:     5 * time.Second,
		HalfOpenMaxProbes:     3,
		WindowSize:            3 * time.Second,
	})

	var (
		callsTotal    atomic.Int64
		callsActual   atomic.Int64
		callsFastFail atomic.Int64
		callsErr      atomic.Int64
	)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Phases.
	phases := []struct {
		name        string
		failureRate float64
		duration    time.Duration
	}{
		{"healthy", 0.00, 15 * time.Second},
		{"100% broken", 1.00, 15 * time.Second},
		{"recovered", 0.00, 15 * time.Second},
		{"flaky 30%", 0.30, 30 * time.Second},
	}

	// Phase driver.
	go func() {
		for _, p := range phases {
			log.Printf(">>> phase: %s (failure_rate=%.0f%%)",
				p.name, p.failureRate*100)
			ds.setFailureRate(p.failureRate)
			select {
			case <-time.After(p.duration):
			case <-ctx.Done():
				return
			}
		}
		cancel()
	}()

	// State logger.
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s, f := breaker.Stats()
				log.Printf("[%s]  total=%d actual=%d fast_fail=%d err=%d  (success=%d failure=%d)",
					breaker.State(),
					callsTotal.Load(), callsActual.Load(),
					callsFastFail.Load(), callsErr.Load(),
					s, f)
			}
		}
	}()

	// Load generator.
	var wg sync.WaitGroup
	concurrency := 50
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				callsTotal.Add(1)
				_, err := breaker.Do(ctx, func(c context.Context) (any, error) {
					return ds.call(c)
				})
				if errors.Is(err, cb.ErrOpen) || errors.Is(err, cb.ErrTooManyProbes) {
					callsFastFail.Add(1)
				} else {
					callsActual.Add(1)
					if err != nil {
						callsErr.Add(1)
					}
				}
				time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
			}
		}()
	}
	wg.Wait()

	fmt.Println("\n=== Final results ===")
	s, f := breaker.Stats()
	fmt.Printf("Total requests:     %d\n", callsTotal.Load())
	fmt.Printf("Actually called ds: %d\n", callsActual.Load())
	fmt.Printf("Fast-failed:        %d\n", callsFastFail.Load())
	fmt.Printf("Errors observed:    %d\n", callsErr.Load())
	fmt.Printf("Breaker final state: %s (success=%d failure=%d)\n",
		breaker.State(), s, f)
}
