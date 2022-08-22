/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 模板操作库
 * @File: Displayer
 * @Version: 1.0.0
 * @Date: 2021/10/15 12:10 下午
 */

package builder

import (
	"bytes"
	"embed"
	"github.com/yqstech/gef/builder/adminTemplates"
	"strings"
	"text/template"
)

// Displayer 渲染器
type Displayer struct {
	//模板数据
	DisplayData map[string]interface{}
	//模板名称
	TplName string
	//模板文件集
	Templates []TemplateParseFS
	//模板函数集
	FuncMap template.FuncMap
}

// TemplateParseFS 页面构建器模板文件对象
type TemplateParseFS struct {
	Fsys     embed.FS
	Patterns []string
}

// SetDate 设置模板数据方法
func (t *Displayer) SetDate(key string, value interface{}) {
	if len(t.DisplayData) == 0 {
		t.DisplayData = map[string]interface{}{}
	}
	if key == "" {
		t.DisplayData = value.(map[string]interface{})
	} else {
		t.DisplayData[key] = value
	}
}

// Display 模板渲染生成html
func (t *Displayer) Display() (string, error) {
	var err error
	//加载默认方法
	insertFuncMap := t.Functions()
	if t.FuncMap == nil {
		t.FuncMap = insertFuncMap
	} else {
		for k, v := range insertFuncMap {
			if _, ok := insertFuncMap[k]; !ok {
				t.FuncMap[k] = v
			}
		}
	}
	//创建模板对象
	tpl := template.New(t.TplName).Funcs(t.FuncMap)
	//加载内置模板
	tpl, err = tpl.ParseFS(adminTemplates.Templates, "layout/*.html")
	if err != nil {
		return "", err
	}
	//加载拓展模板
	for _, f := range t.Templates {
		tpl, err = tpl.ParseFS(f.Fsys, f.Patterns...)
	}
	//加载小组件
	tpl, err = tpl.ParseFS(adminTemplates.Templates, "widget/*.html")
	if err != nil {
		return "", err
	}
	//校验模板
	tpl = template.Must(tpl, err)
	//渲染模板
	var buf bytes.Buffer
	err = tpl.Execute(&buf, t.DisplayData)
	if err != nil {
		return "", err
	}
	//获取html代码
	return buf.String(), nil
}

// Functions 获取内置模板
func (t *Displayer) Functions() template.FuncMap {
	return template.FuncMap{
		//任意值转字符串
		"toString": Interface2String,
		"inArray":  inArray,
		"inMap":    inMap,
		"mapValue": func(k string, data map[string]interface{}) interface{} {
			if v, ok := data[k]; ok {
				return v
			}
			return nil
		},
		"json_encode": func(data interface{}, trans bool) string {
			str := JsonEncode(data)
			if trans {
				//模板里使用需要转义
				str = strings.Replace(str, "\"", "&#34;", -1)
			}
			return str
		},
		"hasType": func(formFields []FormField, types string) bool {
			if types == "" {
				return false
			}
			typeList := strings.Split(types, ",")
			for _, formField := range formFields {
				for _, t := range typeList {
					if formField.Type == t {
						return true
					}
				}
			}
			return false
		},
		"jsValue": func(value interface{}) interface{} {
			s := Interface2String(value)
			if IsNum(s) {
				return s
			} else {
				return "'" + s + "'"
			}
		},
		"isNum": IsNum,
	}
}
