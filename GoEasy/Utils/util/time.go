package util

import (
	"strings"
	"time"
)

//获取时间字符串
func Time() string {
	return Int2String(Int642Int(time.Now().Unix()))
}

func TimeNow() string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	return time.Now().In(cstSh).Format("2006-01-02 15:04:05")
}

func TimeNowFormat(fmt string, y, m, d int) string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	return time.Now().AddDate(y, m, d).In(cstSh).Format(fmt)
}

func TimeNowWeek(y, m, d int) (string, string) {
	var WeekDayMap = map[string]string{
		"Monday":    "周一",
		"Tuesday":   "周二",
		"Wednesday": "周三",
		"Thursday":  "周四",
		"Friday":    "周五",
		"Saturday":  "周六",
		"Sunday":    "周日",
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	var week = time.Now().AddDate(y, m, d).In(cstSh).Weekday().String()
	return week, WeekDayMap[week]
}

func Str2UnixTime(t string) int64 {
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	return stamp.Unix()
}
func UnixTimeFormat(t int64, fmt string) string {
	return time.Unix(t, 0).Format(fmt)
}

// TimeStringSplit 时间分割
func TimeStringSplit(t string) []string {
	//杂乱字符串替换成,
	t = strings.Replace(t, "-", ",", -1)
	t = strings.Replace(t, "/", ",", -1)
	t = strings.Replace(t, "_", ",", -1)
	t = strings.Replace(t, " ", ",", -1)
	t = strings.Replace(t, ":", ",", -1)
	return strings.Split(t, ",")
}

// TimeLeft 计算时间间隔，
// return time2 - time1
func TimeLeft(t1, t2 string) int64 {
	time1 := Str2UnixTime(t1)
	time2 := Str2UnixTime(t2)
	if time1 >= time2 {
		return 0
	} else {
		return time2 - time1
	}
}
