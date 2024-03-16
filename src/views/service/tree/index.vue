<template>
	<div class="container">
		<Breadcrumb :items="['项目服务', '服务中心']" />
		<div class="content">
			<div class="content-left"> <STree></STree></div>
			<div class="content-right">
				<a-row :gutter="20" align="stretch">
					<a-col :span="24">
						<a-card class="general-card" :title="'服务详情'">
							<a-row justify="space-between">
								<a-col :span="24">
									<a-tabs :default-active-tab="1" type="rounded">
										<a-tab-pane key="1" :title="'基本信息'"><Basic></Basic> </a-tab-pane>
										<a-tab-pane key="2" :title="'资产列表'"><Asset></Asset> </a-tab-pane>
										<a-tab-pane key="3" :title="'实时监控'"> </a-tab-pane>
										<a-tab-pane key="4" :title="'部署模板'"> </a-tab-pane>
										<a-tab-pane key="5" :title="'版本列表'"> </a-tab-pane>
										<a-tab-pane key="6" :title="'构建任务'"> </a-tab-pane>
										<a-tab-pane key="7" :title="'部署历史'"> </a-tab-pane>
									</a-tabs>
								</a-col>
							</a-row>
						</a-card>
					</a-col>
				</a-row>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import STree from './components/stree.vue';
	import Asset from './components/asset.vue';
	import Basic from './components/basic.vue';
</script>

<script lang="ts">
	export default {
		name: 'Tree',
	};
</script>

<style scoped lang="less">
	.container {
		padding: 0 20px 20px 20px;
	}

	.content {
		display: flex;
		margin-top: 12px;

		&-left {
			width: 292px;
			margin-right: 16px;
		}

		&-right {
			flex: 1;
			overflow: hidden;
		}

		.tab-pane-wrapper {
			padding: 0 16px 16px 16px;
		}
	}
</style>

<style lang="less" scoped>
	.mobile {
		.content {
			display: block;
			&-left {
				margin-right: 0;
				margin-bottom: 16px;
			}
			&-right {
				width: 100%;
			}
		}
	}
</style>
