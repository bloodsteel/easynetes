<template>
  <div class="container">
    <Breadcrumb :items="['容器集群', '命名空间']" />
  </div>
</template>

<script lang="ts">
  export default {
    name: 'Namespace',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
