<template>
  <div class="container">
    <Breadcrumb :items="['容器集群', '工作负载']" />
  </div>
</template>

<script lang="ts">
  export default {
    name: 'Workload',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
