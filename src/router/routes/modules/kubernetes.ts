import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const KUBE: AppRouteRecordRaw = {
  path: '/kubernetes',
  name: 'kubernetes',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: '容器集群',
    requiresAuth: true,
    icon: 'icon-common',
    order: 3,
    hideInMenu: false,
  },
  children: [
    {
      path: 'cluster',
      name: 'Cluster',
      component: () => import('@/views/kubernetes/cluster/index.vue'),
      meta: {
        locale: '集群管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'namespace',
      name: 'Namespace',
      component: () => import('@/views/kubernetes/namespace/index.vue'),
      meta: {
        locale: '命名空间',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'workload',
      name: 'Workload',
      component: () => import('@/views/kubernetes/workload/index.vue'),
      meta: {
        locale: '工作负载',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default KUBE;
