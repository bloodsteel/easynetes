<template>
	<div class="container">
		<Breadcrumb :items="['系统设置', '仓库设置']" />
	</div>
</template>

<script lang="ts">
	export default {
		name: 'Gitlab',
	};
</script>

<style scoped lang="less">
	.container {
		padding: 0 20px 20px 20px;
	}
</style>
