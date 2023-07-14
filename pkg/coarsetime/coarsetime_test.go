package coarsetime

import (
	"testing"
	"time"
)

func TestCeilingTimezoneTimeNow(t *testing.T) {
	t.Log(CeilingTimezoneTimeNow(0))
	t.Log(CeilingTimezoneTimeNow(1))
	t.Log(CeilingTimezoneTimeNow(2))
	t.Log(CeilingTimezoneTimeNow(3))
	t.Log(CeilingTimezoneTimeNow(4))
	t.Log(CeilingTimezoneTimeNow(5))
	t.Log(CeilingTimezoneTimeNow(6))
	t.Log(CeilingTimezoneTimeNow(7))
	t.Log(CeilingTimezoneTimeNow(8))
	t.Log(CeilingTimezoneTimeNow(9))
	t.Log(CeilingTimezoneTimeNow(10))
	t.Log(CeilingTimezoneTimeNow(11))
	t.Log(CeilingTimezoneTimeNow(12))
	t.Log(CeilingTimezoneTimeNow(-1))
	t.Log(CeilingTimezoneTimeNow(-2))
	t.Log(CeilingTimezoneTimeNow(-3))
	t.Log(CeilingTimezoneTimeNow(-4))
	t.Log(CeilingTimezoneTimeNow(-5))
	t.Log(CeilingTimezoneTimeNow(-6))
	t.Log(CeilingTimezoneTimeNow(-7))
	t.Log(CeilingTimezoneTimeNow(-8))
	t.Log(CeilingTimezoneTimeNow(-9))
	t.Log(CeilingTimezoneTimeNow(-10))
	t.Log(CeilingTimezoneTimeNow(-11))
	t.Log(CeilingTimezoneTimeNow(-12))
}

func TestCeilingTimezoneTime(t *testing.T) {
	ctx, _ := time.Parse("15:04:05", "23:00:00")
	t.Log(ctx)
	startTime := CeilingTimezone(ctx)
	endTime := CeilingTimeNow()
	t.Log(startTime.Unix(), endTime.Unix())
}
