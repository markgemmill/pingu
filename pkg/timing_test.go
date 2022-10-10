package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculatePauseInSeconds(t *testing.T) {
	data := [][]int{
		[]int{1, 1, 1},
		[]int{2, 1, 2},
		[]int{1, 2, 1},
		[]int{2, 2, 4},
		[]int{1, 3, 1},
		[]int{2, 3, 8},
	}
	for _, d := range data {
		fmt.Printf("%d + %d -> %d\n", d[0], d[1], d[2])
		seconds := CalculatePauseInSeconds(d[0], d[1])
		assert.Equal(t, time.Duration(d[2])*time.Second, seconds)
	}
}
