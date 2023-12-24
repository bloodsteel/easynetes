import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';

import { useUserStore } from '@/stores';

// 用户登出 只有在右上角头像中的用户退出中使用
export default function useUser() {
  // 在组合式API中获取路由器
  const router = useRouter();
  const userStore = useUserStore();
  const logout = async (logoutTo?: string) => {
    // 登出
    await userStore.logout();
    // 获取当前的路由
    const currentRoute = router.currentRoute.value;
    Message.success('登出成功');
    router.push({
      name: logoutTo && typeof logoutTo === 'string' ? logoutTo : 'login',
      query: {
        ...router.currentRoute.value.query,
        redirect: currentRoute.name as string,
      },
    });
  };
  return {
    logout,
  };
}
