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
	add := flag.String("add", "", "在给定时间上增减时长，比如 1h30m、2d、90s")
	diff := flag.Bool("d", false, "把两个参数当时间点，输出它们之间的差（如 1700000000 1700003600）")
	ago := flag.String("ago", "", "输出当前时间往前推多少，比如 3h、2d（和 -now 类似但只给一个结果）")
	utc := flag.Bool("u", false, "输出统一用 UTC 时区，而不是本地时区")
	layout := flag.String("f", "", "自定义输出格式，比如 -f 2006/01/02，不设就用默认样式")
	base := flag.String("date", "", "指定基准时间（时间戳或 2006-01-02 15:04:05），-add/-ago 基于它而不是当前时间")
	flag.Parse()

	zone := time.Local
	if *utc {
		zone = time.UTC
	}

	// 当前时间上加减一段时长
	if *now && *add != "" {
		d, err := time.ParseDuration(*add)
		if err != nil {
			fmt.Fprintf(os.Stderr, "时长看不懂: %v\n", err)
			os.Exit(1)
		}
		show(time.Now().In(zone).Add(d), *layout)
		return
	}
	if *now {
		now := time.Now().In(zone)
		fmt.Printf("秒:   %d\n", now.Unix())
		fmt.Printf("毫秒: %d\n", now.UnixMilli())
		fmt.Printf("本地: %s\n", now.Format("2006-01-02 15:04:05"))
		fmt.Printf("UTC:  %s\n", now.UTC().Format(time.RFC3339))
		return
	}
	if *ago != "" {
		d, err := time.ParseDuration(*ago)
		if err != nil {
			fmt.Fprintf(os.Stderr, "时长看不懂: %v\n", err)
			os.Exit(1)
		}
		baseTm := time.Now()
		if *base != "" {
			baseTm, err = toTime(*base)
			if err != nil {
				fmt.Fprintf(os.Stderr, "基准时间看不懂: %v\n", err)
				os.Exit(1)
			}
		}
		tm := baseTm.In(zone).Add(-d)
		fmt.Printf("%s 之前是 %s（秒 %d）\n", *ago, tm.Format("2006-01-02 15:04:05"), tm.Unix())
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "给个时间戳或时间字符串，例如 go-timestamp 1700000000 或 go-timestamp -now")
		os.Exit(2)
	}

	// 两个时间点求差
	if *diff {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "-d 需要两个时间点")
			os.Exit(2)
		}
		a, err1 := toTime(args[0])
		b, err2 := toTime(args[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "时间解析失败: %v / %v\n", err1, err2)
			os.Exit(1)
		}
		dd := b.Sub(a)
		fmt.Printf("%s 到 %s 相差 %s（%d 秒）\n", args[0], args[1], dd, int64(dd.Seconds()))
		return
	}

	for _, a := range args {
		tm, err := toTime(a)
		if err != nil {
			fmt.Printf("%s -> 错误: %v\n", a, err)
			continue
		}
		tm = tm.In(zone)
		// 支持在某个时间点上加时长
		if *add != "" {
			d, err := time.ParseDuration(*add)
			if err != nil {
				fmt.Printf("%s -> 时长看不懂: %v\n", a, err)
				continue
			}
			tm = tm.Add(d)
		}
		show(tm, *layout)
	}
}

// show 按默认样式或自定义格式打印一个时间
func show(tm time.Time, layout string) {
	if layout != "" {
		// Go 的参考布局就是 2006-01-02 15:04:05，直接拿用户给的格式当模板
		fmt.Println(tm.Format(layout))
		return
	}
	fmt.Printf("%s\n  秒:   %d\n  毫秒: %d\n  本地: %s\n  UTC:  %s\n",
		tm.Format("2006-01-02 15:04:05"), tm.Unix(), tm.UnixMilli(),
		tm.Format("2006-01-02 15:04:05"), tm.UTC().Format(time.RFC3339))
}
