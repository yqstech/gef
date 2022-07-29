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
	//发送短信
	"SmsSend": []Listener{
		//#Map tel(string) ip(string) content(string) template_out_id(string)
		SmsSend{}, //发送短信
	},
	//发送短信通道
	"SmsAli": []Listener{
		SmsAli{}, //阿里短信
	},
	"SmsAm": []Listener{
		SmsAm{}, //云市场
	},
	"SmsJdcx": []Listener{
		SmsJdcx{}, //京东万象
	},
	"SmsMock": []Listener{
		SmsMock{}, //模拟短信
	},
}
