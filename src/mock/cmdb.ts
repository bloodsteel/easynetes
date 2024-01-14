import Mock from 'mockjs';
import setupMock, { successResponseWrap } from '@/utils/setup-mock';

// import { HostRecord, HostType,HostStatus } from '@/types/cmdb'
import { HostRecord } from '@/types/cmdb'

const cmdbHosts: HostRecord[] = [
    {
        id: 1,
        hostID: 'host-1-1',
        hostName: 'mysql-01.dev.com',
        hostIP: '192.168.1.1',
        hostSSHPort: 22,
        hostType: '裸金属',
        createdTime: '2023-12-23 15:30:00',
        updatedTime: '2023-12-23 15:30:00',
        status: '在线',
    }
]

setupMock({
    setup() {
        Mock.mock(new RegExp('/api/cmdb'), () => {
            console.log("/api/cmdb", );
            return successResponseWrap({
                status: '0',
                msg: '',
                code: 20000,
                data: cmdbHosts
            })
        });
    }
});
