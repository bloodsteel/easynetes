import axios, { AxiosRequestConfig } from 'axios';
import { HostRecord } from '@/types/cmdb';
import { Stringifiable } from 'query-string';

export interface HostDataRecord {
	total: number;
	data: HostRecord[];
	state: boolean;
	code: number;
}

export function queryCmdbData(params: any) {
	// const config: AxiosRequestConfig<any> = {
	// 	params: params,
	// };
	console.log('params', params);
	// return axios.get<HostDataRecord>('/api/cmdb', config);
	// return axios.get<HostDataRecord>('/api/cmdb', { params });
	// return axios.get<HostDataRecord>('http://127.0.0.1:8080/api/v1/host', { params });
	return axios.get<undefined, HostDataRecord>('/api/v1/host', { params });
}

// { "hostName": "mysql-2.dev.com", "hostIP": "192.168.1.2", "hostType": "Linux", "userName": "centos", "status": true, "comment": "test host", "hostSSHPort": 22 }

export interface HostItemRecord {
	data: HostRecord;
	state: boolean;
	code: number;
}

export interface HostDataToCreate {
	hostName: string;
	hostIP: string;
	// 这里应该是 枚举，先使用 string
	hostType: string;
	userName: string;
	// true 为上线，false 为未上线
	status: boolean;
	comment: string;
	hostSSHPort: number;
}

export function createCmdbData(data: HostDataToCreate) {
	const apiPath = '/api/v1/host';
	// post < T = any, R = AxiosResponse<T>, D = any > (url: string, data ?: D, config ?: AxiosRequestConfig<D>): Promise<R>;
	// 先不传入 config 对象
	// return axios.post<undefined, HostItemRecord, HostDataToCreate>(apiPath, data, config);
	return axios.post<undefined, HostItemRecord, HostDataToCreate>(apiPath, data);
}
export function deleteCmdbData(id: number | string) {
	const apiPath = '/api/v1/host/';
	return axios.delete(apiPath + +id.toString());
}

export function updateCmdbData(data: HostRecord) {
	const apiPath = '/api/v1/host/';
	let id = 0;
	if (data.id) {
		id = data.id;
	}
	return axios.put<undefined, HostItemRecord, HostRecord>(apiPath + id.toString(), data);
}
