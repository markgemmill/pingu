package pkg

import (
	"bytes"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"math"
	"sort"
	"strings"
	"text/template"
	"time"
)

func TimeSegment(d float64, f float64) (float64, float64) {
	wholeAmount := math.Floor(d / f)
	return wholeAmount, d - (wholeAmount * f)
}

func plural(f int) string {
	if f > 1 {
		return "s"
	}
	return ""
}

// DurationString returns a human-readable time duration.
// Example: 23 days, 5 hours and 3 minutes
func DurationString(d time.Duration) string {

	totalHours, minutes := TimeSegment(d.Minutes(), 60)
	days, hours := TimeSegment(totalHours, 24)

	b := strings.Builder{}
	if days > 0.0 {
		_, err := fmt.Fprintf(&b, "%d day%s", int(days), plural(int(days)))
		ExitOnError(err, "")
	}
	if days > 0.0 && hours > 0.0 {
		_, err := fmt.Fprintf(&b, ", %d hour%s", int(hours), plural(int(hours)))
		ExitOnError(err, "")
	} else if hours > 0.0 {
		_, err := fmt.Fprintf(&b, "%d hour%s", int(hours), plural(int(hours)))
		ExitOnError(err, "")
	}

	if (hours > 0.0 || days > 0.0) && minutes > 0 {
		_, err := fmt.Fprintf(&b, " and %d minute%s", int(minutes), plural(int(minutes)))
		ExitOnError(err, "")
	} else if minutes > 0 {
		_, err := fmt.Fprintf(&b, "%v minute%s", int(minutes), plural(int(minutes)))
		ExitOnError(err, "")
	}

	return b.String()
}

func StoreRecordStatusReport(record *StoreRecord) string {
	/*
		2020-09-24 10:34:00  PASSING 245 checks for last 23 days 12 hours and 5 minutes.
		2020-09-01 12:13:00  FAIL     10 checks for 1 hour and 23 minutes.
	*/
	const report = `{{ .Record.Last.Format "2006-01-02 15:04:05" }} {{ .Status }} {{ .Record.Count }} checks for {{ .Duration }}.`
	tmpl := template.Must(template.New("record-status").Parse(report))

	status := "FAILING"
	if record.Status == PASS {
		status = "PASSING"
	}

	d := record.Last.Sub(record.Start)

	data := struct {
		Record   *StoreRecord
		Status   string
		Duration string
	}{Record: record, Status: status, Duration: DurationString(d)}

	var b bytes.Buffer
	err := tmpl.Execute(&b, data)
	PanicOnError(err)

	return b.String()
}

// ReportMessage creates an html and text report of the data store.
type ReportMessage struct {
	Store   *StoreMaster
	context pongo2.Context
}

func (r *ReportMessage) Initialize() {

	current := StoreRecordStatusReport(&r.Store.Current)

	history := make([]StoreRecord, 0)
	history = append(history, r.Store.Passes...)
	history = append(history, r.Store.Failures...)

	sort.Slice(history, func(i, j int) bool {
		return history[i].Last.Before(history[j].Last)
	})

	var hist []string

	for _, record := range history {
		hist = append(hist, StoreRecordStatusReport(&record))
	}

	r.context = pongo2.Context{
		"url":     r.Store.Url,
		"current": current,
		"history": hist,
	}
}

func (r *ReportMessage) Subject() string {
	return fmt.Sprintf("URL CHECK REPORT: %s", r.Store.Url)
}

func (r *ReportMessage) ToHtml() string {
	return RenderTemplate("report-email.html", &r.context, true)
}

func (r *ReportMessage) ToText() string {
	return RenderTemplate("report-email.txt", &r.context, false)
}
