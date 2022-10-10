package pkg

import (
	"math"
	"time"
)

func CalculatePauseInSeconds(retry, increment int) time.Duration {
	seconds := int64(math.Floor(math.Pow(float64(retry), float64(increment))))
	return time.Duration(seconds) * time.Second
}
