package estimator

import (
	"fmt"
	"strings"
)

// FormatText returns a human-readable text report.
func FormatText(s Scenario, r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", s.Name)
	fmt.Fprintf(&b, "%s\n\n", strings.Repeat("=", len(s.Name)))

	fmt.Fprintln(&b, "Inputs:")
	fmt.Fprintf(&b, "  DAU:                %s\n", humanInt64(s.DAU))
	fmt.Fprintf(&b, "  Writes/user/day:    %d\n", s.Actions.WritesPerUserPerDay)
	fmt.Fprintf(&b, "  Reads/user/day:     %d\n", s.Actions.ReadsPerUserPerDay)
	fmt.Fprintf(&b, "  Avg write size:     %s\n", humanBytes(int64(s.AvgWriteBytes)))
	fmt.Fprintf(&b, "  Read amplification: %dx\n", s.ReadAmplification)
	fmt.Fprintf(&b, "  Peak multiplier:    %.1fx\n", s.PeakMultiplier)
	fmt.Fprintf(&b, "  Replication:        %dx\n", s.ReplicationFactor)
	fmt.Fprintf(&b, "  Retention:          %d days\n\n", s.RetentionDays)

	fmt.Fprintln(&b, "Derived:")
	fmt.Fprintf(&b, "  Writes/sec avg:     %s\n", humanFloat(r.WritesPerSecAvg))
	fmt.Fprintf(&b, "  Writes/sec peak:    %s\n", humanFloat(r.WritesPerSecPeak))
	fmt.Fprintf(&b, "  Reads/sec avg:      %s\n", humanFloat(r.ReadsPerSecAvg))
	fmt.Fprintf(&b, "  Reads/sec peak:     %s\n", humanFloat(r.ReadsPerSecPeak))
	fmt.Fprintf(&b, "  Storage/day:        %s\n", humanBytes(r.StoragePerDay))
	fmt.Fprintf(&b, "  Storage/year:       %s\n", humanBytes(r.StoragePerYear))
	fmt.Fprintf(&b, "  Storage retention:  %s\n", humanBytes(r.StorageRetention))
	fmt.Fprintf(&b, "  Bandwidth peak:     %s/sec\n", humanBytes(r.BandwidthReadPeak))

	return b.String()
}

// FormatMarkdown returns a markdown-table report suitable for design docs.
func FormatMarkdown(s Scenario, r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Capacity Estimate: %s\n\n", s.Name)

	fmt.Fprintln(&b, "### Inputs")
	fmt.Fprintln(&b, "| Parameter | Value |")
	fmt.Fprintln(&b, "|---|---|")
	fmt.Fprintf(&b, "| DAU | %s |\n", humanInt64(s.DAU))
	fmt.Fprintf(&b, "| Writes/user/day | %d |\n", s.Actions.WritesPerUserPerDay)
	fmt.Fprintf(&b, "| Reads/user/day | %d |\n", s.Actions.ReadsPerUserPerDay)
	fmt.Fprintf(&b, "| Avg write size | %s |\n", humanBytes(int64(s.AvgWriteBytes)))
	fmt.Fprintf(&b, "| Read amplification | %dx |\n", s.ReadAmplification)
	fmt.Fprintf(&b, "| Peak multiplier | %.1fx |\n", s.PeakMultiplier)
	fmt.Fprintf(&b, "| Replication | %dx |\n", s.ReplicationFactor)
	fmt.Fprintf(&b, "| Retention | %d days |\n\n", s.RetentionDays)

	fmt.Fprintln(&b, "### Derived")
	fmt.Fprintln(&b, "| Metric | Value |")
	fmt.Fprintln(&b, "|---|---|")
	fmt.Fprintf(&b, "| Writes/sec avg | %s |\n", humanFloat(r.WritesPerSecAvg))
	fmt.Fprintf(&b, "| Writes/sec peak | %s |\n", humanFloat(r.WritesPerSecPeak))
	fmt.Fprintf(&b, "| Reads/sec avg | %s |\n", humanFloat(r.ReadsPerSecAvg))
	fmt.Fprintf(&b, "| Reads/sec peak | %s |\n", humanFloat(r.ReadsPerSecPeak))
	fmt.Fprintf(&b, "| Storage/day | %s |\n", humanBytes(r.StoragePerDay))
	fmt.Fprintf(&b, "| Storage/year | %s |\n", humanBytes(r.StoragePerYear))
	fmt.Fprintf(&b, "| Bandwidth peak | %s/sec |\n", humanBytes(r.BandwidthReadPeak))

	return b.String()
}

func humanInt64(n int64) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

func humanFloat(f float64) string {
	switch {
	case f >= 1_000_000:
		return fmt.Sprintf("%.1fM/sec", f/1_000_000)
	case f >= 1_000:
		return fmt.Sprintf("%.1fK/sec", f/1_000)
	}
	return fmt.Sprintf("%.0f/sec", f)
}

func humanBytes(b int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
		pb = tb * 1024
	)
	switch {
	case b >= pb:
		return fmt.Sprintf("%.1f PB", float64(b)/pb)
	case b >= tb:
		return fmt.Sprintf("%.1f TB", float64(b)/tb)
	case b >= gb:
		return fmt.Sprintf("%.1f GB", float64(b)/gb)
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/mb)
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/kb)
	}
	return fmt.Sprintf("%d B", b)
}
