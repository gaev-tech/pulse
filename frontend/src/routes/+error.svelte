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
		background: hsl(var(--background));
		text-align: center;
		padding: 32px;
	}

	h1 {
		font-size: 28px;
		font-weight: 700;
		color: hsl(var(--foreground));
		margin: 0 0 12px;
	}

	p {
		font-size: 14px;
		color: hsl(var(--muted-foreground));
		margin: 0 0 8px;
	}

	.status-code {
		font-size: 12px;
		opacity: 0.6;
	}

	.home-link {
		margin-top: 24px;
		color: hsl(var(--primary));
		font-size: 14px;
		text-decoration: none;
		font-weight: 500;
	}

	.home-link:hover {
		text-decoration: underline;
	}
</style>
