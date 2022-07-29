/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 模板操作库
 * @File: Template
 * @Version: 1.0.0
 * @Date: 2021/10/15 12:10 下午
 */

package EasyApp

import (
	"bytes"
	"github.com/gef/GoEasy/Templates"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
	"html"
	"io/fs"
	"strings"
	"text/template"
)

//节点流匹配的模板

type Template struct {
	DisplayData  map[string]interface{}
	TplName      string
	TemplateSelf []string
}

// TemplateAdd 模板追加地址
func (t *Template) TemplateAdd(tplPath string) {
	have := false
	for _, tpl := range t.TemplateSelf {
		if tpl == tplPath {
			have = true
		}
	}
	if !have {
		t.TemplateSelf = append(t.TemplateSelf, tplPath)
	}
}

func (t *Template) SetDate(key string, value interface{}) {
	if len(t.DisplayData) == 0 {
		t.DisplayData = map[string]interface{}{}
	}
	if key == "" {
		t.DisplayData = value.(map[string]interface{})
	} else {
		t.DisplayData[key] = value
	}
}

// PageData2Display 设置模板数据
func (t *Template) PageData2Display(pageData *PageData) {
	data := t.DisplayData
	if data == nil {
		data = make(map[string]interface{})
	}
	//页面样式
	data["pageStyle"] = pageData.PageStyle
	//设置标题
	if data["title"] == nil {
		data["title"] = pageData.title
	}
	data["pageName"] = pageData.pageName
	data["pageNotice"] = pageData.pageNotice
	//列表列
	data["listColumns"] = pageData.listColumns
	data["listColumnsStyles"] = pageData.listColumnsStyles
	//列表数据地址
	data["listDataUrl"] = pageData.listDataUrl
	//列表按钮
	rightBtns := []Button{}
	for _, button_name := range pageData.listRightBtns {
		if btn, ok := pageData.buttons[button_name]; ok {
			rightBtns = append(rightBtns, btn)
		}
	}
	data["rightBtns"] = rightBtns
	//列表顶部按钮
	topBtns := []Button{}
	for _, button_name := range pageData.listTopBtns {
		if btn, ok := pageData.buttons[button_name]; ok {
			topBtns = append(topBtns, btn)
		}
	}
	data["topBtns"] = topBtns
	//列表搜索项
	data["listSearchFields"] = pageData.listSearchFields
	//列表Tab
	data["pageTabs"] = pageData.pageTabs
	data["pageTabSelect"] = pageData.pageTabSelect
	data["pageTabsLength"] = len(pageData.pageTabs)
	//分页
	data["listPageHide"] = pageData.listPageHide //是否隐藏分页
	data["listPageSize"] = pageData.listPageSize //分页
	//批量操作
	data["listBatchAction"] = pageData.listBatchAction //批量操作
	
	//表单
	//提交地址
	data["editDataUrl"] = pageData.editDataUrl
	data["addDataUrl"] = pageData.addDataUrl
	//表单项列表
	for k, v := range pageData.formFields {
		fkey := v.Key
		ftype := v.Type
		if fkey != "" {
			//初始化表单项的原始数据值
			if dbvalue, ok := pageData.formData[fkey]; ok {
				pageData.formFields[k].Value = util.Interface2String(dbvalue)
			}
			//修正checkbox默认值
			if util.IsInArray(ftype, []interface{}{"checkbox", "checkbox_level", "tags", "images"}) {
				if pageData.formFields[k].Value == "" {
					pageData.formFields[k].Value = "[]"
				}
			}
			//文本类组件需要转义双引号
			if util.IsInArray(ftype, []interface{}{"text", "textarea", "text-disabled", "textarea-disabled", "wangEditor"}) {
				pageData.formFields[k].Value = strings.ReplaceAll(pageData.formFields[k].Value, "\"", "\\\"")
			}
			//文本域,换行\n 避免被html解析造成换行错乱
			if ftype == "textarea" || ftype == "textarea-disabled" {
				pageData.formFields[k].Value = strings.ReplaceAll(pageData.formFields[k].Value, "\n", "\\n")
			}
			if ftype == "code" || ftype == "codeEditor" {
				pageData.formFields[k].Value = html.EscapeString(pageData.formFields[k].Value)
			}
			
		}
	}
	data["formFields"] = pageData.formFields
	
	if pageData.formSubmitTitle == "" {
		pageData.formSubmitTitle = pageData.actionName + pageData.pageName
	}
	data["formSubmitTitle"] = pageData.formSubmitTitle
	
	data["formSubmitHide"] = pageData.formSubmitHide
	
	//组件需要
	//上传图片的地址
	data["uploadImageUrl"] = pageData.uploadImageUrl
	
	//组件默认
	t.DisplayData = data
}

// Display 模板渲染，模板名称，模板数据
func (t *Template) Display() (string, error) {
	var err error
	//加载默认页面
	//tpl := template.Must(template.New(tplName).Funcs(t.Funcs()).ParseGlob("templates/qadmin/layout/*.html"))
	tpl := template.Must(
		template.New(t.TplName).Funcs(t.Functions()).
			ParseFS(Templates.Files, "qadmin/layout/*.html"),
	)
	fd, err := fs.ReadDir(Templates.FilesAdd, "admin")
	if err == nil {
		if len(fd) > 0 {
			//加载替换页面
			tpl, err = tpl.ParseFS(Templates.FilesAdd, "admin/*.html")
			if err != nil {
				logger.Error(err.Error())
				return "", err
			}
		}
	}
	//加载自定义页面
	if len(t.TemplateSelf) > 0 {
		tpl, err = tpl.ParseFS(Templates.FilesSelf, t.TemplateSelf...)
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}
	}
	//加载小组件
	tpl, err = tpl.ParseFS(Templates.Files, "qadmin/widget/*.html")
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	/**
	  渲染页面
	*/
	var buf bytes.Buffer
	err = tpl.Execute(&buf, t.DisplayData)
	if err != nil {
		return "", err
	}
	//获取html代码
	return buf.String(), nil
}

func (t *Template) Functions() template.FuncMap {
	return template.FuncMap{
		"toString": util.Interface2String,
		"inArray":  util.IsInArray,
		"inMap": func(k string, data map[string]interface{}) bool {
			if _, ok := data[k]; ok {
				return true
			}
			return false
		},
		"mapValue": func(k string, data map[string]interface{}) interface{} {
			if v, ok := data[k]; ok {
				return v
			}
			return nil
		},
		"json_encode": func(data interface{}, trans bool) string {
			str := util.JsonEncode(data)
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
			s := util.Interface2String(value)
			if util.IsNum(s) {
				return s
			} else {
				return "'" + s + "'"
			}
		},
	}
}
