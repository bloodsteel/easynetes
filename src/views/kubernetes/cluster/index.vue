<template>
  <div class="container">
    <Breadcrumb :items="['容器集群', '集群管理']" />
  </div>
</template>

<script lang="ts">
  export default {
    name: 'Cluster',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
