<template>
  <a-card :bordered="false">
    <a-input-search v-model="searchKey" style="margin-bottom: 8px" />
    <a-tree
      block-node
      :default-expand-all="false"
      :size="size"
      :data="treeData"
      :show-line="showLine"
      :load-more="loadMore"
    >
    </a-tree>
  </a-card>
</template>

<script lang="ts" setup>
  import { ref, computed } from 'vue';
  import { TreeNodeData } from '@arco-design/web-vue/es/tree/interface';

  type SizeProps = 'mini' | 'small' | 'medium' | 'large';

  const showLine = ref(true);
  const size = ref<SizeProps>('large');
  const searchKey = ref('');
  // 原始数据
  const originTreeData = ref([
    {
      title: 'Trunk 0-0',
      key: '0-0',
    },
    {
      title: 'Trunk 0-1',
      key: '0-1',
      children: [
        {
          title: 'Branch 0-1-1',
          key: '0-1-1',
        },
      ],
    },
  ]);

  // 最终数据
  const treeData = computed(() => {
    if (!searchKey.value) return originTreeData as unknown as TreeNodeData[];
    return searchData(searchKey.value);
  });

  // 搜索数据
  function searchData(keyword: string) {
    const loop = (data: TreeNodeData[]) => {
      const result: TreeNodeData[] = [];
      data.forEach((item: TreeNodeData) => {
        if (
          item.title &&
          item.title.toLowerCase().indexOf(keyword.toLowerCase()) > -1
        ) {
          result.push({ ...item });
        } else if (item.children) {
          const filterData = loop(item.children);
          if (filterData.length) {
            result.push({
              ...item,
              children: filterData,
            });
          }
        }
      });
      return result;
    };
    return loop(originTreeData as unknown as TreeNodeData[]);
  }

  // 动态加载
  const loadMore = (nodeData: TreeNodeData) => {
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        nodeData.children = [
          { title: `leaf`, key: `${nodeData.key}-1`, isLeaf: true },
        ];
        resolve();
      }, 1000);
    });
  };
</script>

<script lang="ts">
  export default {
    name: 'Tree',
  };
</script>
