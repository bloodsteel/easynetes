import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const SERVICE: AppRouteRecordRaw = {
	path: '/service',
	name: 'service',
	component: DEFAULT_LAYOUT,
	meta: {
		locale: '项目服务',
		requiresAuth: true,
		icon: 'icon-apps',
		order: 1,
		hideInMenu: false,
	},
	children: [
		{
			path: 'tree',
			name: 'Tree',
			component: () => import('@/views/service/tree/index.vue'),
			meta: {
				locale: '服务树',
				requiresAuth: true,
				roles: ['*'],
			},
		},
		{
			path: 'metaData',
			name: 'MetaData',
			component: () => import('@/views/service/meta-data/index.vue'),
			meta: {
				locale: '元数据',
				requiresAuth: true,
				roles: ['*'],
			},
		},
	],
};

export default SERVICE;
