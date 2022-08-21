/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 空白页
 * @File: EmptyPage
 * @Version: 1.0.0
 * @Date: 2021/10/18 4:24 下午
 */

package EasyApp

// EmptyPage 默认空白页
type EmptyPage struct {
	Page
}

// PageInit 清空操作列表，自动全部指向Empty操作方法
func (d EmptyPage) PageInit(pageData *PageData) {
	pageData.ActionClear()
}
