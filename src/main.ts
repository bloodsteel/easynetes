import '@/assets/style/global.less';
import '@/api/interceptor';

import { createApp } from 'vue';

import ArcoVue from '@arco-design/web-vue';
import ArcoVueIcon from '@arco-design/web-vue/es/icon';

import App from '@/App.vue';
import router from '@/router';
import stores from '@/stores';

import globalComponents from '@/components';

import '@/mock';

const app = createApp(App);

app.use(ArcoVue, {});
app.use(ArcoVueIcon);

app.use(router);
app.use(stores);
app.use(globalComponents);

app.mount('#app');
