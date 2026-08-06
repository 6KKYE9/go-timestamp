package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// 把命令行给的时间戳或时间字符串互相转一下
// 10 位是秒，13 位是毫秒，其余按 RFC3339 / 常见格式尝试
func toTime(arg string) (time.Time, error) {
	arg = strings.TrimSpace(arg)
	if n, err := strconv.ParseInt(arg, 10, 64); err == nil {
		switch {
		case n >= 1e12: // 毫秒
			return time.UnixMilli(n).UTC(), nil
		case n >= 1e10: // 可能已经是纳秒量级，极少，按秒兜底
			return time.Unix(n, 0).UTC(), nil
		default:
			return time.Unix(n, 0).UTC(), nil
		}
	}
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	} {
		if tm, err := time.ParseInLocation(layout, arg, time.Local); err == nil {
			return tm, nil
		}
	}
	return time.Time{}, fmt.Errorf("看不懂这个时间: %s", arg)
}

func main() {
	now := flag.Bool("now", false, "打印当前时间戳（秒和毫秒）")
	flag.Parse()

	if *now {
		now := time.Now()
		fmt.Printf("秒:   %d\n", now.Unix())
		fmt.Printf("毫秒: %d\n", now.UnixMilli())
		fmt.Printf("本地: %s\n", now.Format("2006-01-02 15:04:05"))
		fmt.Printf("UTC:  %s\n", now.UTC().Format(time.RFC3339))
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "给个时间戳或时间字符串，例如 go-timestamp 1700000000 或 go-timestamp -now")
		os.Exit(2)
	}
	for _, a := range args {
		tm, err := toTime(a)
		if err != nil {
			fmt.Printf("%s -> 错误: %v\n", a, err)
			continue
		}
		fmt.Printf("%s\n  秒:   %d\n  毫秒: %d\n  本地: %s\n  UTC:  %s\n",
			a, tm.Unix(), tm.UnixMilli(), tm.Format("2006-01-02 15:04:05"), tm.UTC().Format(time.RFC3339))
	}
}
