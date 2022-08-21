/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Listeners
 * @Version: 1.0.0
 * @Date: 2021/11/24 3:11 下午
 */

package Event

// Listeners 事件监听列表
var Listeners = map[string][]Listener{
	"AppInit": []Listener{
		ListenAppInit{},
	},
}
