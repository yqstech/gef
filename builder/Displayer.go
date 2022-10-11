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
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder/adminTemplates"
	"github.com/yqstech/gef/util"
	"strings"
	"text/template"
	"time"
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
		"isNum":      IsNum,
		"select":     dbSelect,
		"count":      dbCount,
		"find":       dbFind,
		"optionName": optionName,
		"isIn":       isIn,
		"in":         isIn,
	}
}

// 查询列表 tableName pageSize page where orderBy
func dbSelect(options ...interface{}) []gorose.Data {
	optionsLen := len(options)
	tableName := ""
	if optionsLen > 0 {
		tableName = util.Interface2String(options[0])
	} else {
		return nil
	}
	pageSize := 10
	if optionsLen > 1 {
		pageSize = util.String2Int(util.Interface2String(options[1]))
	}
	page := 1
	if optionsLen > 2 {
		page = util.String2Int(util.Interface2String(options[2]))
	}
	where := "is_delete = 0"
	if optionsLen > 3 {
		getWhere := util.Interface2String(options[3])
		if getWhere != "" {
			where = where + " and " + getWhere
		}
	}
	orderBy := "id desc"
	if optionsLen > 4 {
		orderBy = util.Interface2String(options[4])
	}

	conn := db.New()
	data, err := conn.Table(tableName).
		Where(where).
		Limit(pageSize).
		Page(page).
		OrderBy(orderBy).Get()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	for i, d := range data {
		//!针对于time类型（sqlite数据库）转为字符串
		for k, v := range d {
			switch v.(type) {
			case time.Time:
				d[k] = v.(time.Time).Format("2006-01-02 15:04:05")
			default:
			}
		}
		data[i] = d
	}
	return data
}

// 查询列表 tableName where
func dbCount(options ...interface{}) int64 {
	optionsLen := len(options)
	tableName := ""
	if optionsLen > 0 {
		tableName = util.Interface2String(options[0])
	} else {
		return 0
	}
	where := "is_delete = 0"
	if optionsLen > 3 {
		getWhere := util.Interface2String(options[3])
		if getWhere != "" {
			where = where + " and " + getWhere
		}
	}
	conn := db.New()
	count, err := conn.Table(tableName).Where(where).Count()
	if err != nil {
		logger.Error(err.Error())
		return 0
	}
	return count
}

// 模板查询单条数据,tableName where orderBy
func dbFind(options ...interface{}) gorose.Data {
	optionsLen := len(options)
	tableName := ""
	if optionsLen > 0 {
		tableName = util.Interface2String(options[0])
	}
	where := "is_delete = 0"
	if optionsLen > 1 {
		getWhere := util.Interface2String(options[1])
		if getWhere != "" {
			where = where + " and " + getWhere
		}
	}
	orderBy := "id desc"
	if optionsLen > 2 {
		orderBy = util.Interface2String(options[2])
	}

	conn := db.New()
	data, err := conn.Table(tableName).Where(where).OrderBy(orderBy).First()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	for k, v := range data {
		//时间格式转成字符串格式
		switch v.(type) {
		case time.Time:
			data[k] = v.(time.Time).Format("2006-01-02 15:04:05")
		default:
		}
	}
	return data
}

// 传入选项集key和选项key，得到选项名称
func optionName(optionModelsKey string, optionKey string) string {
	optionList := Models.OptionModels{}.ByKey(optionModelsKey, false)
	optionMap := map[string]interface{}{}
	for _, opt := range optionList {
		optionMap[util.Interface2String(opt["value"])] = opt["name"]
	}
	if v2, ok2 := optionMap[optionKey]; ok2 {
		return util.Interface2String(v2)
	}
	return ""
}

func isIn(options ...interface{}) bool {
	if len(options) < 2 {
		return false
	}
	first := util.Interface2String(options[0])
	for k, v := range options {
		if k > 0 {
			if first == util.Interface2String(v) {
				return true
			}
		}
	}
	return false
}
