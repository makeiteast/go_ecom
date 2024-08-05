package util

import "time"

const TimeFormat = "2006-01-02 15:04:05"
const TimeFormatRSS = "Mon, 02 Jan 2006 15:04:05 -0700"

// ParseDate s - string in DD-MM-YYYY format
func ParseDate(s string) (r time.Time, err error) {
	r, err = time.Parse("02-01-2006", s)
	if err == nil {
		r.In(time.UTC)
	}
	return
}

// ParseDatetimeRss s - string in YYYY-MM-DDT00:00:00Z format
func ParseDatetimeRss(s string) (r time.Time, err error) {
	r, err = time.Parse(TimeFormatRSS, s)
	if err == nil {
		r.In(time.UTC)
	}
	return
}

func NowUTC() time.Time {
	return time.Now().Round(time.Microsecond).UTC()
}

func FormatToDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

func LogDateTimeToStringLocalTime(t time.Time) string {
	return t.Add(5 * time.Hour).Format("2006-01-02 15:04:05.000000")
}

func TimeFromInt64(ts int64) time.Time {
	return time.Unix(ts, 0).UTC()
}
