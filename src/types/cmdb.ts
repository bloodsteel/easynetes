// export enum hostType {
//   // 虚拟机
//   VM = '虚拟机',
//   // Bare Metal
//   BM = '裸金属',
//   // 云主机
//   CloudHost = '云主机',
//   // 容器
//   Container = '容器',
// }

// export enum hostStatus {
//   online = '已上线',
//   created = '已创建',
//   creating = '创建中',
//   running = '运行中',
//   offline = '已下线',
// }

export interface HostRecord {
	id?: number;
	hostID?: string;
	hostName: string;
	hostIP: string;
	userName: string;
	hostSSHPort: number;
	// hostType: HostType;
	hostType: string;
	createdTime?: string;
	updatedTime?: string;
	// status: HostStatus;
	status: boolean;
	comment?: string;
}

export interface HostParams extends Partial<HostRecord> {
	current: number;
	pageSize: number;
}
