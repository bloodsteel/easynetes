import appRoutes from '../routes';

const mixinRoutes = [...appRoutes];

// 用于生成菜单
const appClientMenus = mixinRoutes.map((el) => {
  const { name, path, meta, redirect, children } = el;
  return {
    name,
    path,
    meta,
    redirect,
    children,
  };
});

export default appClientMenus;
