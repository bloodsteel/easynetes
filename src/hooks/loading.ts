import { ref } from 'vue';

// 几乎所有的地方都用到了
// 表示 元素 是否 在加载中状态
// https://arco.design/vue/component/button#API 比如按钮的加载中状态表示已经点击
// https://arco.design/vue/component/table#API  比如表格的状态是加载中
// https://arco.design/vue/component/card  比如卡片是否加载中
// https://arco.design/vue/component/list#API  列表加载中
// 等等。。
export default function useLoading(initValue = false) {
	const loading = ref(initValue);
	const setLoading = (value: boolean) => {
		loading.value = value;
	};
	const toggle = () => {
		loading.value = !loading.value;
	};
	return {
		loading,
		setLoading,
		toggle,
	};
}
