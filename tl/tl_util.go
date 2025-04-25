package tl

import (
	"errors"
	"time"
)

func isWeekday() bool {
	weekday := time.Now().Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

func generateDayTimes() (string, string) {
	now := time.Now()
	// Create morning 10:00 AM and evening 6:00 PM times
	morning := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	evening := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	morningStr := morning.Format("2006-01-02T15:04:05")
	eveningStr := evening.Format("2006-01-02T15:04:05")
	return morningStr, eveningStr
}

func loggedToday(date string) (bool, error) {
	t, err := time.Parse(Layout, date)
	if err != nil {
		return false, errors.New("Error parsing timestamp")
	}
	d := t.Format(Format)
	t1 := time.Now()
	d1 := t1.Format(Format)
	return d == d1, nil
}

func isAfterSixPMInIST() bool {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(ist)
	sixPM := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, ist)
	return now.After(sixPM)
}
