/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-10-20 16:38:23
 * @LastEditTime: 2021-10-20 16:38:24
 * @Description:
 */

package util

import (
	"math/rand"
	"time"
)

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// String2Bool 字符串转布尔
func String2Bool(str string) bool {
	if str == "false" || str == "" {
		return false
	} else {
		return true
	}
}
