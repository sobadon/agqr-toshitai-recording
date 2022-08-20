package program

import "time"

type Program struct {
	ID    int
	Title string
	Start time.Time
	End   time.Time

	// すぐは必要ないならあとで
	// Personality []string
}
