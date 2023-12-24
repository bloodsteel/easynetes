<template>
  <div class="container">
    <Breadcrumb :items="['项目服务', '服务树']" />
  </div>
</template>

<script lang="ts">
  export default {
    name: 'Tree',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
