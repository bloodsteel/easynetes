<template>
  <!-- 整体页面布局 -->
  <a-layout class="layout" :class="{ mobile: appStore.hideMenu }">
    <!-- 顶部header -->
    <div v-if="navbar" class="layout-navbar">
      <NavBar />
    </div>
    <a-layout>
      <a-layout>
        <!-- 侧边栏 -->
        <a-layout-sider
          v-if="renderMenu"
          v-show="!hideMenu"
          class="layout-sider"
          breakpoint="xl"
          :collapsed="collapsed"
          :collapsible="true"
          :width="menuWidth"
          :style="{ paddingTop: navbar ? '60px' : '' }"
          :hide-trigger="true"
          @collapse="setCollapsed"
        >
          <div class="menu-wrapper">
            <Menu />
          </div>
        </a-layout-sider>
        <!-- navbar toggleDrawerMenu -->
        <a-drawer
          v-if="hideMenu"
          :visible="drawerVisible"
          placement="left"
          :footer="false"
          mask-closable
          :closable="false"
          @cancel="drawerCancel"
        >
          <Menu />
        </a-drawer>
        <a-layout class="layout-content" :style="paddingStyle">
          <a-layout-content>
            <!-- 页面内容 -->
            <PageLayout />
          </a-layout-content>
          <!-- 页脚 -->
          <Footer v-if="footer" />
        </a-layout>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" setup>
  import { ref, computed, watch, onMounted, provide } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { useAppStore, useUserStore } from '@/stores';
  import NavBar from '@/components/navbar/index.vue';
  import Menu from '@/components/menu/index.vue';
  import Footer from '@/components/footer/index.vue';
  import usePermission from '@/hooks/permission';
  import useResponsive from '@/hooks/responsive';
  import PageLayout from './page-layout.vue';

  const isInit = ref(false);
  const appStore = useAppStore();
  const userStore = useUserStore();
  const router = useRouter();
  const route = useRoute();
  const permission = usePermission();
  useResponsive(true);
  // 页面顶部的元素高度, 用于给页面内容设置pending-top, 当关闭页面顶部显示的时候有用
  const navbarHeight = `60px`;
  // 页面顶部
  const navbar = computed(() => appStore.navbar);
  // 是否渲染侧边栏菜单
  const renderMenu = computed(() => appStore.menu && !appStore.topMenu);
  // 是否隐藏侧边栏菜单
  const hideMenu = computed(() => appStore.hideMenu);
  // 是否显示页脚
  const footer = computed(() => appStore.footer);
  // 侧边栏宽度
  const menuWidth = computed(() => {
    // 如果侧边栏收起那么为48, 展开的话为配置中的220
    return appStore.menuCollapse ? 48 : appStore.menuWidth;
  });
  // 侧边栏收起状态
  const collapsed = computed(() => {
    return appStore.menuCollapse;
  });
  // 内容页设置CSS样式;  当关掉页面顶部或者是页面侧边栏的时候, 撑开页面内容
  const paddingStyle = computed(() => {
    const paddingLeft =
      renderMenu.value && !hideMenu.value
        ? { paddingLeft: `${menuWidth.value}px` }
        : {};
    const paddingTop = navbar.value ? { paddingTop: navbarHeight } : {};
    return { ...paddingLeft, ...paddingTop };
  });
  // 侧边栏 展开-收起时的事件回调函数; 有 人为点击触发 以及 响应式自动触发两种
  const setCollapsed = (val: boolean) => {
    if (!isInit.value) return;
    appStore.updateSettings({ menuCollapse: val });
  };
  // 当用户的角色发生变化的时候, 判断是否有该路由的权限, 如果没有则跳转到404
  watch(
    () => userStore.role,
    (roleValue) => {
      if (roleValue && !permission.accessRouter(route))
        router.push({ name: 'notFound' });
    },
  );
  // 这里是在navbar中的toggleDrawerMenu需要使用
  const drawerVisible = ref(false);
  const drawerCancel = () => {
    drawerVisible.value = false;
  };
  provide('toggleDrawerMenu', () => {
    drawerVisible.value = !drawerVisible.value;
  });

  onMounted(() => {
    isInit.value = true;
  });
</script>

<style scoped lang="less">
  @nav-size-height: 60px;
  @layout-max-width: 1100px;

  .layout {
    width: 100%;
    height: 100%;
  }

  .layout-navbar {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 100;
    width: 100%;
    height: @nav-size-height;
  }

  .layout-sider {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 99;
    height: 100%;
    transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
    &::after {
      position: absolute;
      top: 0;
      right: -1px;
      display: block;
      width: 1px;
      height: 100%;
      background-color: var(--color-border);
      content: '';
    }

    > :deep(.arco-layout-sider-children) {
      overflow-y: hidden;
    }
  }

  .menu-wrapper {
    height: 100%;
    overflow: auto;
    overflow-x: hidden;
    :deep(.arco-menu) {
      ::-webkit-scrollbar {
        width: 12px;
        height: 4px;
      }

      ::-webkit-scrollbar-thumb {
        border: 4px solid transparent;
        background-clip: padding-box;
        border-radius: 7px;
        background-color: var(--color-text-4);
      }

      ::-webkit-scrollbar-thumb:hover {
        background-color: var(--color-text-3);
      }
    }
  }

  .layout-content {
    min-height: 100vh;
    overflow-y: hidden;
    background-color: var(--color-fill-2);
    transition: padding 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
  }
</style>
