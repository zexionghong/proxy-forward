package coarsetime

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	coarsetime atomic.Value
	frequency  = time.Millisecond * 100
	loc        = map[int]*time.Location{}
)

func init() {
	loc[0], _ = time.LoadLocation("Etc/GMT+0")
	loc[-1], _ = time.LoadLocation("Etc/GMT+1")
	loc[-2], _ = time.LoadLocation("Etc/GMT+2")
	loc[-3], _ = time.LoadLocation("Etc/GMT+3")
	loc[-4], _ = time.LoadLocation("Etc/GMT+4")
	loc[-5], _ = time.LoadLocation("Etc/GMT+5")
	loc[-6], _ = time.LoadLocation("Etc/GMT+6")
	loc[-7], _ = time.LoadLocation("Etc/GMT+7")
	loc[-8], _ = time.LoadLocation("Etc/GMT+8")
	loc[-9], _ = time.LoadLocation("Etc/GMT+9")
	loc[-10], _ = time.LoadLocation("Etc/GMT+10")
	loc[-11], _ = time.LoadLocation("Etc/GMT+11")
	loc[-12], _ = time.LoadLocation("Etc/GMT+12")
	loc[1], _ = time.LoadLocation("Etc/GMT-1")
	loc[2], _ = time.LoadLocation("Etc/GMT-2")
	loc[3], _ = time.LoadLocation("Etc/GMT-3")
	loc[4], _ = time.LoadLocation("Etc/GMT-4")
	loc[5], _ = time.LoadLocation("Etc/GMT-5")
	loc[6], _ = time.LoadLocation("Etc/GMT-6")
	loc[7], _ = time.LoadLocation("Etc/GMT-7")
	loc[8], _ = time.LoadLocation("Etc/GMT-8")
	loc[9], _ = time.LoadLocation("Etc/GMT-9")
	loc[10], _ = time.LoadLocation("Etc/GMT-10")
	loc[11], _ = time.LoadLocation("Etc/GMT-11")
	loc[12], _ = time.LoadLocation("Etc/GMT-12")
	t := time.Now().Truncate(frequency)
	coarsetime.Store(&t)
	go func() {
		for {
			time.Sleep(frequency)
			t := time.Now().Truncate(frequency)
			coarsetime.Store(&t)
		}
	}()
}

// FloorTimeNow returns the current time from the range (now - 100ms, now],
// This is a faster alternative to time.Now().
func FloorTimeNow() time.Time {
	tp := coarsetime.Load().(*time.Time)
	return (*tp)
}

// CeilingTimeNow returns the current time from the range [now, now+100ms).
// This is a faster alternative to time.Now()
func CeilingTimeNow() time.Time {
	tp := coarsetime.Load().(*time.Time)
	return (*tp).Add(frequency)
}

// CeilingDate returns the current date
// this is a faster alternative to time.Date(time_now.Year(), time_now.Month(), time_now.Day(), 0, 0, 0, 0, loc)
func CeilingTimezoneDateToday(timezone int) time.Time {
	l := loc[timezone]
	time_now := CeilingTimeNow().In(l)
	today := time.Date(time_now.Year(), time_now.Month(), time_now.Day(), 0, 0, 0, 0, l)
	return today
}

func CeilingTimezoneTimeNow(timezone int) time.Time {
	l := loc[timezone]
	time_now := CeilingTimeNow().In(l)
	return time_now
}

func CeilingTimezoneTimeNowYYMMDD(timezone int) int {
	l := loc[timezone]
	time_now := CeilingTimeNow().In(l)
	year := ""
	month := ""
	day := ""
	year = fmt.Sprintf("%d", time_now.Year())
	if time_now.Month() < 10 {
		month = fmt.Sprintf("0%d", time_now.Month())
	} else {
		month = fmt.Sprintf("%d", time_now.Month())
	}
	if time_now.Day() < 10 {
		day = fmt.Sprintf("0%d", time_now.Day())
	} else {
		day = fmt.Sprintf("%d", time_now.Day())
	}
	today := fmt.Sprintf("%s%s%s", year, month, day)
	todayInt, _ := strconv.Atoi(today)
	return todayInt
}

func CeilingTimezone(timezone time.Time) time.Time {
	l := loc[8]
	time_now := CeilingTimeNow().In(l)
	// 可以刷新
	if (time_now.Hour() > timezone.Hour()) || (time_now.Hour() == timezone.Hour() && time_now.Minute() > timezone.Minute()) {
		today := time.Date(time_now.Year(), time_now.Month(), time_now.Day(), timezone.Hour(), timezone.Minute(), 0, 0, l)
		return today
	}
	time_now = time_now.Add(time.Duration(-24 * time.Hour))
	today := time.Date(time_now.Year(), time_now.Month(), time_now.Day(), timezone.Hour(), timezone.Minute(), 0, 0, l)
	return today
}
