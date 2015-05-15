package humanize

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func validate(t *testing.T, v string, expect time.Time) {
	t1, err := ParseTime(v)
	assert.NoError(t, err)
	assert.Equal(t, expect, t1, "ParseTime(%s) should equal to %s, got: %s", v, expect.String(), t1.String())
}

func TestParseTime(t *testing.T) {
	validate(t, "2015-05-14T01:02:33Z", mustParseRFC3339("2015-05-14T01:02:33Z"))
	validate(t, "2015-05-14 01:02:33", mustParseRFC3339("2015-05-14T01:02:33Z"))
	validate(t, "2015-05-14", time.Date(2015, 5, 14, 0, 0, 0, 0, time.Local))

	now := time.Now().Local()
	validate(t, "01:02:03", time.Date(now.Year(), now.Month(), now.Day(), 1, 2, 3, 0, now.Location()))

	yesterday := time.Now().Add(-24 * time.Hour).Local()
	validate(t, "yesterday", truncateDay(yesterday))

	tomorrow := time.Now().Add(24 * time.Hour).Local()
	validate(t, "tomorrow", truncateDay(tomorrow))

	day1 := now.Add(-24 * time.Hour)
	validate(t, "1 day ago", truncateDay(day1))

	day2 := now.Add(-3 * 24 * time.Hour)
	validate(t, "3 days ago", truncateDay(day2))

	// this month
	validate(t, "this month", time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()))

	// this week
	d := -time.Duration(time.Now().Weekday()-time.Monday) * 24 * time.Hour
	day3 := now.Add(d)
	// fmt.Fprintln(os.Stderr, d, day3)
	validate(t, "this week", truncateDay(day3))

	// last month
	year := now.Year()
	month := now.Month() - 1
	if month < 1 {
		month = 1
		year -= 1
	}
	validate(t, "last month", time.Date(year, month, 1, 0, 0, 0, 0, now.Location()))

	// last week
	d1 := (7 + time.Duration(now.Weekday()-time.Monday)) * 24 * time.Hour
	day4 := now.Add(-d1)
	validate(t, "last week", time.Date(day4.Year(), day4.Month(), day4.Day(), 0, 0, 0, 0, now.Location()))

	// hours ago
	t1 := now.Add(-1 * time.Hour)
	validate(t, "1 hour ago", time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), 0, 0, 0, now.Location()))

	t2 := now.Add(-3 * time.Hour)
	validate(t, "3 hours ago", time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), 0, 0, 0, now.Location()))

}
