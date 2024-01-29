import { number } from "echarts";

export default ({ mock, setup }: { mock?: boolean; setup: () => void }) => {
  if (mock !== false && true) setup();
};

// 统一数据返回
export const successResponseWrap = (data: unknown) => {
  return {
    data,
    status: 'ok',
    msg: '请求成功',
    code: 20000,
  };
};
// 统一返回列表数据
export const successResponseWrapForList = (total: number, data: unknown) => {
  return {
    data,
    total,
    status: 'ok',
    msg: '请求成功',
    code: 20000,
  };
};

export const failResponseWrap = (data: unknown, msg: string, code = 50000) => {
  return {
    data,
    status: 'fail',
    msg,
    code,
  };
};
