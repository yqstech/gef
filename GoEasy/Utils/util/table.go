/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: table.go
 * @Version: 1.0.0
 * @Date: 2022/2/17 7:37 下午
 */

package util

type Table struct {
	Columns []string
	Rows    []string
	Html    string
}

// Row 生成表格一行的html代码并存储到Html内
func (that *Table) Row(data ...string) error {
	html := "<tr>"
	for _, v := range data {
		html += "<td>" + v + "</td>"
	}
	if len(data) < len(that.Columns) {
		c := len(that.Columns) - len(data)
		for i := c; i > 1; i-- {
			html += "<td> </td>"
		}
	}
	html += "</tr>"
	that.Html += html
	return nil
}

// Display 包装表格的html代码
func (that *Table) Display(class string) string {
	html := "<table class=\"xtable " + class + "\">"
	if len(that.Columns) > 0 {
		html += "<tr>"
		for _, v := range that.Columns {
			html += "<td>" + v + "</td>"
		}
		html += "</tr>"
	}
	html += that.Html
	html += "</table>"
	return html
}
