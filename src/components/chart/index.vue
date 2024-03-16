<template>
	<!-- https://github.com/ecomfe/vue-echarts#example -->
	<!-- option 是ECharts的通用接口, 修改这个属性会触发ECharts的setOption方法  -->
	<!-- autoResize 布尔值, 当调整根的大小时, 是否应自动调整图标大小 -->
	<v-charts v-if="renderChart" :option="options" :autoresize="autoResize" :style="{ width, height }" />
</template>

<script lang="ts" setup>
	import { ref, nextTick } from 'vue';
	import VCharts from 'vue-echarts';
	// import { useAppStore } from '@/store';

	// 从父组件接收数据
	defineProps({
		options: {
			type: Object,
			default() {
				return {};
			},
		},
		autoResize: {
			type: Boolean,
			default: true,
		},
		width: {
			type: String,
			default: '100%',
		},
		height: {
			type: String,
			default: '100%',
		},
	});
	// const appStore = useAppStore();
	// const theme = computed(() => {
	//   if (appStore.theme === 'dark') return 'dark';
	//   return '';
	// });
	const renderChart = ref(false);
	// wait container expand  https://cn.vuejs.org/api/general.html#nexttick
	// TODO: 不理解
	nextTick(() => {
		renderChart.value = true;
	});
</script>

<style scoped lang="less"></style>
