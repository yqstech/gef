package commHandle

import (
	"fmt"
	"github.com/yqstech/gef/Utils/util"
	"net/http"
	"os"
	"os/exec"
	
	"github.com/julienschmidt/httprouter"
)

// Server 服务管理
type Server struct {
}

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

func (that Server) Pid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := os.Getpid()
	count = count + 1
	fmt.Fprint(w, pid)
}
