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
