<script lang="tsx">
	import { defineComponent, ref, h, compile, computed } from 'vue';
	import { useRoute, useRouter, RouteRecordRaw } from 'vue-router';
	import type { RouteMeta } from 'vue-router';
	import { useAppStore } from '@/stores';
	import { listenerRouteChange } from '@/utils/route-listener';
	import { openWindow, regexUrl } from '@/utils';
	import useMenuTree from './use-menu-tree';

	// 定义组件 https://cn.vuejs.org/guide/typescript/overview.html#definecomponent
	export default defineComponent({
		emit: ['collapse'],
		setup() {
			const appStore = useAppStore();
			const router = useRouter();
			const route = useRoute();
			const { menuTree } = useMenuTree();
			// 菜单折叠状态, 表示菜单是否折叠
			const collapsed = computed({
				get() {
					if (appStore.device === 'desktop') return appStore.menuCollapse;
					return false;
				},
				set(value: boolean) {
					appStore.updateSettings({ menuCollapse: value });
				},
			});
			// 是否是顶部菜单
			const topMenu = computed(() => appStore.topMenu);
			// 表示展开的子菜单的key数组, 用于保存记录, 当菜单折叠和打开的时候菜单状态一致
			const openKeys = ref<string[]>([]);
			// 表示选中的菜单的key数组，用于保存记录，当菜单折叠和打开的时候菜单状态一致
			const selectedKey = ref<string[]>([]);

			const goto = (item: RouteRecordRaw) => {
				// 打开外部连接, 首先判断URL是否符合正则表达式
				if (regexUrl.test(item.path)) {
					openWindow(item.path);
					selectedKey.value = [item.name as string];
					return;
				}
				// 消除外部连接的副作用
				const { hideInMenu, activeMenu } = item.meta as RouteMeta;
				if (route.name === item.name && !hideInMenu && !activeMenu) {
					selectedKey.value = [item.name as string];
					return;
				}
				router.push({
					name: item.name,
				});
			};
			const findMenuOpenKeys = (target: string) => {
				const result: string[] = [];
				let isFind = false;
				const backtrack = (item: RouteRecordRaw, keys: string[]) => {
					if (item.name === target) {
						isFind = true;
						result.push(...keys);
						return;
					}
					if (item.children?.length) {
						item.children.forEach((el) => {
							backtrack(el, [...keys, el.name as string]);
						});
					}
				};
				menuTree.value.forEach((el: RouteRecordRaw) => {
					if (isFind) return;
					backtrack(el, [el.name as string]);
				});
				return result;
			};
			//
			listenerRouteChange((newRoute) => {
				const { requiresAuth, activeMenu, hideInMenu } = newRoute.meta;
				if (requiresAuth && (!hideInMenu || activeMenu)) {
					const menuOpenKeys = findMenuOpenKeys((activeMenu || newRoute.name) as string);

					const keySet = new Set([...menuOpenKeys, ...openKeys.value]);
					openKeys.value = [...keySet];

					selectedKey.value = [activeMenu || menuOpenKeys[menuOpenKeys.length - 1]];
				}
			}, true);
			const setCollapse = (val: boolean) => {
				if (appStore.device === 'desktop') appStore.updateSettings({ menuCollapse: val });
			};

			// 遍历渲染菜单
			const renderSubMenu = () => {
				function travel(_route: RouteRecordRaw[], nodes = []) {
					if (_route) {
						_route.forEach((element) => {
							// icon 表示 菜单图标
							const icon = element?.meta?.icon ? () => h(compile(`<${element?.meta?.icon}/>`)) : null;
							// node 是 菜单节点; 如果有子菜单, 那么渲染成a-sub-menu , 否则 渲染成 a-menu-item
							const node =
								element?.children && element?.children.length !== 0 ? (
									<a-sub-menu
										key={element?.name}
										v-slots={{
											icon,
											title: () => h(compile(element?.meta?.locale || '')),
										}}
									>
										{travel(element?.children)}
									</a-sub-menu>
								) : (
									<a-menu-item key={element?.name} v-slots={{ icon }} onClick={() => goto(element)}>
										{element?.meta?.locale || ''}
									</a-menu-item>
								);
							nodes.push(node as never);
						});
					}
					return nodes;
				}
				return travel(menuTree.value);
			};

			return () => (
				<a-menu
					mode={topMenu.value ? 'horizontal' : 'vertical'}
					v-model:collapsed={collapsed.value}
					v-model:open-keys={openKeys.value}
					show-collapse-button={appStore.device !== 'mobile'}
					auto-open={false}
					selected-keys={selectedKey.value}
					auto-open-selected={true}
					level-indent={34}
					style="height: 100%;width:100%;"
					onCollapse={setCollapse}
				>
					{renderSubMenu()}
				</a-menu>
			);
		},
	});
</script>

<style lang="less" scoped>
	:deep(.arco-menu-inner) {
		.arco-menu-inline-header {
			display: flex;
			align-items: center;
		}
		.arco-icon {
			&:not(.arco-icon-down) {
				font-size: 18px;
			}
		}
	}
</style>
