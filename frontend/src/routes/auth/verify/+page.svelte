<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '../../../api/auth';
	import { addAccount } from '../../../stores/session';

	let error = '';
	let verifying = true;

	onMount(async () => {
		const token = $page.url.searchParams.get('token');

		if (!token) {
			goto('/auth');
			return;
		}

		try {
			const result = await auth.verifyMagicLink(token);
			addAccount(result.user, result.access_token, result.refresh_token);
			goto('/');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Ссылка недействительна или уже использована';
			verifying = false;
		}
	});
</script>

<div class="verify-screen">
	{#if verifying && !error}
		<div class="status">
			<span class="spinner" />
			<p>Выполняется вход…</p>
		</div>
	{:else}
		<div class="error-state">
			<p class="error-message">{error}</p>
			<a href="/auth" class="btn-link">Запросить новую ссылку</a>
		</div>
	{/if}
</div>

<style>
	.verify-screen {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: #f5f5f5;
	}

	.status {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		color: #374151;
		font-size: 14px;
	}

	.status p {
		margin: 0;
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 3px solid #e5e7eb;
		border-top-color: #6366f1;
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
		display: inline-block;
	}

	.error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 20px;
		background: #fff;
		border-radius: 12px;
		padding: 40px;
		max-width: 400px;
		text-align: center;
		box-shadow: 0 2px 16px rgba(0, 0, 0, 0.08);
	}

	.error-message {
		font-size: 14px;
		color: #ef4444;
		margin: 0;
	}

	.btn-link {
		color: #6366f1;
		text-decoration: none;
		font-size: 14px;
		font-weight: 500;
	}

	.btn-link:hover {
		text-decoration: underline;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
