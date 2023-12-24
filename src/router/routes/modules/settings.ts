import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const SETTINGS: AppRouteRecordRaw = {
  path: '/settings',
  name: 'settings',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: '系统设置',
    requiresAuth: true,
    icon: 'icon-tool',
    order: 4,
    hideInMenu: false,
  },
  children: [
    {
      path: 'jenkins',
      name: 'Jenkins',
      component: () => import('@/views/settings/jenkins/index.vue'),
      meta: {
        locale: 'Jenkins',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'git',
      name: 'Git',
      component: () => import('@/views/settings/gitlab/index.vue'),
      meta: {
        locale: 'Git',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default SETTINGS;
