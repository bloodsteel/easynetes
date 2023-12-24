<template>
  <a-card class="general-card" :title="'站内通知'">
    <!-- 将加载中的数据用灰色占位。 https://arco.design/vue/component/skeleton -->
    <a-skeleton v-if="loading" :animation="true">
      <a-skeleton-line :rows="3" />
    </a-skeleton>
    <a-result v-else status="404">
      <template #subtitle>
        {{ '暂无数据' }}
      </template>
    </a-result>
  </a-card>
</template>

<script lang="ts" setup>
  import useLoading from '@/hooks/loading';

  // 这里表示是否在加载中
  const { loading, setLoading } = useLoading(true);
  setTimeout(() => {
    setLoading(false);
  }, 500);
</script>

<style lang="less" scoped>
  :deep(.arco-result) {
    padding: 40px 32px 108px;
  }
</style>
