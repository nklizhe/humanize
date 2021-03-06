package humanize

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	regNow             = regexp.MustCompile("^\\s*now\\s*(\\(\\s*\\))?\\s*$")
	regNowWithDays     = regexp.MustCompile("^\\s*now\\s*(\\(\\s*\\))?\\s+([+-]{1})\\s+(\\d+)\\s*(d|day[s]?)")
	regNowWithDuration = regexp.MustCompile("^\\s*now\\s*(\\(\\s*\\))?\\s+([+-]{1})\\s+(.*)$")

	regDate        = regexp.MustCompile("(\\d{4})-(\\d{2})-(\\d{2})")
	regTime        = regexp.MustCompile("(\\d{2}):(\\d{2}):(\\d{2})")
	regThisMonth   = regexp.MustCompile("this\\s+month")
	regLastMonth   = regexp.MustCompile("last\\s+month")
	regNextMonth   = regexp.MustCompile("next\\s+month")
	regMonthsAgo   = regexp.MustCompile("(\\d+)\\s+(month[s]?)\\s+ago")
	regMonthsLater = regexp.MustCompile("(\\d+)\\s+(month[s]?)\\s+later")
	regThisWeek    = regexp.MustCompile("this\\s+week")
	regLastWeek    = regexp.MustCompile("last\\s+week")
	regDaysAgo     = regexp.MustCompile("(\\d+)\\s+(d|day[s]?)\\s+ago")
	regDaysLater   = regexp.MustCompile("(\\d+)\\s+(d|day[s]?)\\s+later")
	regYesterday   = regexp.MustCompile("yesterday")
	regToday       = regexp.MustCompile("today")
	regTomorrow    = regexp.MustCompile("tomorrow")
	regHoursAgo    = regexp.MustCompile("(\\d+)\\s+(h|hour[s]?)\\s+ago")
	regHoursLater  = regexp.MustCompile("(\\d+)\\s+(h|hour[s]?)\\s+later")
)

func truncateYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func truncateMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func truncateWeek(t time.Time) time.Time {
	days := t.Weekday() - time.Monday
	return truncateDay(t.Add(time.Duration(-days) * 24 * time.Hour))
}

func truncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func truncateHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// ParseTime parses a string into time
func ParseTime(s string) (time.Time, error) {
	// try "2006-01-02T15:04:05Z"
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}

	// try "2006-01-02 15:04:05"
	t, err = time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		return t, nil
	}

	// try "2006-01-02"
	if regDate.MatchString(s) {
		m := regDate.FindStringSubmatch(s)
		if len(m) > 3 {
			year, _ := strconv.ParseInt(m[1], 10, 16)
			month, _ := strconv.ParseInt(m[2], 10, 8)
			day, _ := strconv.ParseInt(m[3], 10, 8)
			today := time.Now().Local()
			return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, today.Location()), nil
		}
	}

	// try "15:04:05"
	if regTime.MatchString(s) {
		m := regTime.FindStringSubmatch(s)
		if len(m) > 3 {
			hour, _ := strconv.ParseInt(m[1], 10, 8)
			minute, _ := strconv.ParseInt(m[2], 10, 8)
			second, _ := strconv.ParseInt(m[3], 10, 8)
			today := time.Now().Local().Truncate(24 * time.Hour)
			return time.Date(today.Year(), today.Month(), today.Day(), int(hour), int(minute), int(second), 0, today.Location()), nil
		}
	}

	// try this month
	if regThisMonth.MatchString(s) {
		return truncateMonth(time.Now()), nil
	}

	// try last month
	if regLastMonth.MatchString(s) {
		today := time.Now().Local()
		y2 := today.Year()
		m2 := today.Month() - 1
		for m2 <= 0 {
			m2 += 12
			y2 -= 1
		}
		return time.Date(y2, m2, 1, 0, 0, 0, 0, today.Location()), nil
	}

	// try next month
	if regNextMonth.MatchString(s) {
		today := time.Now().Local()
		y2 := today.Year()
		m2 := today.Month() + 1
		for m2 > 12 {
			m2 -= 12
			y2 += 1
		}
		return time.Date(y2, m2, 1, 0, 0, 0, 0, today.Location()), nil
	}

	// try "n months ago"
	if regMonthsAgo.MatchString(s) {
		m := regMonthsAgo.FindStringSubmatch(s)
		if len(m) > 2 {
			months, _ := strconv.ParseInt(m[1], 10, 16)
			today := time.Now().Local()
			y2 := today.Year()
			m2 := int64(today.Month()) - months
			for m2 <= 0 {
				m2 += 12
				y2 -= 1
			}
			return time.Date(y2, time.Month(m2), 1, 0, 0, 0, 0, today.Location()), nil
		}
	}

	// try "n months later"
	if regMonthsAgo.MatchString(s) {
		m := regMonthsAgo.FindStringSubmatch(s)
		if len(m) > 2 {
			months, _ := strconv.ParseInt(m[1], 10, 16)
			today := time.Now().Local()
			y2 := today.Year()
			m2 := int64(today.Month()) + months
			for m2 > 0 {
				m2 -= 12
				y2 += 1
			}
			return time.Date(y2, time.Month(m2), 1, 0, 0, 0, 0, today.Location()), nil
		}
	}

	// try "this week"
	if regThisWeek.MatchString(s) {
		return truncateWeek(time.Now()), nil
	}

	// try "last week"
	if regLastWeek.MatchString(s) {
		return truncateWeek(time.Now().Add(-7 * 24 * time.Hour)), nil
	}

	// try "n days ago"
	if regDaysAgo.MatchString(s) {
		m := regDaysAgo.FindStringSubmatch(s)
		if len(m) > 2 {
			days, _ := strconv.ParseInt(m[1], 10, 16)
			dur := -1 * 24 * time.Duration(days) * time.Hour
			today := time.Now().Add(dur).Local()
			return time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()), nil
		}
	}

	// try "yesterday"
	if regYesterday.MatchString(s) {
		return truncateDay(time.Now().Add(-24 * time.Hour).Local()), nil
	}

	// try "today"
	if regToday.MatchString(s) {
		return truncateDay(time.Now().Local()), nil
	}

	// try "tomorrow"
	if regTomorrow.MatchString(s) {
		return truncateDay(time.Now().Add(24 * time.Hour).Local()), nil
	}

	// try "n days later"
	if regDaysLater.MatchString(s) {
		m := regDaysLater.FindStringSubmatch(s)
		if len(m) > 2 {
			days, _ := strconv.ParseInt(m[1], 10, 16)
			dur := 1 * 24 * time.Duration(days) * time.Hour
			today := time.Now().Add(dur).Local()
			return time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()), nil
		}
	}

	// try "n hours ago"
	if regHoursAgo.MatchString(s) {
		m := regHoursAgo.FindStringSubmatch(s)
		if len(m) > 2 {
			hour, _ := strconv.ParseInt(m[1], 10, 16)
			dur := -1 * time.Duration(hour) * time.Hour
			today := time.Now().Add(dur).Local().Truncate(time.Hour)
			return today, nil
		}
	}

	// try "n hours later"
	if regHoursLater.MatchString(s) {
		m := regHoursLater.FindStringSubmatch(s)
		if len(m) > 2 {
			hour, _ := strconv.ParseInt(m[1], 10, 16)
			dur := time.Duration(hour) * time.Hour
			today := time.Now().Add(dur).Local().Truncate(time.Hour)
			return today, nil
		}
	}

	// try "now"
	if regNow.MatchString(s) {
		return time.Now(), nil
	}

	// try "now +/- days"
	if regNowWithDays.MatchString(s) {
		m := regNowWithDays.FindStringSubmatch(s)
		if len(m) > 3 {
			sign := m[2]
			days, err := strconv.ParseInt(m[3], 10, 16)
			if err != nil {
				return time.Unix(0, 0), err
			}
			dur := time.Duration(days*24) * time.Hour
			if sign == "-" {
				dur = -1 * dur
			}
			return time.Now().Add(dur), nil
		}
	}

	// try "now +/- duration"
	if regNowWithDuration.MatchString(s) {
		m := regNowWithDuration.FindStringSubmatch(s)
		// fmt.Fprintln(os.Stderr, "regNowWithDuration", s, m)
		if len(m) > 3 {
			sign := m[2]
			durStr := m[3]
			dur, err := time.ParseDuration(durStr)
			if err != nil {
				return time.Unix(0, 0), err
			}
			if sign == "-" {
				dur = -1 * dur
			}
			return time.Now().Add(dur), nil
		}
	}

	// failed
	return time.Unix(0, 0), fmt.Errorf("cannot parse \"%s\"", s)
}
