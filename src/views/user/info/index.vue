<template>
	<div class="container">
		<!-- 面包屑 https://arco.design/vue/component/breadcrumb 注意：这里是全局组件 -->
		<Breadcrumb :items="['个人中心', '用户信息']" />
		<UserInfoHeader />
		<div class="content">
			<!-- 我的项目 -->
			<div class="content-left">
				<!-- 栅格 https://arco.design/vue/component/grid -->
				<!-- cols 表示每一行展示的列数  col-gap表示列与列之间的间距 row-gap 表示行于行的间距 -->
				<a-grid :cols="24" :col-gap="16" :row-gap="16">
					<!-- span 跨越的格数 -->
					<a-grid-item :span="24">
						<!-- 我的项目 -->
						<MyProject />
					</a-grid-item>
					<a-grid-item :span="24">
						<!-- 最新动态 -->
						<LatestActivity />
					</a-grid-item>
				</a-grid>
			</div>
			<!-- 我的团队 -->
			<div class="content-right">
				<a-grid :cols="24" :row-gap="16">
					<a-grid-item :span="24">
						<!-- 我的团队 -->
						<MyTeam />
					</a-grid-item>
					<a-grid-item class="panel" :span="24">
						<!-- 站内通知 -->
						<LatestNotification />
					</a-grid-item>
				</a-grid>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import UserInfoHeader from './components/user-info-header.vue';
	import LatestNotification from './components/latest-notification.vue';
	import MyProject from './components/my-project.vue';
	import LatestActivity from './components/latest-activity.vue';
	import MyTeam from './components/my-team.vue';
</script>

<script lang="ts">
	export default {
		name: 'Info',
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
			flex: 1;
			margin-right: 16px;
			overflow: hidden;
			// background-color: var(--color-bg-2);

			:deep(.arco-tabs-nav-tab) {
				margin-left: 16px;
			}
		}

		&-right {
			width: 332px;
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
