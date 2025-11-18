package analytics

import (
	"time"
)

type AnalyticsRecord struct {
	Timestamp  string    `json:"timestamp"`
	Username   string    `json:"username"`
	Request    string    `json:"request"`
	Pool       string    `json:"pool"`
	Deciders   string    `json:"deciders"`
	Conclusion string    `json:"conclusion"`
	Effect     string    `json:"effect"`
	ExpireAt   time.Time `json:"expireAt" bson:"expireAt"`
}

type AnalyticsFilters struct {
	SkippedUsernames []string `json:"skippingUsername" mapstructure:"skippingUsername"`
}

func (f *AnalyticsFilters) HasFilter() bool {
	return len(f.SkippedUsernames) > 0
}

func (f *AnalyticsFilters) ShouldFilter(record *AnalyticsRecord) bool {
	return inSlice(record.Username, f.SkippedUsernames)
}

func inSlice(s string, ss []string) bool {
	for _, item := range ss {
		if item == s {
			return true
		}
	}

	return false
}
