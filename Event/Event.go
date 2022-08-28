/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 事件处理
 * @File: tag
 * @Version: 1.0.0
 * @Date: 2021/11/11 10:43 上午
 */

package Event

// Listener 事件监听者,必须有Do方法
type Listener interface {
	Do(eventName string, data ...interface{}) (error, int)
}

// BindEvents 批量绑定事件
func BindEvents(eventList map[string][]Listener) {
	for eventName, ListenerList := range eventList {
		if _, ok := Listeners[eventName]; !ok {
			Listeners[eventName] = []Listener{}
		}
		for _, listener := range ListenerList {
			Listeners[eventName] = append(Listeners[eventName], listener)
		}
	}
}

// BindEvent 绑定事件
func BindEvent(eventName string, listener Listener) {
	if _, ok := Listeners[eventName]; !ok {
		Listeners[eventName] = []Listener{
			listener,
		}
	} else {
		Listeners[eventName] = append(Listeners[eventName], listener)
	}
}

// Trigger 触发一个事件
func Trigger(eventName string, data ...interface{}) (error, int) {
	lastCode := 0
	if ListenerList, ok := Listeners[eventName]; ok {
		for _, listener := range ListenerList {
			//todo:拷贝一个事件处理对象？
			err, code := listener.Do(eventName, data...)
			if err != nil {
				return err, code
			}
			lastCode = code
		}
	}
	return nil, lastCode
}
