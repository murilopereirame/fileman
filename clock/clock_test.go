package clock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRealClockUnixReturnsCurrentUnixTime(t *testing.T) {
	realClock := RealClock{}
	now := time.Now().Unix()
	unixTime := realClock.Unix()

	assert.True(t, unixTime == now)
}

func TestCalculateAgeReturnsZeroForCurrentTime(t *testing.T) {
	realClock := RealClock{}
	now := time.Now().Unix()
	age := realClock.CalculateAge(now)

	assert.Zero(t, age)
}

func TestCalculateAgeReturnsPositiveForPastReference(t *testing.T) {
	realClock := RealClock{}
	past := time.Now().Add(-48 * time.Hour).Unix()
	age := realClock.CalculateAge(past)

	assert.Positive(t, age, "Expected age to be positive, got %f", age)
}

func TestCalculateAgeReturnsNegativeForFutureReference(t *testing.T) {
	realClock := RealClock{}
	future := time.Now().Add(48 * time.Hour).Unix()
	age := realClock.CalculateAge(future)

	assert.Negative(t, age, "Expected age to be negative, got %f", age)
}
