import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const USER: AppRouteRecordRaw = {
	path: '/user',
	name: 'user',
	component: DEFAULT_LAYOUT,
	meta: {
		locale: '用户中心',
		requiresAuth: true,
		icon: 'icon-user',
		order: 5,
		hideInMenu: false,
	},
	children: [
		{
			path: 'info',
			name: 'Info',
			component: () => import('@/views/user/info/index.vue'),
			meta: {
				locale: '用户信息',
				requiresAuth: true,
				roles: ['*'],
			},
		},
		{
			path: 'setting',
			name: 'Setting',
			component: () => import('@/views/user/setting/index.vue'),
			meta: {
				locale: '用户设置',
				requiresAuth: true,
				roles: ['*'],
			},
		},
	],
};

export default USER;
