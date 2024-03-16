import axios from 'axios';
import type { RouteRecordNormalized } from 'vue-router';
import { UserState } from '@/stores/modules/user/types';

export interface LoginData {
	username: string;
	password: string;
}

export interface LoginRes {
	token: string;
}
// 请求登录接口
export function login(data: LoginData) {
	return axios.post<LoginRes>('/api/user/login', data);
}
// 请求登出接口
export function logout() {
	return axios.post<LoginRes>('/api/user/logout');
}
// 请求用户信息接口
export function getUserInfo() {
	return axios.post<UserState>('/api/user/info');
}

export function getMenuList() {
	return axios.post<RouteRecordNormalized[]>('/api/user/menu');
}
