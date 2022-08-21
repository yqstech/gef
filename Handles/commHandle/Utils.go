/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Utils
 * @Version: 1.0.0
 * @Date: 2022/2/7 8:59 下午
 */

package commHandle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"net/http"
	"net/http/httputil"
)

type Utils struct {
}

func (this Utils) RequestPrint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println(string(data))
}
