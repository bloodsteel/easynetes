export interface HostRecord {
  id: number;
  hostID: string;
  hostName: string;
  hostType: '虚拟机' | '裸金属' | '云主机' | '容器';
  createdTime: string;
  status: '已上线' | '已创建' | '创建中' | '运行中' | '已下线';
}

export interface HostParams extends Partial<HostRecord> {
  current: number;
  pageSize: number;
}
