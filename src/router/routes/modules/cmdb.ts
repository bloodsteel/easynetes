import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const CMDB: AppRouteRecordRaw = {
	path: '/cmdb',
	name: 'cmdb',
	component: DEFAULT_LAYOUT,
	meta: {
		locale: '主机资产',
		requiresAuth: true,
		icon: 'icon-storage',
		order: 2,
		hideInMenu: false,
	},
	children: [
		{
			path: 'zone',
			name: 'Zone',
			component: () => import('@/views/cmdb/zone/index.vue'),
			meta: {
				locale: '可用区',
				requiresAuth: true,
				roles: ['*'],
			},
		},
		{
			path: 'host',
			name: 'Host',
			component: () => import('@/views/cmdb/host/index.vue'),
			meta: {
				locale: '主机',
				requiresAuth: true,
				roles: ['*'],
			},
		},
	],
};

export default CMDB;
