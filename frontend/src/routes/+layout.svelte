<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getActiveAccount } from '../stores/session';
	import { initTheme } from '../stores/theme';
	import ErrorToast from '../components/ui/ErrorToast.svelte';
	import NavigationSidebar from '../components/ui/NavigationSidebar.svelte';

	const PUBLIC_ROUTES = ['/auth', '/auth/verify'];

	$: isPublic = PUBLIC_ROUTES.some((route) => $page.url.pathname.startsWith(route));

	onMount(() => {
		initTheme();
		if (!isPublic && !getActiveAccount()) {
			goto('/auth');
		}
	});
</script>

<ErrorToast />

{#if isPublic}
	<slot />
{:else}
	<div class="app-shell">
		<NavigationSidebar currentPath={$page.url.pathname} />
		<main class="main-content">
			<slot />
		</main>
	</div>
{/if}

<style>
	.app-shell {
		display: flex;
		height: 100vh;
		overflow: hidden;
	}

	.main-content {
		flex: 1;
		overflow: auto;
		background: hsl(var(--background));
	}
</style>
