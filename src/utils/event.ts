// 给指定的元素添加事件监听, event是事件类型, handler是事件回调函数
export function addEventListen(
	target: Window | HTMLElement,
	event: string,
	handler: EventListenerOrEventListenerObject,
	capture = false,
) {
	if (target.addEventListener && typeof target.addEventListener === 'function') {
		// https://developer.mozilla.org/zh-CN/docs/Web/API/EventTarget/addEventListener
		target.addEventListener(event, handler, capture);
	}
}

// 给指定的元素删除事件监听是addEventListen的反向操作
export function removeEventListen(
	target: Window | HTMLElement,
	event: string,
	handler: EventListenerOrEventListenerObject,
	capture = false,
) {
	if (target.removeEventListener && typeof target.removeEventListener === 'function') {
		// https://developer.mozilla.org/zh-CN/docs/Web/API/EventTarget/removeEventListener
		target.removeEventListener(event, handler, capture);
	}
}
