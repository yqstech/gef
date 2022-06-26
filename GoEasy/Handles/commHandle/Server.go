/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-04-01 16:21:23
 * @LastEditTime: 2021-06-17 14:59:14
 * @Description: 服务管理
 */
package commHandle

import (
	"fmt"
	"github.com/gef/GoEasy/Utils/util"
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
}

/**
 * @description: 重启服务实例
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @param {httprouter.Params} ps
 * @return {*}
 */
func (that Server) Restart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := os.Getpid()
	command := exec.Command("bash", "-c", "kill -USR2 "+util.Int2String(pid))
	e2 := command.Start()
	if e2 != nil {
		panic(e2.Error())
	}
	fmt.Fprint(w, "ok")
}

var count = 0

/**
 * @description: 获取服务实例PID
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @param {httprouter.Params} ps
 * @return {*}
 */
func (that Server) Pid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := os.Getpid()
	count = count + 1
	fmt.Fprint(w, pid)
}
