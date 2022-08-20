package program

import "time"

type Program struct {
	ID     int
	Title  string
	Start  time.Time
	End    time.Time
	Status Status

	// すぐは必要ないならあとで
	// Personality []string
}

type Status string

const (
	StatusScheduled = Status("scheduled")
	StatusRecording = Status("scheduled")
	StatusDone      = Status("done")
	StatusFailed    = Status("failed")
)

func (s Status) String() string {
	return string(s)
}
