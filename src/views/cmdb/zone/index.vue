<template>
	<div class="container">
		<Breadcrumb :items="['资产数据', '可用区']" />
		<a-card class="general-card" :title="'主机资产'">
			<a-divider style="margin-top: 0" />
			<a-row style="margin-bottom: 16px">
				<a-col :span="12">
					<a-space>
						<a-button type="primary">
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
						<div class="action-icon" @click="search"><icon-refresh size="18" /></div>
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
				:pagination="pagination"
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
				<!-- 操作 slot -->
				<template #operations>
					<a-button type="text" size="small">一般</a-button>
					<a-button type="text" size="small" status="success">成功</a-button>
					<a-button type="text" size="small" status="warning">警告</a-button>
					<a-button type="text" size="small" status="danger">危险</a-button>
				</template>
			</a-table>
		</a-card>
	</div>
</template>

<script lang="ts" setup>
	import { computed, ref, reactive } from 'vue';
	import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
	import useLoading from '@/hooks/loading';
	import { Pagination } from '@/types/global';

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
			title: '地域',
			dataIndex: 'region',
		},
		{
			title: '可用区',
			dataIndex: 'zone',
		},
		{
			title: '创建时间',
			dataIndex: 'createdTime',
		},
		{
			title: '备注',
			dataIndex: 'remark',
		},
		{
			title: '操作',
			dataIndex: 'operations',
			slotName: 'operations',
		},
	]);
	// 获取数据  TODO: 这里没搞通
	const renderData = ref([]);
	const fetchData = async () => {
		setLoading(true);
		try {
			renderData.value = [
				{
					id: 0,
					region: 'aliyun-huabei2',
					zone: 'A',
					createdTime: '2024-1-1',
					remark: '阿里云华北二',
				},
			];
			pagination.current = 1;
			pagination.total = 3;
		} catch (err) {
			// 错误处理
		} finally {
			setLoading(false);
		}
	};
	fetchData();
	// 改变页码
	const onPageChange = (current: number) => {
		fetchData({ ...basePagination, current });
	};
</script>

<script lang="ts">
	export default {
		name: 'Zone',
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
