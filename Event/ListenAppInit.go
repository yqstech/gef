/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: ListenAppInit
 * @Version: 1.0.0
 * @Date: 2021/11/11 11:46 上午
 */

package Event

import (
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/config"
)

type ListenAppInit struct {
}

func (that ListenAppInit) Do(eventName string, data ...interface{}) (error, int) {
	if config.Debug != "" {
		logger.Debug("应用开始")
	}
	//*data[0].(*string) = "中国"
	//*data[1].(*string) = "加油"
	return nil, 0
}
