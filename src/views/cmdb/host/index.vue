<template>
	<div class="container">
		<Breadcrumb :items="['资产数据', '主机']" />
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
						<div class="action-icon" @click="refresh"><icon-refresh size="18" /></div>
					</a-tooltip>
					<a-dropdown @select="handleSelectDensity">
						<a-tooltip :content="'密度'">
							<div class="action-icon"><icon-line-height size="18" /></div>
						</a-tooltip>
						<template #content>
							<a-doption v-for="item in densityList" :key="item.value" :value="item.value"
								:class="{ active: item.value === size }">
								<span>{{ item.name }}</span>
							</a-doption>
						</template>
					</a-dropdown>
				</a-col>
			</a-row>
			<!-- 表格 -->
			<a-table row-key="id" :loading="loading" :pagination="false" :columns="columns" :data="renderData" :bordered="false"
				:size="size" @page-change="onPageChange">
				<!-- 索引 slot -->
				<template #index="{ rowIndex }">
					{{ rowIndex + 1 + (pagination.current - 1) * pagination.pageSize }}
				</template>
				<!-- 状态 slot -->
				<template #status="{ record }">
					<span v-if="record.status === 'offline'" class="circle"></span>
					<span v-else class="circle pass"></span>
					{{ '已上线' }}
				</template>
				<!-- 操作 slot -->
				<template #operations>
					<a-button type="text" size="small">一般</a-button>
					<a-button type="text" size="small" status="success">成功</a-button>
					<a-button type="text" size="small" status="warning">警告</a-button>
					<a-button type="text" size="small" status="danger">危险</a-button>
				</template>
			</a-table>
			<!-- 分页组件 -->
			<a-pagination :total="pagination.total ? pagination.total : 0" v-model:current="pagination.current"
				v-model:page-size="pagination.pageSize" show-jumper show-total show-page-size @change="onPageChange"
				@page-size-change="onPageSizeChange">

			</a-pagination>
		</a-card>
	</div>
</template>

<script lang="ts" setup>
import { computed, ref, reactive, onMounted } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import useLoading from '@/hooks/loading';
import { Pagination } from '@/types/global';
import { HostRecord, HostParams } from '@/types/cmdb';
import { queryCmdbData } from '@/api/cmdb';



// 数据
const { loading, setLoading } = useLoading(true);
// 表格密度 框架自带的值
type SizeProps = 'mini' | 'small' | 'medium' | 'large';
const size = ref<SizeProps>('medium');
const handleSelectDensity = (
	val: string | number | Record<string, any> | undefined,
	e: Event,
) => {
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
		dataIndex: 'createdTime',
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
	queryCmdbData({ current: pagination.current, pageSize: pagination.pageSize }).then((res) => {
		console.log("the res is ", res);
		pagination.total = res.data.total;
		renderData.value = res.data.data;
	}).catch((err) => {
		console.log("get api error", err);
	}).finally(() => {
		setLoading(false);
	})
};
// 初始化时加载数据
onMounted(() => {
	fetchData();
	console.log("the renderData is ", renderData)
})
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
		pagination.total = 0
	}
	console.log("pageSize is 1: ", pagination);
	if (pagination.current > 1) {
		if (((pagination.current - 1) * pagination.pageSize) > pagination.total) {
			pagination.current = 1
		}
	}
	fetchData();
	console.log("pageSize is 2: ", pagination);
	// return true;
}
// 刷新
const refresh = () => {

	fetchData();
}
// 搜索
const search = () => {
	console.log('searching...');
}
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
