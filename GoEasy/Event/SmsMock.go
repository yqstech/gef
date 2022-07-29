/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsAli
 * @Version: 1.0.0
 * @Date: 2022/7/29 23:03
 */

package Event

import (
	"github.com/wonderivan/logger"
)

type SmsMock struct {
}

func (that SmsMock) Do(eventName string, data ...interface{}) (error, int) {
	programs := data[0].(map[string]interface{})
	for k, v := range programs {
		logger.Info("SmsMock", k, v)
	}
	return nil, 200
}
