import axios, { AxiosRequestConfig } from 'axios';
import { HostRecord } from '@/types/cmdb';

export interface HostDataRecord {
	total: number,
	data: HostRecord[]
}

export function queryCmdbData(params: any) {
	// const config: AxiosRequestConfig<any> = {
	// 	params: params,
	// };
	console.log("params", params);
	// return axios.get<HostDataRecord>('/api/cmdb', config);
	return axios.get<HostDataRecord>('/api/cmdb', { params });
}
