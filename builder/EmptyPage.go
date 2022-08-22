/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EmptyPage
 * @Version: 1.0.0
 * @Date: 2022/8/22 11:20
 */

package builder

import "github.com/julienschmidt/httprouter"

// EmptyPage 默认空白页
type EmptyPage struct {
	NodePage
}

// NodeInit 清空操作列表，自动全部指向Empty操作方法
func (d *EmptyPage) NodeInit(pageData *PageBuilder) {
	d.NodePageActions = map[string]httprouter.Handle{}
}
