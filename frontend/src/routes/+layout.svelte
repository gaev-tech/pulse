<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getActiveAccount } from '../stores/session';
	import ErrorToast from '../components/ui/ErrorToast.svelte';

	const PUBLIC_ROUTES = ['/auth', '/auth/verify'];

	onMount(() => {
		const isPublic = PUBLIC_ROUTES.some((route) => $page.url.pathname.startsWith(route));
		if (!isPublic && !getActiveAccount()) {
			goto('/auth');
		}
	});
</script>

<ErrorToast />
<slot />
