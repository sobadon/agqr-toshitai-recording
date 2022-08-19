package program

import "time"

type Program struct {
	Title string
	Start time.Time
	End   time.Time

	// すぐは必要ないならあとで
	// Personality []string
}
