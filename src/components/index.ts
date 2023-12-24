import { App } from 'vue';
import { use } from 'echarts/core';
// 渲染
import { CanvasRenderer } from 'echarts/renderers';
// 图表
import { BarChart, LineChart, PieChart, RadarChart } from 'echarts/charts';
// 基础组件
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent,
  GraphicComponent,
} from 'echarts/components';
import Chart from './chart/index.vue';
import Breadcrumb from './breadcrumb/index.vue';

// 这里是手动导入组件和图表，保证更小的体积
// https://github.com/ecomfe/vue-echarts#example
use([
  CanvasRenderer,
  BarChart,
  LineChart,
  PieChart,
  RadarChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent,
  GraphicComponent,
]);

export default {
  install(Vue: App) {
    // 注册全局组件
    Vue.component('Chart', Chart);
    Vue.component('Breadcrumb', Breadcrumb);
  },
};
