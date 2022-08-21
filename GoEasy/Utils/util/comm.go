package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/yqstech/gef/GoEasy/Utils/snowflake"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

//密码加密
func GetPassword(str string) string {
	return MD5(MD5("*cradmin*" + MD5(str)))
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	return string(path[0 : i+1]), nil
}

//int64转int
func Int642Int(v int64) int {
	r, err := strconv.Atoi(strconv.FormatInt(v, 10))
	if err != nil {
		return 0
	}
	return r
}

//字符串转int
func String2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

//int转字符串
func Int2String(i int) string {
	s := strconv.Itoa(i)
	return s
}
func Int642String(i int64) string {
	s := strconv.FormatInt(i, 10)
	return s
}

//json字符串编码
func JsonEncode(data interface{}) string {
	if data == nil {
		return "[]"
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%d", err.Error())
		return ""
	} else {
		return string(jsondata)
	}
}

//json字符串解码
func JsonDecode(jsonstr string, data interface{}) {
	json.Unmarshal([]byte(jsonstr), data)
}

//任意对象转字符串
func Interface2String(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// Arr2String 判断是数组对象转json字符串
func Array2String(value interface{}) interface{} {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		return value
	case float32:
		return value
	case int:
		return value
	case uint:
		return value
	case int8:
		return value
	case uint8:
		return value
	case int16:
		return value
	case uint16:
		return value
	case int32:
		return value
	case uint32:
		return value
	case int64:
		return value
	case uint64:
		return value
	case string:
		return value
	case []byte:
		return value
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func Float642String(a float64) string {
	b := strconv.FormatFloat(a, 'f', -1, 64)
	return b
}

//字符串截取中间部分
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

//数字转中文称谓
func ChineseNum(str string) int {
	if len(str) == 0 {
		return 0
	}
	NumL := map[string]int{
		"亿": 100000000,
		"万": 10000,
		"千": 1000,
		"百": 100,
		"十": 10,
		"零": 0}
	NumLKeys := []string{"亿", "万", "千", "百", "十", "零"}
	CNums := map[string]int{
		"二十": 20,
		"三十": 30,
		"四十": 40,
		"五十": 50,
		"六十": 60,
		"七十": 70,
		"八十": 80,
		"九十": 90,
		"零":  0,
		"一":  1,
		"二":  2,
		"三":  3,
		"四":  4,
		"五":  5,
		"六":  6,
		"七":  7,
		"八":  8,
		"九":  9,
		"十":  10,
		"十一": 11,
		"十二": 12,
		"十三": 13,
		"十四": 14,
		"十五": 15,
		"十六": 16,
		"十七": 17,
		"十八": 18,
		"十九": 19}
	for k, v := range CNums {
		if str == k {
			return v
		}
	}
	//循环亿万千百
	for _, l := range NumLKeys {
		v := NumL[l]
		n := strings.Index(str, l)
		if n >= 0 {
			t1 := ChineseNum(string([]byte(str)[:n]))
			t2 := v
			m := n + len(l)
			t3 := ChineseNum(string([]byte(str)[m:]))
			return t1*t2 + t3
		}
	}
	return 0
}

//数组去重
func ArrayIntOnly(a []int) []int {
	b := []int{}
	c := map[int]int{}
	for _, v := range a {
		if _, ok := c[v]; !ok {
			c[v] = 1
			b = append(b, v)
		}
	}
	return b
}

//数组转字符串
func ArrayInt2String(a []int, sp string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", sp, -1)
}

//屏幕输出
func Echo(name string, value interface{}) {
	fmt.Printf(name+"=%d\n", value)
}

//雪花算法（服务器相关）
func GetOnlyID() int64 {
	snowflake.SetMachineId(1)
	id := snowflake.GetSnowflakeId()
	return id
}

/**
 * @description: 判断是否是数字字符串
 * @param {string} s
 * @return {bool}
 */
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

/**
 * @description: 浮点字符串转Int64
 * @param {string} fs 浮点格式的字符串
 * @param {float64} mult 乘数倍数
 * @return {*}
 */
func FloatString2Int64(fs string, mult float64) (int64, error) {
	//string转float64
	floatData, err := strconv.ParseFloat(fs, 64)
	if err != nil {
		return 0, err
	}
	//乘积后转INT64
	v := int64(math.Ceil(floatData * mult))
	return v, nil
}

// GenValidateCode 随机数
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// Paginator 分页器
// 当前页
// 每页数量
// 总数量
// 显示分页数量 ShowPageNum
// 链接地址
// 分页关键词
func Paginator(CurrentPageNum, PageSize int, Total, ShowPageNum int, Url, pageKey string) []PagerItem {

	//计算得到最大页码
	maxPageNum := int(math.Ceil(float64(Total) / float64(PageSize))) //总页数
	//矫正页码数值
	if maxPageNum < 1 {
		maxPageNum = 1
	}
	if CurrentPageNum > maxPageNum {
		CurrentPageNum = maxPageNum
	}
	if CurrentPageNum <= 0 {
		CurrentPageNum = 1
	}
	//计算起始页码
	beginPageNum := CurrentPageNum - (ShowPageNum / 2)
	if beginPageNum < 1 {
		//矫正起始页码
		beginPageNum = 1
	}
	//计算截止页码
	endPageNum := beginPageNum + ShowPageNum
	if endPageNum > maxPageNum {
		//矫正截止页码
		endPageNum = maxPageNum
	}
	//定义分页页码列表
	var pageNums []int
	for i := beginPageNum; i <= endPageNum; i++ {
		pageNums = append(pageNums, i)
	}
	//计算前一页后一页
	prevPageNum := CurrentPageNum - 1
	if prevPageNum < 1 {
		prevPageNum = CurrentPageNum
	}
	nextPageNum := CurrentPageNum + 1
	if nextPageNum > maxPageNum {
		nextPageNum = CurrentPageNum
	}
	//链接前缀
	PageURL := "?" + pageKey + "="
	if strings.Contains(Url, "?") {
		PageURL = "&" + pageKey + "="
	}
	var PageData []PagerItem
	PageData = append(PageData, PagerItem{
		Title:   "上一页",
		Url:     Interface2String(Is(prevPageNum == 1, Url, Url+PageURL+Int2String(prevPageNum))),
		PageNum: prevPageNum,
		Current: Is(prevPageNum == CurrentPageNum, true, false).(bool),
		Tag:     "prev",
	})
	for _, num := range pageNums {
		PageData = append(PageData, PagerItem{
			Title:   Int2String(num),
			Url:     Interface2String(Is(num == 1, Url, Url+PageURL+Int2String(num))),
			PageNum: num,
			Current: Is(num == CurrentPageNum, true, false).(bool),
			Tag:     "",
		})
	}
	PageData = append(PageData, PagerItem{
		Title:   "下一页",
		Url:     Interface2String(Is(nextPageNum == 1, Url, Url+PageURL+Int2String(nextPageNum))),
		PageNum: nextPageNum,
		Current: Is(nextPageNum == CurrentPageNum, true, false).(bool),
		Tag:     "next",
	})
	return PageData
}

// PagerItem 分页页码结构
type PagerItem struct {
	Title   string //显示名称 如 上一页 1
	Url     string //链接地址
	PageNum int    //页码
	Current bool   //是否是当前页码
	Tag     string //标记 prev 上一页 next 下一页
}

//模拟三元运算

func Is(is bool, opt1 interface{}, opt2 interface{}) interface{} {
	if is {
		return opt1
	} else {
		return opt2
	}
}
