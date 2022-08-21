/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: math.go
 * @Version: 1.0.0
 * @Date: 2022/1/10 11:55 下午
 */

package util

import "math"

// FloatRound2Int 浮点数四舍五入变整数
// 例如：util.FloatRound2Int(float64(990*81)/100.0) ==> 802
func FloatRound2Int(x float64) int64 {
	return int64(math.Floor(x + 0.5))
}
