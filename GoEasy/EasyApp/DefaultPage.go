/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 默认操作handle，所有未知的操作均转发至此
 * @File: DefaultPage
 * @Version: 1.0.0
 * @Date: 2021/10/18 4:24 下午
 */

package EasyApp

// DefaultPage 默认空页面
type DefaultPage struct {
	Page
}

// PageInit 清空操作列表，自动全部指向Empty操作方法
func (d DefaultPage) PageInit(pageData *PageData) {
	pageData.ActionClear()
}
