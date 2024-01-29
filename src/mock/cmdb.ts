import Mock from 'mockjs';
// import fs from 'fs';
import setupMock, { successResponseWrap, successResponseWrapForList } from '@/utils/setup-mock';

// import { HostRecord, HostType,HostStatus } from '@/types/cmdb'
import { HostRecord } from '@/types/cmdb'
import { throttle } from 'lodash';

const hostJson = `[{
"id": 1,
	"hostID": "host-1-1",
		"hostName": "mysql-1.dev.com",
			"hostIP": "192.168.1.1",
				"hostSSHPort": 22,
					"hostType": "裸金属",
						"createdTime": "2023-12-1 15:30:00",
							"updatedTime": "2023-12-1 15:30:00",
								"status": "在线"
},
{
	"id": 2,
		"hostID": "host-1-2",
			"hostName": "mysql-2.dev.com",
				"hostIP": "192.168.1.2",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-2 15:30:00",
								"updatedTime": "2023-12-2 15:30:00",
									"status": "在线"
},
{
	"id": 3,
		"hostID": "host-1-3",
			"hostName": "mysql-3.dev.com",
				"hostIP": "192.168.1.3",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-3 15:30:00",
								"updatedTime": "2023-12-3 15:30:00",
									"status": "在线"
},
{
	"id": 4,
		"hostID": "host-1-4",
			"hostName": "mysql-4.dev.com",
				"hostIP": "192.168.1.4",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-4 15:30:00",
								"updatedTime": "2023-12-4 15:30:00",
									"status": "在线"
},
{
	"id": 5,
		"hostID": "host-1-5",
			"hostName": "mysql-5.dev.com",
				"hostIP": "192.168.1.5",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-5 15:30:00",
								"updatedTime": "2023-12-5 15:30:00",
									"status": "在线"
},
{
	"id": 6,
		"hostID": "host-1-6",
			"hostName": "mysql-6.dev.com",
				"hostIP": "192.168.1.6",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-6 15:30:00",
								"updatedTime": "2023-12-6 15:30:00",
									"status": "在线"
},
{
	"id": 7,
		"hostID": "host-1-7",
			"hostName": "mysql-7.dev.com",
				"hostIP": "192.168.1.7",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-7 15:30:00",
								"updatedTime": "2023-12-7 15:30:00",
									"status": "在线"
},
{
	"id": 8,
		"hostID": "host-1-8",
			"hostName": "mysql-8.dev.com",
				"hostIP": "192.168.1.8",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-8 15:30:00",
								"updatedTime": "2023-12-8 15:30:00",
									"status": "在线"
},
{
	"id": 9,
		"hostID": "host-1-9",
			"hostName": "mysql-9.dev.com",
				"hostIP": "192.168.1.9",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-9 15:30:00",
								"updatedTime": "2023-12-9 15:30:00",
									"status": "在线"
},
{
	"id": 10,
		"hostID": "host-1-10",
			"hostName": "mysql-10.dev.com",
				"hostIP": "192.168.1.10",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-10 15:30:00",
								"updatedTime": "2023-12-10 15:30:00",
									"status": "在线"
},
{
	"id": 11,
		"hostID": "host-1-11",
			"hostName": "mysql-11.dev.com",
				"hostIP": "192.168.1.11",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-11 15:30:00",
								"updatedTime": "2023-12-11 15:30:00",
									"status": "在线"
},
{
	"id": 12,
		"hostID": "host-1-12",
			"hostName": "mysql-12.dev.com",
				"hostIP": "192.168.1.12",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-12 15:30:00",
								"updatedTime": "2023-12-12 15:30:00",
									"status": "在线"
},
{
	"id": 13,
		"hostID": "host-1-13",
			"hostName": "mysql-13.dev.com",
				"hostIP": "192.168.1.13",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-13 15:30:00",
								"updatedTime": "2023-12-13 15:30:00",
									"status": "在线"
},
{
	"id": 14,
		"hostID": "host-1-14",
			"hostName": "mysql-14.dev.com",
				"hostIP": "192.168.1.14",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-14 15:30:00",
								"updatedTime": "2023-12-14 15:30:00",
									"status": "在线"
},
{
	"id": 15,
		"hostID": "host-1-15",
			"hostName": "mysql-15.dev.com",
				"hostIP": "192.168.1.15",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-15 15:30:00",
								"updatedTime": "2023-12-15 15:30:00",
									"status": "在线"
},
{
	"id": 16,
		"hostID": "host-1-16",
			"hostName": "mysql-16.dev.com",
				"hostIP": "192.168.1.16",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-16 15:30:00",
								"updatedTime": "2023-12-16 15:30:00",
									"status": "在线"
},
{
	"id": 17,
		"hostID": "host-1-17",
			"hostName": "mysql-17.dev.com",
				"hostIP": "192.168.1.17",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-17 15:30:00",
								"updatedTime": "2023-12-17 15:30:00",
									"status": "在线"
},
{
	"id": 18,
		"hostID": "host-1-18",
			"hostName": "mysql-18.dev.com",
				"hostIP": "192.168.1.18",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-18 15:30:00",
								"updatedTime": "2023-12-18 15:30:00",
									"status": "在线"
},
{
	"id": 19,
		"hostID": "host-1-19",
			"hostName": "mysql-19.dev.com",
				"hostIP": "192.168.1.19",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-19 15:30:00",
								"updatedTime": "2023-12-19 15:30:00",
									"status": "在线"
},
{
	"id": 20,
		"hostID": "host-1-20",
			"hostName": "mysql-20.dev.com",
				"hostIP": "192.168.1.20",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-20 15:30:00",
								"updatedTime": "2023-12-20 15:30:00",
									"status": "在线"
},
{
	"id": 21,
		"hostID": "host-1-21",
			"hostName": "mysql-21.dev.com",
				"hostIP": "192.168.1.21",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-21 15:30:00",
								"updatedTime": "2023-12-21 15:30:00",
									"status": "在线"
},
{
	"id": 22,
		"hostID": "host-1-22",
			"hostName": "mysql-22.dev.com",
				"hostIP": "192.168.1.22",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-22 15:30:00",
								"updatedTime": "2023-12-22 15:30:00",
									"status": "在线"
},
{
	"id": 23,
		"hostID": "host-1-23",
			"hostName": "mysql-23.dev.com",
				"hostIP": "192.168.1.23",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-23 15:30:00",
								"updatedTime": "2023-12-23 15:30:00",
									"status": "在线"
},
{
	"id": 24,
		"hostID": "host-1-24",
			"hostName": "mysql-24.dev.com",
				"hostIP": "192.168.1.24",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-24 15:30:00",
								"updatedTime": "2023-12-24 15:30:00",
									"status": "在线"
},
{
	"id": 25,
		"hostID": "host-1-25",
			"hostName": "mysql-25.dev.com",
				"hostIP": "192.168.1.25",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-25 15:30:00",
								"updatedTime": "2023-12-25 15:30:00",
									"status": "在线"
},
{
	"id": 26,
		"hostID": "host-1-26",
			"hostName": "mysql-26.dev.com",
				"hostIP": "192.168.1.26",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-26 15:30:00",
								"updatedTime": "2023-12-26 15:30:00",
									"status": "在线"
},
{
	"id": 27,
		"hostID": "host-1-27",
			"hostName": "mysql-27.dev.com",
				"hostIP": "192.168.1.27",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-27 15:30:00",
								"updatedTime": "2023-12-27 15:30:00",
									"status": "在线"
},
{
	"id": 28,
		"hostID": "host-1-28",
			"hostName": "mysql-28.dev.com",
				"hostIP": "192.168.1.28",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-28 15:30:00",
								"updatedTime": "2023-12-28 15:30:00",
									"status": "在线"
},
{
	"id": 29,
		"hostID": "host-1-29",
			"hostName": "mysql-29.dev.com",
				"hostIP": "192.168.1.29",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-29 15:30:00",
								"updatedTime": "2023-12-29 15:30:00",
									"status": "在线"
},
{
	"id": 30,
		"hostID": "host-1-30",
			"hostName": "mysql-30.dev.com",
				"hostIP": "192.168.1.30",
					"hostSSHPort": 22,
						"hostType": "裸金属",
							"createdTime": "2023-12-30 15:30:00",
								"updatedTime": "2023-12-30 15:30:00",
									"status": "在线"
}
]`

const hostJsonData: HostRecord[] = JSON.parse(hostJson)

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
	},
	{
		id: 2,
		hostID: "host-1-2",
		hostName: "mysql-2.dev.com",
		hostIP: "192.168.1.2",
		hostSSHPort: 22,
		hostType: "裸金属",
		createdTime: "2023-12-2 15: 30:00",
		updatedTime: "2023-12-2 15: 30:00",
		status: "在线"
	},
	{
		id: 3,
		hostID: "host-1-3",
		hostName: "mysql-3.dev.com",
		hostIP: "192.168.1.3",
		hostSSHPort: 22,
		hostType: "裸金属",
		createdTime: "2023-12-3 15:30:00",
		updatedTime: "2023-12-3 15:30:00",
		status: "在线"
	},
]
// const jsonStr = fs.readFileSync('cmdb/hosts.json', 'utf-8')
// const jsonData = JSON.parse(jsonStr);
// cmdbHosts.concat(jsonData);


setupMock({
	setup() {
		Mock.mock(new RegExp('/api/cmdb'), (params: {
			url: string
		}) => {
			// const pageSize = params.pageSize;
			const urlArray = params.url.split('?');
			const paramArr = urlArray[1].split('&');
			const paramObj: Record<string, string> = {};
			paramArr.forEach(item => {
				const theKVArr = item.split("=");
				const k = theKVArr[0];
				const v = theKVArr[1];
				paramObj[k] = v;
			})
			console.log("/api/cmdb", paramObj);
			const current = Number(paramObj.current);
			const pageSize = Number(paramObj.pageSize);
			let theOffset = 0;
			let sliceEnd = pageSize;
			if (current > 1) {
				theOffset = (current - 1) * pageSize;
				sliceEnd = current * pageSize;
			}
			return successResponseWrap({
				total: hostJsonData.length,
				data: hostJsonData.slice(theOffset, sliceEnd),
			})
			// return successResponseWrapForList(cmdbHosts.length, cmdbHosts);
		});
	}
});
