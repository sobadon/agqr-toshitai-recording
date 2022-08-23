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

func Dummies(now time.Time) []Program {
	return []Program{
		{
			ID:    1,
			Title: "ダミー1",
			Start: now.Add(20 * time.Second),
			End:   now.Add(20 * time.Second).Add(1 * time.Minute),
		},
		// {
		// 	ID: 2,
		// 	Title: "ダミー2",
		// 	Start: now.Add(20 * time.Second),
		// 	End: now.Add(20 * time.Second).Add(1 * time.Minute),
		// },
	}
}

type Status string

const (
	StatusScheduled = Status("scheduled")
	StatusRecording = Status("recording")
	StatusDone      = Status("done")
	StatusFailed    = Status("failed")
)

func (s Status) String() string {
	return string(s)
}
