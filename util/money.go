/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: money
 * @Version: 1.0.0
 * @Date: 2022/3/2 10:55 上午
 */

package util

// Money 分(int64)转元(string)
func Money(m int64) string {
	return Float642String(float64(m) / 100.0)
}

// Money2Cent 元转分
func Money2Cent(money string) (int64, error) {
	return FloatString2Int64(money, 100.0)
}
