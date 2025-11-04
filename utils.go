package main

import (
	"strconv"
	"time"
)

type TimerSummary struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Mode     string // "timer" or "counter"
	Finished bool   // true if completed, false if quit/interrupted
	Name     string // optional name for the timer
}

func addSuffixIfArgIsNumber(s *string, suffix string) {
	_, err := strconv.ParseFloat(*s, 64)
	if err == nil {
		*s = *s + suffix
	}
}
