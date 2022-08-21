package EasyApp

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// AppBoot 应用启动器
type AppBoot struct {
}

// BootHandle 定义启动器方法结构
type BootHandle func(map[string]AppPage, http.ResponseWriter, *http.Request, httprouter.Params)

// Run 应用启动器默认启动方法
// 传入应用页面集和http参数
func (ar AppBoot) Run(AppPages map[string]AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//根据url链接，获取当前页面名称
	pageName := ps.ByName("pageName")
	//在页面集中检索页面，并执行页面的Run方法
	if activePage, ok := AppPages[pageName]; ok {
		activePage.Run(&activePage, w, r, ps)
	} else {
		//检索不到，转发到空白页
		activePage = EmptyPage{}
		activePage.Run(&activePage, w, r, ps)
	}
}

func (ar AppBoot) SafeHandle(next BootHandle) BootHandle {
	return func(appPages map[string]AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//生产环境需要打开
		//defer func() {
		//	if r := recover(); r != nil {
		//		fmt.Fprint(w, "程序异常：", r)
		//	}
		//}()
		if r.Method == "POST" {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			_body := httprouter.Param{Key: "_body", Value: string(bodyBytes)}
			ps = append(ps, _body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		next(appPages, w, r, ps)
	}
}
