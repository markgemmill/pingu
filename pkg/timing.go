package pkg

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CalculatePauseInSeconds(retry, increment int) time.Duration {
	seconds := int64(math.Floor(math.Pow(float64(retry), float64(increment))))
	return time.Duration(seconds) * time.Second
}

/*
DayTime struct represents a time on a given day of the week.
*/
type DayTime struct {
	Dow    string
	Hour   int
	Minute int
	AmPm   string
}

func (d DayTime) String() string {
	return fmt.Sprintf("%s %d:%d", d.Dow, d.twentyFour(), d.Minute)
}

// twentyFour returns the 24-hour version of the hour.
func (d DayTime) twentyFour() int {
	if d.AmPm == "PM" {
		return d.Hour + 12
	}
	if d.AmPm == "AM" && d.Hour == 12 {
		return 0
	}
	return d.Hour
}

/*
priorToDate calculates the specific date the DayTime would be in relation
to the given Time. If DayTime specifies a day of the week, then it will
be either the current given day of the week or earlier. If DayTime does not
specify a day of the week, then it will be the same as the current given
day of the week.
*/
func (d DayTime) priorToDate(given time.Time) (time.Time, error) {
	var days = []int{0, -1, -2, -3, -4, -5, -6}
	for _, i := range days {
		cd := given.AddDate(0, 0, i)
		dow := d.Dow
		if dow == "" {
			dow = strings.ToUpper(cd.Weekday().String())
		}
		if strings.HasPrefix(strings.ToUpper(cd.Weekday().String()), dow) {
			return time.Date(cd.Year(), cd.Month(), cd.Day(), d.twentyFour(), d.Minute, 0, 0, cd.Location()), nil
		}
	}

	return given, errors.New("No prior date found!")
}

/*
afterDate calculates the specific date the DayTime would be in relation
to the given Time. If DayTime specifies a day of the week, then it will
be either the current given day of the week or later. If DayTime does not
specify a day of the week, then it will be the same as the current given
day of the week.
*/
func (d DayTime) afterDate(given time.Time) (time.Time, error) {
	var days = []int{0, 1, 2, 3, 4, 5, 6}
	for _, i := range days {
		cd := given.AddDate(0, 0, i)
		dow := d.Dow
		if dow == "" {
			dow = strings.ToUpper(cd.Weekday().String())
		}
		if strings.HasPrefix(strings.ToUpper(cd.Weekday().String()), dow) {
			return time.Date(cd.Year(), cd.Month(), cd.Day(), d.twentyFour(), d.Minute, 0, 0, cd.Location()), nil
		}
	}

	return given, errors.New("No after date found!")
}

func NewDayTime(day, hour, minute, amPm string) (DayTime, error) {

	hr, err := strconv.ParseInt(hour, 0, 64)
	if err != nil {
		return DayTime{}, errors.New(fmt.Sprintf("'%s' is invalid hour value.", hour))
	}
	mn, err := strconv.ParseInt(minute, 0, 64)
	if err != nil {
		return DayTime{}, errors.New(fmt.Sprintf("'%s' is invalid minute value.", minute))
	}

	dayTime := DayTime{
		Dow:    day,
		Hour:   int(hr),
		Minute: int(mn),
		AmPm:   amPm,
	}

	return dayTime, nil
}

/*
TimePeriod takes a start and end DayTime values and calculates the
specific start and end Time values from a given Time.
*/
type TimePeriod struct {
	starts    DayTime
	ends      DayTime
	startTime time.Time
	endTime   time.Time
	current   time.Time
}

func NewTimeout(starts, ends DayTime) *TimePeriod {
	t := TimePeriod{
		starts: starts,
		ends:   ends,
	}
	return &t
}

func (t *TimePeriod) Init(given time.Time) {
	t.current = given
	t.startTime, _ = t.starts.priorToDate(t.current)
	t.endTime, _ = t.ends.afterDate(t.startTime)
}

func (t *TimePeriod) String() string {
	isWithinTimePeriod := t.InTimePeriod()

	return fmt.Sprintf(
		"%s %v %s - %s",
		t.current.Format("2006-01-02 03:04PM"),
		isWithinTimePeriod,
		t.startTime.Format("2006-01-02 03:04PM"),
		t.endTime.Format("2006-01-02 03:04PM"),
	)
}

func (t *TimePeriod) InTimePeriod() bool {
	return (t.startTime.Equal(t.current) || t.startTime.Before(t.current)) &&
		(t.endTime.Equal(t.current) || t.endTime.After(t.current))
}

/*
ParseTimePeriod accepts a string in the format and returns DayTime objects
representing the start and ending of the time period:

	DOW HH:MM AM - DOW HH:MM AM

DOW can be of MON, TUE, WED, THU, FRI, SAT, SUN.

AM can be either AM or PM - if omitted, then expect the HH:MM to be 24 hr.

DOW can be omitted which would indicate the date would be daily.

Examples:

	SAT 23:00 - SUN 01:00
	SAT 11:00 PM - SUN 01:00 AM
*/
func ParseTimePeriod(timeoutStatement string) (DayTime, DayTime, error) {

	rx := `^(?P<beginDay>MON|TUE|WED|THU|FRI|SAT|SUN)? ?(?P<beginHr>\d{1,2}):(?P<beginMin>\d\d) ?(?P<beginAM>AM|PM)? ?- ?(?P<endDay>MON|TUE|WED|THU|FRI|SAT|SUN)? ?(?P<endHr>\d{1,2}):(?P<endMin>\d\d) ?(?P<endAM>AM|PM)?$`
	re := regexp.MustCompile(rx)

	match := re.FindStringSubmatch(timeoutStatement)
	if len(match) == 0 {
		return DayTime{}, DayTime{}, errors.New("TimePeriod string invalid.")
	}
	mapping := make(map[string]int)
	for i, name := range re.SubexpNames() {
		mapping[name] = i
	}

	beginning, err := NewDayTime(
		match[mapping["beginDay"]],
		match[mapping["beginHr"]],
		match[mapping["beginMin"]],
		match[mapping["beginAM"]],
	)
	if err != nil {
		return DayTime{}, DayTime{}, err
	}

	ending, err := NewDayTime(
		match[mapping["endDay"]],
		match[mapping["endHr"]],
		match[mapping["endMin"]],
		match[mapping["endAM"]],
	)
	if err != nil {
		return beginning, DayTime{}, err
	}

	return beginning, ending, nil
}

func IsIgnorePeriodActive(ignorePeriodStr string, currentTime time.Time) bool {
	b, e, _ := ParseTimePeriod(ignorePeriodStr)
	to := NewTimeout(b, e)
	to.Init(currentTime)
	return to.InTimePeriod()
}
