package util

import (
	"regexp"
	"strings"
	"time"
)

// TimeGo 时光机
func TimeGo(options ...interface{}) string {
	optionsLen := len(options)

	fmt := "2006-01-02 15:04:05"
	if optionsLen > 0 && options[0].(string) != "" {
		fmt = options[0].(string)
	}
	y := 0
	m := 0
	d := 0
	h := 0
	i := 0
	s := 0
	if optionsLen > 1 {
		y = options[1].(int)
	}
	if optionsLen > 2 {
		m = options[2].(int)
	}
	if optionsLen > 3 {
		d = options[3].(int)
	}
	if optionsLen > 4 {
		h = options[4].(int)
	}
	if optionsLen > 5 {
		i = options[5].(int)
	}
	if optionsLen > 6 {
		s = options[6].(int)
	}

	cstSh, _ := time.LoadLocation("Asia/Shanghai") //上海
	t := time.Now().AddDate(y, m, d).In(cstSh)

	if h != 0 || i != 0 || s != 0 {
		//时分秒需要变化,先转成时间戳计算完再格式化
		tStr := t.Format("2006-01-02 15:04:05")
		tInt64 := Str2UnixTime(tStr) + int64(h)*3600 + int64(i)*60 + int64(s)
		return UnixTimeFormat(tInt64, fmt)
	}
	return t.Format(fmt)
}

// 获取时间字符串
func Time() string {
	return Int2String(Int642Int(time.Now().Unix()))
}

func IfTimeFmt(t interface{}) interface{} {
	switch t.(type) {
	case time.Time:
		return t.(time.Time).Format("2006-01-02 15:04:05")
	default:
		return t
	}
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

// Second2Dhms 秒格式化成 天 时 分 秒
func Second2Dhms(second int64, d, h, m, s string) string {
	if second < 0 {
		second = 0
	}
	day := second / 86400
	hour := (second - (day * 86400)) / 3600
	minute := (second - day*86400 - hour*3600) / 60
	sec := second - day*86400 - hour*3600 - minute*60
	result := ""
	//天数大于0时才显示天
	if day > 0 {
		result = result + Int642String(day) + d
	}
	//天或小时大于0时显示
	if hour > 0 || day > 0 {
		result = result + Int642String(hour) + h
	}
	//天时分大于0时显示分
	if minute > 0 || hour > 0 || day > 0 {
		result = result + Int642String(minute) + m
	}
	//一直显示秒
	result = result + Int642String(sec) + s

	return result
}

// Dhms2Second 时分秒转换成秒
func Dhms2Second(str string) int64 {
	second := int64(0)
	//正则匹配所有数字
	re := regexp.MustCompile("[0-9]+")
	nums := re.FindAllString(str, -1)
	//数字倒序
	var timeArr []int64
	l := int64(len(nums))
	if l > 0 {
		l--
		for true {
			if l < 0 {
				break
			}
			timeArr = append(timeArr, int64(String2Int(nums[l])))
		}
	}
	//倒序后的顺序为秒、分、时、天，累加秒数
	for k, v := range timeArr {
		if k == 0 {
			second = second + v
		} else if k == 1 {
			second = second + v*60
		} else if k == 2 {
			second = second + v*3600
		} else if k == 3 {
			second = second + v*86400
		} else {
			break
		}
	}
	return second
}
