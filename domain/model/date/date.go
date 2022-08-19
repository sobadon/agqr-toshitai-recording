package date

import (
	"time"

	"github.com/sobadon/agqr-toshitai-recording/internal/timeutil"
)

// 年月日
type Date time.Time

func New(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, timeutil.LocationJST()))
}

func (d Date) Format(layout string) string {
	return time.Time(d).Format(layout)
}
