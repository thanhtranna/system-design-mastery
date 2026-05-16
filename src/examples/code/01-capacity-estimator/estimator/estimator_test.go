package estimator

import (
	"math"
	"testing"
)

func TestEstimate_TwitterShape(t *testing.T) {
	// Twitter-scale: 200M DAU, 2 posts/user/day, 100x read amplification
	s := Scenario{
		Name: "Twitter",
		DAU:  200_000_000,
		Actions: Actions{
			WritesPerUserPerDay: 2,
			ReadsPerUserPerDay:  50, // baseline reads per user
		},
		AvgWriteBytes:     280,
		ReadAmplification: 2,  // each read fetches ~2 tweets on average
		PeakMultiplier:    3,
		ReplicationFactor: 3,
		RetentionDays:     365 * 5,
	}

	r := Estimate(s)

	// Writes/sec avg: 200M * 2 / 86400 ≈ 4629
	if math.Abs(r.WritesPerSecAvg-4629.6) > 1 {
		t.Errorf("WritesPerSecAvg = %f, want ~4629", r.WritesPerSecAvg)
	}
	// Peak = avg * 3 ≈ 13888
	if math.Abs(r.WritesPerSecPeak-13888.8) > 1 {
		t.Errorf("WritesPerSecPeak = %f, want ~13888", r.WritesPerSecPeak)
	}

	// Reads/sec avg: 200M * 50 * 2 / 86400 ≈ 231481
	if math.Abs(r.ReadsPerSecAvg-231481) > 100 {
		t.Errorf("ReadsPerSecAvg = %f, want ~231481", r.ReadsPerSecAvg)
	}

	// Storage/day = 4629.6 * 86400 * 280 * 3 ≈ 336 GB/day
	expectedDayBytes := int64(4629.6 * 86400 * 280 * 3)
	if math.Abs(float64(r.StoragePerDay-expectedDayBytes)) > 1_000_000 {
		t.Errorf("StoragePerDay = %d, want ~%d", r.StoragePerDay, expectedDayBytes)
	}
}

func TestEstimate_HandlesZeroReadAmplification(t *testing.T) {
	s := Scenario{
		Name:          "default",
		DAU:           1000,
		Actions:       Actions{WritesPerUserPerDay: 1, ReadsPerUserPerDay: 10},
		AvgWriteBytes: 100,
		// ReadAmplification omitted (zero) — should default to 1
	}
	r := Estimate(s)
	if r.ReadsPerSecAvg == 0 {
		t.Errorf("ReadsPerSecAvg should be non-zero with default amplification, got 0")
	}
}

func TestEstimate_HandlesPeakMultiplierMinimum(t *testing.T) {
	s := Scenario{
		DAU: 1000, Actions: Actions{WritesPerUserPerDay: 1, ReadsPerUserPerDay: 1},
		AvgWriteBytes: 100,
		PeakMultiplier: 0.5, // invalid - should be clamped to 1
	}
	r := Estimate(s)
	if r.WritesPerSecPeak != r.WritesPerSecAvg {
		t.Errorf("With PeakMultiplier < 1, peak should equal avg; got peak=%f avg=%f",
			r.WritesPerSecPeak, r.WritesPerSecAvg)
	}
}

func TestValidate_RejectsZeroDAU(t *testing.T) {
	s := Scenario{DAU: 0, AvgWriteBytes: 100}
	if err := s.Validate(); err == nil {
		t.Error("expected error for DAU=0, got nil")
	}
}

func TestValidate_RejectsZeroWriteBytes(t *testing.T) {
	s := Scenario{DAU: 100, AvgWriteBytes: 0}
	if err := s.Validate(); err == nil {
		t.Error("expected error for AvgWriteBytes=0, got nil")
	}
}

// Benchmark — useful in interviews to mention "I measured X"
func BenchmarkEstimate(b *testing.B) {
	s := Scenario{
		DAU: 1_000_000_000,
		Actions: Actions{WritesPerUserPerDay: 100, ReadsPerUserPerDay: 1000},
		AvgWriteBytes:     500,
		ReadAmplification: 10,
		PeakMultiplier:    5,
		ReplicationFactor: 3,
		RetentionDays:     365 * 3,
	}
	for i := 0; i < b.N; i++ {
		_ = Estimate(s)
	}
}
