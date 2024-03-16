import { computed } from 'vue';
import { RouteRecordRaw, RouteRecordNormalized } from 'vue-router';
import usePermission from '@/hooks/permission';
import appClientMenus from '@/router/app-menus';
import { cloneDeep } from 'lodash';

export default function useMenuTree() {
	const permission = usePermission();
	// 返回路由
	const appRoute = computed(() => {
		return appClientMenus;
	});
	// 返回一个有效的(权限判断和元数据配置等)、排序后的路由树
	const menuTree = computed(() => {
		const copyRouter = cloneDeep(appRoute.value) as RouteRecordNormalized[];
		// 根据路由元数据中的order进行排序
		copyRouter.sort((a: RouteRecordNormalized, b: RouteRecordNormalized) => {
			return (a.meta.order || 0) - (b.meta.order || 0);
		});
		// 遍历每个路由, 按层递归遍历
		function travel(_routes: RouteRecordRaw[], layer: number) {
			if (!_routes) return null;

			const collector: any = _routes.map((element) => {
				// 判断给定的路由是否有权限访问, 没权限则直接返回null
				if (!permission.accessRouter(element)) {
					return null;
				}

				// leaf node  禁止显示子菜单, 或者没有子菜单; 则直接返回, 并让element的子菜单为空
				// 表达式为真说明没有子菜单, 直接返回
				if (element.meta?.hideChildrenInMenu || !element.children) {
					element.children = [];
					return element;
				}

				// 上一步已经把没有子菜单的路由过滤掉了, 这里过滤子菜单为不隐藏的; 只要没有明确设置hideInMenu 都显示
				element.children = element.children.filter((x) => x.meta?.hideInMenu !== true);

				// 递归处理子节点，返回的是有效的(有权限访问等)子菜单
				const subItem = travel(element.children, layer + 1);

				// 经过处理后返回的子菜单，加入element
				if (subItem.length) {
					element.children = subItem;
					return element;
				}
				// 如果没有有效的子菜单, 并且当前层级为1层以上 ? TODO: 没看懂 似乎不会执行到
				if (layer > 1) {
					element.children = subItem;
					return element;
				}

				// TODO: 没看懂  似乎不会执行到
				if (element.meta?.hideInMenu === false) {
					return element;
				}

				return null;
			});
			// 将collector中为null的元素过滤掉, 保留有效的数组元素
			return collector.filter(Boolean);
		}
		return travel(copyRouter, 0);
	});

	return {
		menuTree,
	};
}
