<template>
  <div class="container">
    <Breadcrumb :items="['系统设置', '引擎设置']" />
  </div>
</template>

<script lang="ts">
  export default {
    name: 'Jenkins',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
