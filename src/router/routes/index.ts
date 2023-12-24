import type { RouteRecordNormalized } from 'vue-router';

const modules = import.meta.glob('./modules/*.ts', { eager: true });

// 将所有模块中的export default 整合到一个路由数组中
function formatModules(_modules: any, result: RouteRecordNormalized[]) {
  // Object.keys 遍历 对象中的所有key, 返回数组
  Object.keys(_modules).forEach((key) => {
    // 这里表示获取模块中的export default
    // 说明在配置路由和菜单的时候, 必须export default
    const defaultModule = _modules[key].default;
    if (!defaultModule) return;
    const moduleList = Array.isArray(defaultModule)
      ? [...defaultModule]
      : [defaultModule];
    result.push(...moduleList);
  });
  return result;
}

const appRoutes: RouteRecordNormalized[] = formatModules(modules, []);

export default appRoutes;
