<template>
	<div class="container">
		<Breadcrumb :items="['资产数据', '主机']" />
		<a-card class="general-card" :title="'主机资产'">
			<a-divider style="margin-top: 0" />
			<a-row style="margin-bottom: 16px">
				<a-col :span="12">
					<a-space>
						<a-button type="primary" @click="addAsset">
							<template #icon>
								<icon-plus />
							</template>
							{{ '新建' }}
						</a-button>
					</a-space>
				</a-col>
				<a-col :span="12" style="display: flex; align-items: center; justify-content: end">
					<a-button>
						<template #icon>
							<icon-download />
						</template>
						{{ '下载' }}
					</a-button>
					<a-tooltip :content="'刷新'">
						<div class="action-icon" @click="refresh"><icon-refresh size="18" /></div>
					</a-tooltip>
					<a-dropdown @select="handleSelectDensity">
						<a-tooltip :content="'密度'">
							<div class="action-icon"><icon-line-height size="18" /></div>
						</a-tooltip>
						<template #content>
							<a-doption
								v-for="item in densityList"
								:key="item.value"
								:value="item.value"
								:class="{ active: item.value === size }"
							>
								<span>{{ item.name }}</span>
							</a-doption>
						</template>
					</a-dropdown>
				</a-col>
			</a-row>
			<!-- 表格 -->
			<a-table
				row-key="id"
				:loading="loading"
				:pagination="false"
				:columns="columns"
				:data="renderData"
				:bordered="false"
				:size="size"
				@page-change="onPageChange"
			>
				<!-- 索引 slot -->
				<template #index="{ rowIndex }">
					{{ rowIndex + 1 + (pagination.current - 1) * pagination.pageSize }}
				</template>
				<!-- 状态 slot -->
				<template #status="{ record }">
					<!--<span v-if="record.status === 'offline'" class="circle"></span>-->
					<span v-if="record.status === false" class="circle"></span>
					<span v-else class="circle pass"></span>
					{{ '已上线' }}
				</template>
				<!-- 操作 slot -->
				<template #operations="{ record }">
					<a-button type="text" size="small" @click="updateAsset(record)">编辑</a-button>
					<!-- <a-button type="text" size="small" status="success">成功</a-button>
					<a-button type="text" size="small" status="warning">警告</a-button> -->
					<a-popconfirm content="确认删除吗？" type="error" @ok="deleteAsset(record.id)">
						<a-button type="text" size="small" status="danger">删除</a-button>
					</a-popconfirm>
				</template>
			</a-table>
			<!-- 分页组件 -->
			<a-pagination
				:total="pagination.total ? pagination.total : 0"
				v-model:current="pagination.current"
				v-model:page-size="pagination.pageSize"
				show-jumper
				show-total
				show-page-size
				:page-size-options="[2, 3, 5, 10, 20, 30, 40, 50]"
				@change="onPageChange"
				@page-size-change="onPageSizeChange"
			>
			</a-pagination>
		</a-card>
		<a-drawer
			:visible="state.formVisible"
			:width="800"
			@ok="handleOk(state.formData)"
			@cancel="handleCancel"
			unmountOnClose
			:footer="true"
		>
			<template #header>
				header
				<a-space>
					<a-button type="primary" style="z-index: 2000" @cilck="handleOk(state.formData)">Submit</a-button>
				</a-space>
			</template>
			<!-- <template #title> -->
			<!-- Title -->
			<!-- 这里的按钮不生效 -->
			<!-- <a-button type="primary" style="z-index: 2000" @cilck="handleOk(form)">Submit</a-button> -->
			<!-- </template> -->

			<div>
				<!-- handleSubmit -->
				<!-- @submit="handleOk(form)" -->
				<a-form ref="formRef" :model="state.formData" :style="{ width: '600px' }">
					<a-form-item field="hostName" tooltip="hostname" label="主机名">
						<a-input v-model="state.formData.hostName" placeholder="hostName" />
					</a-form-item>
					<a-form-item field="hostIP" label="ip地址">
						<a-input v-model="state.formData.hostIP" placeholder="ip" />
					</a-form-item>
					<a-form-item field="hostType" label="主机类型">
						<a-select
							:style="{ width: '320px' }"
							v-model="state.formData.hostType"
							placeholder="Please select ..."
							allow-clear
						>
							<a-option>Linux</a-option>
							<a-option>Windows</a-option>
							<a-option>Switch</a-option>
							<a-option>Router</a-option>
						</a-select>
					</a-form-item>
					<a-form-item field="userName" label="用户名">
						<a-input v-model="state.formData.userName" placeholder="用户名" />
					</a-form-item>
					<a-form-item field="status">
						<a-checkbox v-model="state.formData.status">是否上线</a-checkbox>
					</a-form-item>
					<a-form-item field="hostSSHPort" label="ssh端口">
						<a-input-number
							v-model="state.formData.hostSSHPort"
							:style="{ width: '320px' }"
							placeholder="ssh端口"
							allow-clear
							hide-button
						></a-input-number>
					</a-form-item>
					<a-form-item field="comment" label="备注">
						<a-textarea
							v-model="state.formData.comment"
							placeholder="Please enter something"
							:max-length="{ length: 256, errorOnly: true }"
							allow-clear
							show-word-limit
						/>
					</a-form-item>
					<!-- <a-form-item>
						<a-button html-type="submit">Submit</a-button>
					</a-form-item> -->
				</a-form>
			</div>
		</a-drawer>
	</div>
</template>

<script lang="ts" setup>
	import { AxiosResponse } from 'axios';
	import { computed, ref, reactive, onMounted } from 'vue';
	import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
	import useLoading from '@/hooks/loading';
	import { Pagination } from '@/types/global';
	import { HostRecord, HostParams } from '@/types/cmdb';
	import { HostDataRecord, queryCmdbData, createCmdbData } from '@/api/cmdb';
	//
	import { HostDataToCreate, HostItemRecord, deleteCmdbData, updateCmdbData } from '@/api/cmdb';
	import { FormInstance } from '@arco-design/web-vue/es/form';
	import { Notification } from '@arco-design/web-vue';

	const defaultHost: HostRecord = {
		id: 0,
		hostName: '',
		hostIP: '',
		hostType: 'Linux',
		userName: '',
		status: true,
		comment: '',
		hostSSHPort: 22,
	};

	const state = reactive({
		formVisible: false,
		formData: { ...defaultHost },
	});
	// 数据
	const { loading, setLoading } = useLoading(true);
	// 表格密度 框架自带的值
	type SizeProps = 'mini' | 'small' | 'medium' | 'large';
	const size = ref<SizeProps>('medium');
	const handleSelectDensity = (val: string | number | Record<string, any> | undefined, e: Event) => {
		size.value = val as SizeProps;
	};
	const densityList = computed(() => [
		{
			name: '迷你',
			value: 'mini',
		},
		{
			name: '偏小',
			value: 'small',
		},
		{
			name: '中等',
			value: 'medium',
		},
		{
			name: '偏大',
			value: 'large',
		},
	]);
	// 表格分页
	const basePagination: Pagination = {
		current: 1,
		pageSize: 2,
		total: 0,
	};
	const pagination = reactive({
		...basePagination,
	});
	// 表格列描述信息
	const columns = computed<TableColumnData[]>(() => [
		{
			title: '#',
			dataIndex: 'index',
			slotName: 'index',
		},
		{
			title: '资产ID',
			dataIndex: 'hostID',
		},
		{
			title: '主机名',
			dataIndex: 'hostName',
		},
		{
			title: '主机类型',
			dataIndex: 'hostType',
		},
		{
			title: '创建时间',
			dataIndex: 'createTime',
		},
		{
			title: '主机状态',
			dataIndex: 'status',
			slotName: 'status',
		},
		{
			title: '操作',
			dataIndex: 'operations',
			slotName: 'operations',
		},
	]);
	// 获取数据
	const renderData = ref<HostRecord[]>([]);
	const fetchData = async () => {
		setLoading(true);
		// pagination
		queryCmdbData({ current: pagination.current, pageSize: pagination.pageSize })
			.then((res: HostDataRecord) => {
				// https://deepinout.com/typescript/typescript-questions/280_typescript_how_to_use_a_type_for_the_response_from_axiosget.html : TypeScript 如何使用 axios.get 的返回类型
				// 这里的例子，也不符合预期
				// https://lembdadev.com/posts/http-request-axios-typescript : 有 post 带类型的例子
				// https://bobbyhadz.com/blog/typescript-http-request-axios#making-http-get-requests-with-axios-in-typescript :
				// 这里的是使用的 try 的方式, 有一个公共的 api , 即里面请求的接口，可以直接获取到数据
				// https://axios-http.com/docs/example : 官网的例子，没有类型
				// https://upmostly.com/typescript/how-to-use-axios-in-your-typescript-apps : 有带类型的例子，不过也是 try 的形式
				console.log('the res is ', res);
				console.log('the res.data is ', res.data);
				// pagination.total = res.data.total;
				pagination.total = res.total;
				// renderData.value = res.data.data;
				renderData.value = res.data;
			})
			.catch((err) => {
				console.log('get api error', err);
			})
			.finally(() => {
				setLoading(false);
			});
	};
	// 初始化时加载数据
	onMounted(() => {
		fetchData();
		console.log('the renderData is ', renderData);
	});
	// fetchData();
	// 改变页码
	const onPageChange = (current: number) => {
		fetchData();
	};
	// 改变每页大小
	const onPageSizeChange = (pageSize: number) => {
		// 页面变化时，考虑当前页是否要变化
		// 如果从当前第 2 页， 每页 10 个，变为 每页 20 个, 那么是否要变化当前页的位置呢？
		// 考虑如果换页之后，当前页仍然有数据，那么就仍然是当前页，否则跳至第一页。
		// 实际情况是：这个函数一执行，页面总是会换到第一页, 但是实际的 pagination.current 却没变
		// 这里还是有点问题的
		if (!pagination.total) {
			pagination.total = 0;
		}
		console.log('pageSize is 1: ', pagination);
		// 如果当前页大于1
		if (pagination.current > 1) {
			// 如果总数
			// if (((pagination.current - 1) * pagination.pageSize) > pagination.total) {
			// 	pagination.current = 1
			// }
			if (pagination.total / pagination.pageSize < pagination.current) {
				pagination.current = 1;
			}
		}
		fetchData();
		console.log('pageSize is 2: ', pagination);
		// return true;
	};
	// 刷新
	const refresh = () => {
		fetchData();
	};
	// 搜索
	const search = () => {
		console.log('searching...');
	};

	// const form = reactive({
	// 	...defaultHost
	// });
	const formRef = ref<FormInstance>();
	const handleSubmit = (data: any) => {
		console.log('data is ', data);
	};
	// 表单相关的函数
	const addAsset = () => {
		console.log('addAsset');
		state.formVisible = true;
	};
	const updateAsset = (data: HostRecord) => {
		console.log('updateAsset, data is ', data);
		state.formData = data;
		state.formVisible = true;
	};
	const deleteAsset = (id: number) => {
		console.log('deleteAsset, id is ', id);
		deleteCmdbData(id)
			.then((res) => {
				// 增加删除成功的通知消息
				Notification.success({
					title: '删除成功',
					content: '',
					closable: true,
					style: { width: '500px' },
				});
				fetchData();
			})
			.catch((err) => {
				console.log(err);
			})
			.finally(() => {});
	};
	//
	const handleOk = (data: any) => {
		console.log('handleOk');
		console.log('data is', data);
		if (data.id === 0) {
			createCmdbData(data)
				.then((res) => {
					console.log('the create res is', res);
					fetchData();
					formRef.value?.resetFields();
					state.formVisible = false;
				})
				.catch((err) => {
					console.log('get api error', err);
				})
				.finally(() => {
					setLoading(false);
				});
		} else {
			updateCmdbData(data)
				.then((res) => {
					formRef.value?.resetFields();
					fetchData();
					state.formVisible = false;
				})
				.catch((err) => {
					console.log('get api error', err);
				})
				.finally(() => {
					setLoading(false);
				});
		}
	};
	const handleCancel = () => {
		console.log('handleCancel');
		formRef.value?.resetFields();
		state.formVisible = false;
		state.formData = defaultHost;
	};
</script>

<script lang="ts">
	export default {
		name: 'Host',
	};
</script>

<style scoped lang="less">
	.container {
		padding: 0 20px 20px 20px;
	}

	:deep(.arco-table-th) {
		&:last-child {
			.arco-table-th-item-title {
				margin-left: 16px;
			}
		}
	}

	.action-icon {
		margin-left: 12px;
		cursor: pointer;
	}

	.active {
		color: #0960bd;
		background-color: #e3f4fc;
	}

	.setting {
		display: flex;
		align-items: center;
		width: 200px;

		.title {
			margin-left: 12px;
			cursor: pointer;
		}
	}
</style>
