package main

import (
	"testing"
	"time"
)

func TestToTimeUnix(t *testing.T) {
	tm, err := toTime("1700000000")
	if err != nil {
		t.Fatal(err)
	}
	if tm.Unix() != 1700000000 {
		t.Fatalf("秒时间戳解析错: %d", tm.Unix())
	}
}

func TestToTimeMilli(t *testing.T) {
	tm, err := toTime("1700000000000")
	if err != nil {
		t.Fatal(err)
	}
	if tm.UnixMilli() != 1700000000000 {
		t.Fatalf("毫秒时间戳解析错: %d", tm.UnixMilli())
	}
}

func TestToTimeLayout(t *testing.T) {
	tm, err := toTime("2023-11-14 22:13:20")
	if err != nil {
		t.Fatal(err)
	}
	if tm.Year() != 2023 || tm.Month() != time.November || tm.Day() != 14 {
		t.Fatalf("时间字符串解析错: %v", tm)
	}
}

func TestToTimeBad(t *testing.T) {
	if _, err := toTime("不是时间"); err == nil {
		t.Fatal("应该解析失败")
	}
}

func TestParseDurationHelpers(t *testing.T) {
	// -add / -ago 走的都是 time.ParseDuration，这里验证其覆盖常见写法
	cases := map[string]bool{
		"1h30m": true, "48h": true, "90s": true, "1h": true,
	}
	for d, wantOK := range cases {
		_, err := time.ParseDuration(d)
		if (err == nil) != wantOK {
			t.Fatalf("ParseDuration(%s) 结果不符合预期", d)
		}
	}
}
