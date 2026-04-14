<script lang="ts">
	let activeTab: 'guide' | 'api' = 'guide';
</script>

<div class="docs-page">
	<div class="tabs">
		<button class:active={activeTab === 'guide'} on:click={() => (activeTab = 'guide')}>
			Руководство
		</button>
		<button class:active={activeTab === 'api'} on:click={() => (activeTab = 'api')}> API </button>
	</div>

	{#if activeTab === 'guide'}
		<div class="guide-content">
			<h1>Что такое Pulse</h1>
			<p>
				Pulse — трекер задач для команд и личного использования. Позволяет организовывать работу
				через гибкие фильтры, управлять командами и отслеживать активность.
			</p>

			<h2>Ключевые сущности</h2>

			<h3>Задачи</h3>
			<p>
				Задача — основная единица работы. Каждая задача имеет название, описание, статус и может
				быть привязана к команде или оставаться личной.
			</p>

			<h3>Фильтры</h3>
			<p>
				Фильтр — сохранённый набор условий для отображения задач. Фильтры бывают личными (видны
				только вам) и командными (доступны участникам команды). Кликнув по фильтру, вы переходите
				на экран результатов с отфильтрованным списком задач.
			</p>

			<h3>Команды</h3>
			<p>
				Команда объединяет нескольких пользователей и их общие задачи. Участники команды имеют
				разные уровни доступа, определяемые через разрешения.
			</p>

			<h2>Авторизация через magic link</h2>
			<p>
				Pulse не использует пароли. При входе на указанный email отправляется одноразовая ссылка.
				По клику на ссылку выполняется автоматический вход. Ссылка действительна 15 минут.
			</p>

			<h2>Основные принципы работы</h2>
			<ul>
				<li>Навигация осуществляется через боковую панель слева.</li>
				<li>Личные фильтры отображаются в разделе «Личные фильтры» сайдбара.</li>
				<li>Командные фильтры отображаются под соответствующей командой.</li>
				<li>Ошибки всегда сопровождаются уведомлением в правом нижнем углу экрана.</li>
				<li>Опасные операции (удаление и т.п.) требуют подтверждения во всплывающем окне.</li>
			</ul>
		</div>
	{:else}
		<div class="api-content">
			<iframe src="/api/swagger/index.html" title="API Documentation" class="swagger-frame" />
		</div>
	{/if}
</div>

<style>
	.docs-page {
		padding: 32px;
		max-width: 900px;
	}

	.tabs {
		display: flex;
		gap: 4px;
		margin-bottom: 24px;
		border-bottom: 1px solid #e5e7eb;
	}

	.tabs button {
		padding: 8px 16px;
		background: none;
		border: none;
		border-bottom: 2px solid transparent;
		margin-bottom: -1px;
		cursor: pointer;
		font-size: 14px;
		color: #6b7280;
		transition: color 0.15s;
	}

	.tabs button.active {
		color: #6366f1;
		border-bottom-color: #6366f1;
		font-weight: 500;
	}

	.tabs button:hover:not(.active) {
		color: #374151;
	}

	.guide-content {
		font-size: 14px;
		line-height: 1.7;
		color: #374151;
	}

	.guide-content h1 {
		font-size: 22px;
		color: #111827;
		margin: 0 0 16px;
	}

	.guide-content h2 {
		font-size: 17px;
		color: #111827;
		margin: 28px 0 12px;
	}

	.guide-content h3 {
		font-size: 14px;
		font-weight: 600;
		color: #374151;
		margin: 20px 0 8px;
	}

	.guide-content p {
		margin: 0 0 12px;
	}

	.guide-content ul {
		margin: 0;
		padding-left: 20px;
	}

	.guide-content li {
		margin-bottom: 6px;
	}

	.api-content {
		height: calc(100vh - 140px);
	}

	.swagger-frame {
		width: 100%;
		height: 100%;
		border: none;
		border-radius: 8px;
	}
</style>
