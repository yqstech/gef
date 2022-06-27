package util

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/valyala/fasthttp"
)

// FileSystems 组合多个http.FileSystem
type FileSystems []http.FileSystem

func (fs FileSystems) Open(name string) (file http.File, err error) {
	for _, i := range fs {
		file, err = i.Open(name)
		if err == nil {
			return file, err
		}
	}
	return
}

// FastHttpGet 使用fasthttp包里的get方式进行请求，CPU资源占用低，速度快
func FastHttpGet(url string) (string, error) {
	client := &fasthttp.Client{
		ReadTimeout: time.Second * 8,
	}
	status, resp, err := client.Get(nil, url)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return "", err
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return "", errors.New("请求失败！")
	}
	return string(resp), nil
}
func FastHttpPost(url string, data map[string]string) (string, error) {
	//填充表单
	args := &fasthttp.Args{}
	if data != nil {
		for k, v := range data {
			args.Add(k, v)
		}
	}
	client := &fasthttp.Client{
		ReadTimeout: time.Second * 8,
	}
	//执行POST
	status, resp, err := client.Post(nil, url, args)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return "", errors.New("请求失败")
	}
	if status != fasthttp.StatusOK {
		fmt.Println(url+"请求没有成功:", status)
		return "", errors.New("请求没有成功")
	}
	return string(resp), nil
}

// FastHttpPostH 可定义header的post请求
func FastHttpPostH(url string, data, headers map[string]string) (string, error) {
	req := &fasthttp.Request{}
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	//设置headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	//设置数据
	args := &fasthttp.Args{}
	if data != nil {
		for k, v := range data {
			args.Add(k, v)
		}
	}
	args.WriteTo(req.BodyWriter())
	//设置客户端
	client := &fasthttp.Client{
		ReadTimeout: time.Second * 8,
	}
	//设置响应对象
	resp := &fasthttp.Response{}
	//执行post
	err := client.Do(req, resp)
	if err != nil {
		return "", err
	}
	return string(resp.Body()), nil
}

// PostValue 获取单个post值
func PostValue(r *http.Request, k string) string {
	v := r.PostFormValue(k)
	return html.EscapeString(v)
}

// PostValueDef 获取post值并带默认值
func PostValueDef(r *http.Request, k string, defValue string) string {
	v := r.PostFormValue(k)
	v = html.EscapeString(v)
	if v == "" {
		return defValue
	} else {
		return v
	}
}

func PostJson(r *http.Request, k string) map[string]interface{} {
	post_data := r.PostFormValue(k)
	trs_data := make(map[string]interface{})
	JsonDecode(post_data, &trs_data)

	return trs_data
}
func Clearmap(a map[string]interface{}) map[string]interface{} {
	d := make(map[string]interface{})
	return d
}

func PostJsonArray(r *http.Request, k string) []map[string]string {
	post_data := r.PostFormValue(k)
	trs_data := []map[string]string{}
	rst_data := []map[string]string{}
	JsonDecode(post_data, &trs_data)
	for _, item := range trs_data {
		d := make(map[string]string)
		for key, value := range item {
			d[key] = html.EscapeString(value)
		}
		rst_data = append(rst_data, d)
	}
	return rst_data
}

//获取post数组数据
func PostValues(r *http.Request, k string, num int) []string {
	var d []string
	var v string
	for i := 0; i < num; i++ {
		v = r.PostFormValue(k + "[" + strconv.Itoa(i) + "]")
		if v != "" {
			d = append(d, html.EscapeString(v))
		}
	}
	return d
}

//获取单个get值
func GetValue(r *http.Request, k string) string {
	query := r.URL.Query()
	if len(query[k]) > 0 {
		v := query[k][0]
		return html.EscapeString(v)
	}
	return ""
}
func GetValueInt(r *http.Request, k string) int {
	return String2Int(GetValue(r, k))
}

//获取单个路由参数值
func RouterValue(ps httprouter.Params, k string) string {
	v := ps.ByName(k)
	return html.EscapeString(v)
}

//! 获取真实ip地址，兼容nginx反向代理
func GetRealIP(r *http.Request) string {
	IP := r.RemoteAddr
	IPinfo := strings.Split(IP, ":")
	IP_ADDR := IPinfo[0]
	RealIP := r.Header.Get("X-Real-IP")
	if RealIP != "" && IP_ADDR == "127.0.0.1" {
		IP_ADDR = RealIP
	}
	return IP_ADDR
}
func UrlEncode(data string) string {
	post_data := url.Values{}
	post_data.Add("data", data)
	q := post_data.Encode()
	return q[5:]
}

// UrlJoin URL拼接，自动判断？和 &
func UrlJoin(url *string, add string) {
	if strings.Contains(*url, "?") {
		*url = *url + "&" + add
	} else {
		*url = *url + "?" + add
	}
}

// UrlRemoveParam 去掉链接里的一部分
func UrlRemoveParam(url, param string) string {
	//字符串分割
	parts := strings.Split(url, param+"=")
	if len(parts) == 1 {
		return url
	}
	//分割符号
	parts2 := strings.Split(parts[1], "&")
	parts2[0] = ""
	return parts[0] + strings.Join(parts2, "&")
}

// UrlScreenParam url筛选参数
// saveNull bool 保留空值
// endSymbol string 保留结尾分割符号，用于手动拼接其他参数
func UrlScreenParam(r *http.Request, params []string, saveNull bool, endSymbol bool) string {
	requestUrl := r.RequestURI
	//首页默认添加斜杠
	if requestUrl == "" {
		requestUrl = "/"
	}
	urlParts := strings.Split(requestUrl, "?")
	url := urlParts[0]
	var paramUrl []string
	for _, param := range params {
		v := GetValue(r, param)
		if v != "" || saveNull {
			paramUrl = append(paramUrl, param+"="+v)
		}
	}
	p := strings.Join(paramUrl, "&")
	if p == "" {
		if endSymbol {
			return url + "?"
		}
		return url
	}
	url = strings.Join([]string{url, p}, "?")
	if endSymbol {
		return url + "&"
	}
	return url
}

// UrlCompletion url补全
func UrlCompletion(url, baseUrl string) string {
	//检查地址是否更换了
	if len(url) > 4 && url[0:4] != "http" {
		//相对地址补全斜杠
		if len(url) > 0 && url[0:1] != "/" {
			url = "/" + url
		}
		return baseUrl + url
	}
	return url
}

// UrlRemoveDomain 地址删除前面的http域名
func UrlRemoveDomain(url string) string {
	if len(url) > 10 && url[0:4] == "http" {
		url = url[10:] //去掉http://或https://
		//查找下一个斜杠
		index := strings.Index(url, "/")
		if index > -1 {
			//取斜杠后面的内容
			url = url[index:]
		}
		return url
	}
	return url
}
