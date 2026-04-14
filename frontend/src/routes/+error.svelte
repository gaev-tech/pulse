<script lang="ts">
	import { page } from '$app/stores';
</script>

<div class="error-page">
	{#if $page.status === 403}
		<h1>Нет доступа</h1>
		<p>У вас нет прав для просмотра этой страницы.</p>
	{:else if $page.status === 404}
		<h1>Не найдено</h1>
		<p>Запрашиваемая страница не существует или была удалена.</p>
	{:else}
		<h1>Произошла ошибка</h1>
		<p>{$page.error?.message ?? 'Что-то пошло не так. Попробуйте обновить страницу.'}</p>
	{/if}
	<p class="status-code">Код ошибки: {$page.status}</p>
	<a href="/" class="home-link">На главную</a>
</div>

<style>
	.error-page {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: #f9fafb;
		text-align: center;
		padding: 32px;
	}

	h1 {
		font-size: 28px;
		color: #111827;
		margin: 0 0 12px;
	}

	p {
		font-size: 14px;
		color: #6b7280;
		margin: 0 0 8px;
	}

	.status-code {
		font-size: 12px;
		color: #9ca3af;
	}

	.home-link {
		margin-top: 24px;
		color: #6366f1;
		font-size: 14px;
		text-decoration: none;
		font-weight: 500;
	}

	.home-link:hover {
		text-decoration: underline;
	}
</style>
