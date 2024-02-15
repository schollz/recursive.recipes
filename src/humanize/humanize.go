package humanize

import (
	"fmt"
	"strings"
	"time"

	"github.com/hako/durafmt"
)

// Unit represents duration unit measurements in milliseconds
var UnitMeasures = map[string]int64{
	"month":  2629800000,
	"week":   604800000,
	"day":    86400000,
	"hour":   3600000,
	"minute": 60000,
}

func Duration(ms int64) string {
	if ms < 60000 {
		return "Less than a minute"
	}
	if ms > 2629800000*100 {
		return "a very long time"
	}

	timeduration := time.Duration(ms) * time.Millisecond
	fmt.Println(ms, timeduration)
	return durafmt.Parse(timeduration).LimitFirstN(2).String()
}

// HumanizeDuration converts duration in milliseconds to a human-readable format
func Duration2(ms int64) string {
	if ms < 0 {
		ms = -ms
	}
	if ms < 60000 {
		return "Less than a minute"
	}
	var parts []string
	for unit, value := range UnitMeasures {
		if ms >= value {
			count := ms / value
			ms %= value
			part := fmt.Sprintf("%d %s", count, pluralize(unit, count))
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, ", ")
}

// pluralize adds "s" to unit if count is not 1
func pluralize(unit string, count int64) string {
	if count == 1 {
		return unit
	}
	return unit + "s"
}
