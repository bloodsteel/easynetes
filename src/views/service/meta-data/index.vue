<template>
	<div class="container">
		<Breadcrumb :items="['项目服务', '元数据']" />
	</div>
</template>

<script lang="ts">
	export default {
		name: 'MetaData',
	};
</script>

<style scoped lang="less">
	.container {
		padding: 0 20px 20px 20px;
	}
</style>
