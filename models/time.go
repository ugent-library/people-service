package models

import "time"

var BeginningOfTime = time.Date(
	0, time.January, 1, 0, 0, 0, 0, time.UTC,
)
var EndOfTime = time.Date(
	9999, time.December, 31, 23, 59, 59, 999, time.UTC,
)

func IsMaybeEndOfTime(t time.Time) bool {
	return t.Year() >= 9999
}

func NormalizeEndOfTime(t time.Time) *time.Time {
	if IsMaybeEndOfTime(t) {
		return nil
	}
	return &t
}

func copyTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	t2 := *t
	return &t2
}
