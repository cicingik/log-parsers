package utils

import "time"

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func InBound(t, startTime, endTime time.Time) bool {
	if t.Equal(startTime) {
		return true
	}
	if t.After(startTime) && t.Before(endTime) {
		return true
	}

	return false
}

// DateInBound is force check in bound with exclude time
func DateInBound(t, startTime, endTime time.Time) bool {
	startDate := truncateToDay(startTime)
	endDate := truncateToDay(endTime).AddDate(0, 0, 1)

	if t.Equal(startDate) {
		return true
	}

	if t.After(startDate) && t.Before(endDate) {
		return true
	}

	return false
}
