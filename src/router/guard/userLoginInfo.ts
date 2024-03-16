import type { Router, LocationQueryRaw } from 'vue-router';
import NProgress from 'nprogress'; // progress bar

import { useUserStore } from '@/stores';
import { isLogin } from '@/utils/auth';

// 验证登录的路由前置守卫
export default function setupUserLoginInfoGuard(router: Router) {
	router.beforeEach(async (to, from, next) => {
		NProgress.start();
		const userStore = useUserStore();
		if (isLogin()) {
			if (userStore.role) {
				next();
			} else {
				try {
					await userStore.info();
					next();
				} catch (error) {
					await userStore.logout();
					next({
						name: 'login',
						query: {
							redirect: to.name,
							...to.query,
						} as LocationQueryRaw,
					});
				}
			}
		} else {
			if (to.name === 'login') {
				next();
				return;
			}
			next({
				name: 'login',
				query: {
					redirect: to.name,
					...to.query,
				} as LocationQueryRaw,
			});
		}
	});
}
