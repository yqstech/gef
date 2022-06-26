/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: map
 * @Version: 1.0.0
 * @Date: 2022/5/21 17:56
 */

package util

func MapStringValueDef(data map[string]string, key string, defValue string) string {
	if v, ok := data[key]; !ok || v == "" {
		return defValue
	} else {
		return v
	}
}
