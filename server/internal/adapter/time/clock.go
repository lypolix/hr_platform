package clock

import "time"

type utcClock struct{}

func NewUTCClock() *utcClock {
	return &utcClock{}
}

func (*utcClock) Now() time.Time {
	return time.Now().UTC()
}
