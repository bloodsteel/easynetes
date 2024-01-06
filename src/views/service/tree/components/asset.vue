<template>
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
</template>

<script lang="ts" setup>
  import { computed, ref, reactive } from 'vue';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import useLoading from '@/hooks/loading';
  import { Pagination } from '@/types/global';
  import { HostRecord, HostParams } from '@/types/cmdb';

  const { loading, setLoading } = useLoading(true);
  // 表格密度 框架自带的值
  type SizeProps = 'mini' | 'small' | 'medium' | 'large';
  const size = ref<SizeProps>('medium');
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
  // 获取数据  TODO: 这里没搞通
  const renderData = ref<HostRecord[]>([]);
  const fetchData = async (
    params: HostParams = { current: 1, pageSize: 20 },
  ) => {
    setLoading(true);
    try {
      renderData.value = [
        {
          id: 0,
          hostID: 'awefw',
          hostName: 'xxxx',
          hostType: '虚拟机',
          createdTime: '2024-1-1',
          status: '已上线',
        },
        {
          id: 1,
          hostID: 'awefw',
          hostName: 'xxxx',
          hostType: '虚拟机',
          createdTime: '2024-1-2',
          status: '已下线',
        },
        {
          id: 2,
          hostID: 'awefw',
          hostName: 'xxxx',
          hostType: '虚拟机',
          createdTime: '2024-1-3',
          status: '已上线',
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
    name: 'Asset',
  };
</script>
