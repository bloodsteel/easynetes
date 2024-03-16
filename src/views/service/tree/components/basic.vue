<template>
	<a-space direction="vertical" fill>
		<a-card class="general-card">
			<ProfileItem :loading="loading" :render-data="currentData" />
		</a-card>
		<a-divider style="margin-bottom: 0" />
	</a-space>
</template>

<script lang="ts" setup>
	import { ref } from 'vue';
	import useLoading from '@/hooks/loading';
	import { queryProfileBasic, ProfileBasicRes } from '@/api/service';
	import ProfileItem from './test.vue';

	const { loading, setLoading } = useLoading(true);
	// 基本信息
	const currentData = ref<ProfileBasicRes>({} as ProfileBasicRes);
	const fetchCurrentData = async () => {
		try {
			const { data } = await queryProfileBasic();
			currentData.value = data;
		} catch (err) {
			// you can report use errorHandler or other
		} finally {
			setLoading(false);
		}
	};
	fetchCurrentData();
</script>

<script lang="ts">
	export default {
		name: 'Basic',
	};
</script>
