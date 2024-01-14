import axios from 'axios';

import { HostRecord } from '@/types/cmdb'

export interface HostDataRecord {
    status: 'ok',
    msg: '请求成功',
    code: 20000,
    data: HostRecord[],
}
 

export function queryCmdbData() {
    return axios.get<HostDataRecord>('/api/cmdb');
  }
  