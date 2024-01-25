package core

import "time"

func ParseDate(date string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04", date, zoneUTC)
	if err == nil {
		return t.In(zoneUTC)
	}
	return time.Now().In(zoneUTC)
}

func SetZoneOffset(offset int) {
	zone = time.FixedZone("CST", offset*3600)
}

func DateFormat(t time.Time, layout string) string {
	return t.In(zone).Format(layout)
}

func Now() time.Time {
	return time.Now().In(zone)
}
