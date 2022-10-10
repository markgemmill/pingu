package pkg

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestTimeSegment(t *testing.T) {
	d, _ := time.ParseDuration("1446m")
	hours, minutes := TimeSegment(d.Minutes(), 60)
	assert.Equal(t, 24.0, hours)
	assert.Equal(t, 6.0, minutes)

	days, h := TimeSegment(hours, 24)

	assert.Equal(t, 1.0, days)
	assert.Equal(t, 0.0, h)

	d, _ = time.ParseDuration("83h4m")
	assert.Equal(t, 4984.0, d.Minutes())
	assert.Equal(t, "83h4m0s", d.String())

	hours, minutes = TimeSegment(d.Minutes(), 60)
	assert.Equal(t, 83.0, hours)
	assert.Equal(t, 4.0, minutes)

	days, h = TimeSegment(hours, 24)
	assert.Equal(t, 3.0, days)
	assert.Equal(t, 11.0, h)
}

func TestDurationString(t *testing.T) {

	d, _ := time.ParseDuration("83h4m")
	assert.Equal(t, "3 days, 11 hours and 4 minutes", DurationString(d))

	d, _ = time.ParseDuration("13h55m")
	assert.Equal(t, "13 hours and 55 minutes", DurationString(d))

	d, _ = time.ParseDuration("60m")
	assert.Equal(t, "1 hour", DurationString(d))

	d, _ = time.ParseDuration("1440m")
	assert.Equal(t, "1 day", DurationString(d))

	d, _ = time.ParseDuration("1446m")
	assert.Equal(t, "1 day and 6 minutes", DurationString(d))
}

func TestStoreRecordReport(t *testing.T) {
	record := StoreRecord{
		Start:  time.Date(2022, 9, 1, 10, 23, 28, 0, time.Local),
		Last:   time.Date(2022, 9, 21, 23, 12, 9, 0, time.Local),
		Status: "PASS",
		Count:  456,
	}

	d := record.Last.Sub(record.Start)
	assert.Equal(t, 29568.0, math.Floor(d.Minutes()))

	assert.Equal(
		t,
		"2022-09-21 23:12:09 PASSING 456 checks for 20 days, 12 hours and 48 minutes.",
		StoreRecordStatusReport(&record),
	)
}
