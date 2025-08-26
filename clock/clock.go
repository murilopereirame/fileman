package clock

import "time"

const DAY_IN_SECOND = 24 * 60 * 60

type Clock interface {
	Unix() int64
	CalculateAge(reference int64) float64
}

type RealClock struct{}

func (r RealClock) Unix() int64 {
	return time.Now().Unix()
}

// CalculateAge Takes a reference date and returns the difference
// between now and the given date in days
func (r RealClock) CalculateAge(reference int64) float64 {
	diff := float64(r.Unix() - reference)
	return diff / DAY_IN_SECOND
}
