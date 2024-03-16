<template>
	<div class="navbar">
		<!-- 页面顶部最左侧 -->
		<div class="left-side">
			<a-space>
				<img
					alt="logo"
					src="//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/dfdba5317c0c20ce20e64fac803d52bc.svg~tplv-49unhts6dw-image.image"
				/>
				<a-typography-title :style="{ margin: 0, fontSize: '18px' }" :heading="5">
					Easynetes Newbee
				</a-typography-title>
				<icon-menu-fold
					v-if="!topMenu && appStore.device === 'mobile'"
					style="font-size: 22px; cursor: pointer"
					@click="toggleDrawerMenu"
				/>
			</a-space>
		</div>
		<!-- 页面顶部最中间 -->
		<div class="center-side">
			<!-- 如果设置为顶部菜单栏, 这会显示 -->
			<Menu v-if="topMenu" />
		</div>
		<!-- 页面顶部最右侧 -->
		<ul class="right-side">
			<li>
				<!-- 切换主题 -->
				<a-tooltip :content="theme === 'light' ? '点击切换为暗黑模式' : '点击切换为亮色模式'">
					<a-button class="nav-btn" type="outline" :shape="'circle'" @click="handleToggleTheme">
						<template #icon>
							<icon-moon-fill v-if="theme === 'dark'" />
							<icon-sun-fill v-else />
						</template>
					</a-button>
				</a-tooltip>
			</li>
			<li>
				<!-- 全屏切换 -->
				<a-tooltip :content="isFullscreen ? '点击退出全屏模式' : '点击切换全屏模式'">
					<a-button class="nav-btn" type="outline" :shape="'circle'" @click="toggleFullScreen">
						<template #icon>
							<icon-fullscreen-exit v-if="isFullscreen" />
							<icon-fullscreen v-else />
						</template>
					</a-button>
				</a-tooltip>
			</li>
			<li>
				<!-- 用户图标 -->
				<a-dropdown trigger="click">
					<!-- 头像 https://arco.design/vue/component/avatar -->
					<a-avatar :size="32" :style="{ marginRight: '8px', cursor: 'pointer' }">
						<img alt="avatar" :src="avatar" />
					</a-avatar>
					<template #content>
						<a-doption>
							<a-space @click="switchRoles">
								<icon-tag />
								<span>
									{{ '切换角色' }}
								</span>
							</a-space>
						</a-doption>
						<a-doption>
							<a-space @click="$router.push({ name: 'Info' })">
								<icon-user />
								<span>
									{{ '用户中心' }}
								</span>
							</a-space>
						</a-doption>
						<a-doption>
							<a-space @click="$router.push({ name: 'Setting' })">
								<icon-settings />
								<span>
									{{ '用户设置' }}
								</span>
							</a-space>
						</a-doption>
						<a-doption>
							<a-space @click="handleLogout">
								<icon-export />
								<span>
									{{ '登出登录' }}
								</span>
							</a-space>
						</a-doption>
					</template>
				</a-dropdown>
			</li>
		</ul>
	</div>
</template>

<script lang="ts" setup>
	import { computed, inject } from 'vue';
	import { Message } from '@arco-design/web-vue';
	import { useDark, useToggle, useFullscreen } from '@vueuse/core';
	import { useAppStore, useUserStore } from '@/stores';
	import useUser from '@/hooks/user';
	// 菜单组件
	import Menu from '@/components/menu/index.vue';

	const appStore = useAppStore();
	const userStore = useUserStore();
	const { logout } = useUser();
	// isFullscreen 布尔值   toggle 切换函数
	const { isFullscreen, toggle: toggleFullScreen } = useFullscreen();
	// 用户头像
	const avatar = computed(() => {
		return userStore.avatar;
	});
	// 用于在菜单栏顶部的 主题切换
	const theme = computed(() => {
		return appStore.theme;
	});
	// 是否开启顶部菜单栏
	const topMenu = computed(() => appStore.topMenu && appStore.menu);
	// 暗黑模式切换
	const isDark = useDark({
		selector: 'body',
		attribute: 'arco-theme',
		valueDark: 'dark',
		valueLight: 'light',
		storageKey: 'arco-theme',
		onChanged(dark: boolean) {
			appStore.toggleTheme(dark);
		},
	});
	// useToggle函数的参数是布尔值, 返回一个切换函数
	const toggleTheme = useToggle(isDark);
	// 事件回调，切换主题
	const handleToggleTheme = () => {
		toggleTheme();
	};
	// 事件回调, 控制页面显示的抽屉页面是否显示; 只需要修改store中的值即可
	// TODO: 是不是可以删了
	// const setVisible = () => {
	//   appStore.updateSettings({ globalSettings: true });
	// };
	// const refBtn = ref();
	// TODO: 是不是可以删了
	// const triggerBtn = ref();
	// 当点击顶部菜单栏的通知按钮的时候, 通过编程的方式触发气泡卡片
	// TODO: 是不是可以删了
	/*
  const setPopoverVisible = () => {
    // 原生ES语法, 创建一个鼠标点击事件
    // https://developer.mozilla.org/en-US/docs/Web/API/MouseEvent
    // TODO: 有没有更好的文档
    const event = new MouseEvent('click', {
      view: window,
      bubbles: true,
      cancelable: true,
    });
    // TODO: 是给这个按钮触发这个事件？
    refBtn.value.dispatchEvent(event);
  };
  */
	// 登出按钮回调
	const handleLogout = () => {
		logout();
	};
	// 用于切换角色
	const switchRoles = async () => {
		const res = await userStore.switchRoles();
		Message.success(res as string);
	};
	// 依赖注入, 这个是从layout里面的provide提供的
	const toggleDrawerMenu = inject('toggleDrawerMenu') as () => void;
</script>

<style scoped lang="less">
	.navbar {
		display: flex;
		justify-content: space-between;
		height: 100%;
		background-color: var(--color-bg-2);
		border-bottom: 1px solid var(--color-border);
	}

	.left-side {
		display: flex;
		align-items: center;
		padding-left: 20px;
	}

	.center-side {
		flex: 1;
	}

	.right-side {
		display: flex;
		padding-right: 20px;
		list-style: none;
		:deep(.locale-select) {
			border-radius: 20px;
		}
		li {
			display: flex;
			align-items: center;
			padding: 0 10px;
		}

		a {
			color: var(--color-text-1);
			text-decoration: none;
		}
		.nav-btn {
			border-color: rgb(var(--gray-2));
			color: rgb(var(--gray-8));
			font-size: 16px;
		}
		.trigger-btn,
		.ref-btn {
			position: absolute;
			bottom: 14px;
		}
		.trigger-btn {
			margin-left: 14px;
		}
	}
</style>

<style lang="less">
	.message-popover {
		.arco-popover-content {
			margin-top: 0;
		}
	}
</style>
