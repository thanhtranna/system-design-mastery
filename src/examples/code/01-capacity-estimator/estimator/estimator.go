// Package estimator computes capacity numbers from behavioral inputs.
// Pure functions — no I/O, fully unit-testable.
package estimator

import "fmt"

// Scenario describes the inputs to a capacity estimation.
type Scenario struct {
	Name              string  `yaml:"name"`
	DAU               int64   `yaml:"dau"`
	Actions           Actions `yaml:"actions"`
	AvgWriteBytes     int     `yaml:"avg_write_bytes"`
	ReadAmplification int     `yaml:"read_amplification"`
	PeakMultiplier    float64 `yaml:"peak_multiplier"`
	ReplicationFactor int     `yaml:"replication_factor"`
	RetentionDays     int     `yaml:"retention_days"`
}

// Actions describes per-user daily activity.
type Actions struct {
	WritesPerUserPerDay int `yaml:"writes_per_user_per_day"`
	ReadsPerUserPerDay  int `yaml:"reads_per_user_per_day"`
}

// Validate returns an error if the scenario is malformed.
func (s Scenario) Validate() error {
	if s.DAU <= 0 {
		return fmt.Errorf("dau must be > 0, got %d", s.DAU)
	}
	if s.Actions.WritesPerUserPerDay < 0 {
		return fmt.Errorf("writes_per_user_per_day must be >= 0")
	}
	if s.Actions.ReadsPerUserPerDay < 0 {
		return fmt.Errorf("reads_per_user_per_day must be >= 0")
	}
	if s.AvgWriteBytes <= 0 {
		return fmt.Errorf("avg_write_bytes must be > 0")
	}
	if s.ReadAmplification <= 0 {
		s.ReadAmplification = 1 // sensible default
	}
	if s.PeakMultiplier < 1 {
		s.PeakMultiplier = 1
	}
	if s.ReplicationFactor < 1 {
		s.ReplicationFactor = 1
	}
	return nil
}

// Result is the computed estimate from a Scenario.
type Result struct {
	WritesPerSecAvg   float64
	WritesPerSecPeak  float64
	ReadsPerSecAvg    float64
	ReadsPerSecPeak   float64
	StoragePerDay     int64   // bytes/day, post-replication
	StoragePerYear    int64   // bytes/year, post-replication
	StorageRetention  int64   // bytes over retention window
	BandwidthReadPeak int64   // bytes/sec peak
}

const secondsPerDay = 86400.0

// Estimate computes the derived capacity numbers from a Scenario.
func Estimate(s Scenario) Result {
	if s.ReadAmplification <= 0 {
		s.ReadAmplification = 1
	}
	if s.PeakMultiplier < 1 {
		s.PeakMultiplier = 1
	}
	if s.ReplicationFactor < 1 {
		s.ReplicationFactor = 1
	}
	if s.RetentionDays <= 0 {
		s.RetentionDays = 365
	}

	writesPerDay := float64(s.DAU) * float64(s.Actions.WritesPerUserPerDay)
	readsPerDay := float64(s.DAU) * float64(s.Actions.ReadsPerUserPerDay) * float64(s.ReadAmplification)

	wpsAvg := writesPerDay / secondsPerDay
	rpsAvg := readsPerDay / secondsPerDay

	bytesPerDay := int64(writesPerDay * float64(s.AvgWriteBytes) * float64(s.ReplicationFactor))
	bytesPerYear := bytesPerDay * 365

	return Result{
		WritesPerSecAvg:   wpsAvg,
		WritesPerSecPeak:  wpsAvg * s.PeakMultiplier,
		ReadsPerSecAvg:    rpsAvg,
		ReadsPerSecPeak:   rpsAvg * s.PeakMultiplier,
		StoragePerDay:     bytesPerDay,
		StoragePerYear:    bytesPerYear,
		StorageRetention:  bytesPerDay * int64(s.RetentionDays),
		BandwidthReadPeak: int64(rpsAvg * s.PeakMultiplier * float64(s.AvgWriteBytes)),
	}
}
