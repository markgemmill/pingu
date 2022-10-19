package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculatePauseInSeconds(t *testing.T) {
	data := [][]int{
		{1, 1, 1},
		{2, 1, 2},
		{1, 2, 1},
		{2, 2, 4},
		{1, 3, 1},
		{2, 3, 8},
	}
	for _, d := range data {
		fmt.Printf("%d + %d -> %d\n", d[0], d[1], d[2])
		seconds := CalculatePauseInSeconds(d[0], d[1])
		assert.Equal(t, time.Duration(d[2])*time.Second, seconds)
	}
}

func assertDayTime(t *testing.T, d DayTime, expectedDOW string, expectedAmPm string, expectedHour, expectedMinute int) {
	assert.Equal(t, expectedDOW, d.Dow)
	assert.Equal(t, expectedAmPm, d.AmPm)
	assert.Equal(t, expectedHour, d.Hour)
	assert.Equal(t, expectedMinute, d.Minute)
}

func TestParseTimeout(t *testing.T) {
	b, e, _ := ParseTimePeriod("MON 10:30 AM - TUE 12:00 AM")
	assertDayTime(t, b, "MON", "AM", 10, 30)
	assertDayTime(t, e, "TUE", "AM", 12, 0)

	b, e, _ = ParseTimePeriod("MON 10:30AM-TUE 12:00AM")
	assertDayTime(t, b, "MON", "AM", 10, 30)
	assertDayTime(t, e, "TUE", "AM", 12, 0)

	// NO SECOND DAY
	b, e, _ = ParseTimePeriod("MON 10:30 AM - 12:00 PM")
	assertDayTime(t, b, "MON", "AM", 10, 30)
	assertDayTime(t, e, "", "PM", 12, 0)

	b, e, _ = ParseTimePeriod("MON 10:30AM-12:00PM")
	assertDayTime(t, b, "MON", "AM", 10, 30)
	assertDayTime(t, e, "", "PM", 12, 0)

	// NO DAY OF WEEK
	b, e, _ = ParseTimePeriod("10:30 AM - 12:00 PM")
	assertDayTime(t, b, "", "AM", 10, 30)
	assertDayTime(t, e, "", "PM", 12, 0)

	b, e, _ = ParseTimePeriod("10:30AM-12:00PM")
	assertDayTime(t, b, "", "AM", 10, 30)
	assertDayTime(t, e, "", "PM", 12, 0)

	// NO AM/PM
	b, e, _ = ParseTimePeriod("10:30 - 12:00")
	assertDayTime(t, b, "", "", 10, 30)
	assertDayTime(t, e, "", "", 12, 0)

	b, e, _ = ParseTimePeriod("10:30-12:00")
	assertDayTime(t, b, "", "", 10, 30)
	assertDayTime(t, e, "", "", 12, 0)

	b, e, _ = ParseTimePeriod("TUE 10:00 PM - WED 12:00 AM")
	assertDayTime(t, b, "TUE", "PM", 10, 0)
	assertDayTime(t, e, "WED", "AM", 12, 0)

	b, e, _ = ParseTimePeriod("TUE 10:00 PM - WED 01:00 AM")
	assertDayTime(t, b, "TUE", "PM", 10, 0)
	assertDayTime(t, e, "WED", "AM", 1, 0)
}

func TestDayTimePriorTo(t *testing.T) {
	dt := DayTime{
		Dow:    "MON",
		AmPm:   "PM",
		Hour:   10,
		Minute: 0,
	}

	given := time.Date(2022, 10, 17, 12, 12, 0, 0, time.Now().Location())
	prior, _ := dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 17, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "TUE"
	prior, _ = dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 11, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "WED"
	prior, _ = dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 12, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "THU"
	prior, _ = dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 13, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "FRI"
	prior, _ = dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 14, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "SAT"
	prior, _ = dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 15, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "SUN"
	prior, _ = dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 16, 22, 0, 0, 0, time.Now().Location()), prior)

}

func TestDayTimePriorTo_Special_Cases(t *testing.T) {
	dt := DayTime{
		Dow:    "TUE",
		AmPm:   "PM",
		Hour:   10,
		Minute: 0,
	}

	given := time.Date(2022, 10, 18, 10, 0, 0, 0, time.Now().Location())
	prior, _ := dt.priorToDate(given)
	assert.Equal(t, time.Date(2022, 10, 18, 22, 0, 0, 0, time.Now().Location()), prior)
}

func TestDayTimeAfterDate(t *testing.T) {
	dt := DayTime{
		Dow:    "MON",
		AmPm:   "PM",
		Hour:   10,
		Minute: 0,
	}

	given := time.Date(2022, 10, 17, 12, 12, 0, 0, time.Now().Location())
	prior, _ := dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 17, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "TUE"
	prior, _ = dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 18, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "WED"
	prior, _ = dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 19, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "THU"
	prior, _ = dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 20, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "FRI"
	prior, _ = dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 21, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "SAT"
	prior, _ = dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 22, 22, 0, 0, 0, time.Now().Location()), prior)

	dt.Dow = "SUN"
	prior, _ = dt.afterDate(given)
	assert.Equal(t, time.Date(2022, 10, 23, 22, 0, 0, 0, time.Now().Location()), prior)
}

func assertInTimeout(t *testing.T, timeout string, day, hour, minute int, expectedResult bool) {
	b, e, _ := ParseTimePeriod(timeout)

	to := NewTimeout(b, e)
	loc := time.Now().Location()
	givenTime := time.Date(2022, 10, day, hour, minute, 0, 0, loc)
	to.Init(givenTime)
	fmt.Println(to.String())
	fmt.Printf("  > %s\n", to.starts.String())
	fmt.Printf("  > %s\n", to.startTime)
	fmt.Printf("  > %s\n", to.endTime)
	result := to.InTimePeriod()
	assert.Equal(t, expectedResult, result, fmt.Sprintf("Expected %s in '%s' to be %v, but got %v\n", givenTime.Format("2006-01-02 03:04PM"), timeout, expectedResult, result))
}

func TestTimeout(t *testing.T) {

	assertInTimeout(t, "10:00 AM - 11:00 AM", 18, 10, 0, true)
	assertInTimeout(t, "10:00 AM - 11:00 AM", 18, 10, 30, true)
	assertInTimeout(t, "10:00 AM - 11:00 AM", 18, 10, 59, true)
	assertInTimeout(t, "10:00 AM - 11:00 AM", 18, 11, 0, true)

	assertInTimeout(t, "10:00 AM - 11:00 AM", 18, 9, 59, false)
	assertInTimeout(t, "10:00 AM - 11:00 AM", 18, 11, 01, false)

	assertInTimeout(t, "TUE 10:00 PM - WED 1:00 AM", 18, 22, 0, true)
	assertInTimeout(t, "TUE 10:00PM - WED 1:00AM", 19, 0, 0, true)

	assertInTimeout(t, "TUE 10:00PM - WED 1:00AM", 18, 9, 59, false)
	assertInTimeout(t, "TUE 22:00 - WED 1:00", 18, 9, 59, false)

	assertInTimeout(t, "22:00 - 23:59", 18, 9, 59, false)
	assertInTimeout(t, "22:00 - 23:59", 18, 22, 0, true)
	assertInTimeout(t, "22:00 - 23:59", 18, 23, 0, true)
	assertInTimeout(t, "22:00 - 23:59", 19, 0, 0, false)

	assertInTimeout(t, "SAT 10:00PM - SUN 1:00AM", 15, 9, 59, false)
	assertInTimeout(t, "SAT 10:00PM - SUN 1:00AM", 15, 22, 0, true)
	assertInTimeout(t, "SAT 10:00PM - SUN 1:00AM", 15, 23, 59, true)
	assertInTimeout(t, "SAT 10:00PM - SUN 1:00AM", 16, 0, 0, true)
	assertInTimeout(t, "SAT 10:00PM - SUN 1:00AM", 16, 1, 0, true)
	assertInTimeout(t, "SAT 10:00PM - SUN 1:00AM", 16, 1, 1, false)

}
