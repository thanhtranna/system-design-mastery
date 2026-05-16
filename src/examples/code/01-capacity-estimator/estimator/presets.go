package estimator

// Presets contains canonical scenarios used as starting points for analysis.
var Presets = map[string]Scenario{
	"twitter": {
		Name:              "Twitter / X",
		DAU:               200_000_000,
		Actions:           Actions{WritesPerUserPerDay: 2, ReadsPerUserPerDay: 50},
		AvgWriteBytes:     280,
		ReadAmplification: 2,
		PeakMultiplier:    3,
		ReplicationFactor: 3,
		RetentionDays:     365 * 5,
	},
	"instagram": {
		Name:              "Instagram",
		DAU:               500_000_000,
		Actions:           Actions{WritesPerUserPerDay: 1, ReadsPerUserPerDay: 60},
		AvgWriteBytes:     500_000, // ~500KB per photo upload
		ReadAmplification: 1,
		PeakMultiplier:    2,
		ReplicationFactor: 3,
		RetentionDays:     365 * 10,
	},
	"whatsapp": {
		Name:              "WhatsApp",
		DAU:               1_500_000_000,
		Actions:           Actions{WritesPerUserPerDay: 50, ReadsPerUserPerDay: 50},
		AvgWriteBytes:     200,
		ReadAmplification: 1,
		PeakMultiplier:    3,
		ReplicationFactor: 3,
		RetentionDays:     365 * 2,
	},
	"uber": {
		Name:              "Uber",
		DAU:               20_000_000,
		Actions:           Actions{WritesPerUserPerDay: 5, ReadsPerUserPerDay: 200}, // GPS pings
		AvgWriteBytes:     50,
		ReadAmplification: 1,
		PeakMultiplier:    5, // friday night, NYE
		ReplicationFactor: 3,
		RetentionDays:     365 * 1,
	},
	"propertyhub": {
		Name:              "PropertyHub (real estate platform)",
		DAU:               50_000,
		Actions:           Actions{WritesPerUserPerDay: 5, ReadsPerUserPerDay: 100},
		AvgWriteBytes:     5_000,
		ReadAmplification: 1,
		PeakMultiplier:    3,
		ReplicationFactor: 3,
		RetentionDays:     365 * 3,
	},
}
