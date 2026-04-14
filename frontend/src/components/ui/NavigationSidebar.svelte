<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getActiveAccount, session } from '../../stores/session';
	import ProfilePopup from './ProfilePopup.svelte';
	import ZeroState from './ZeroState.svelte';

	export let currentPath: string = '/';

	let profileOpen = false;
	let account = getActiveAccount();
	let filtersOpen = true;
	let teamsOpen = true;

	const STORAGE_KEY = 'pulse_sidebar_state';

	onMount(() => {
		try {
			const raw = localStorage.getItem(STORAGE_KEY);
			if (raw) {
				const state = JSON.parse(raw);
				filtersOpen = state.filtersOpen ?? true;
				teamsOpen = state.teamsOpen ?? true;
			}
		} catch {
			/* ignore */
		}
	});

	session.subscribe(() => {
		account = getActiveAccount();
	});

	function saveSidebarState() {
		try {
			localStorage.setItem(STORAGE_KEY, JSON.stringify({ filtersOpen, teamsOpen }));
		} catch {
			/* ignore */
		}
	}

	function toggleFilters() {
		filtersOpen = !filtersOpen;
		saveSidebarState();
	}

	function toggleTeams() {
		teamsOpen = !teamsOpen;
		saveSidebarState();
	}
</script>

<aside class="sidebar">
	<div class="sidebar-header">
		<span class="logo-icon">◆</span>
		<span class="logo-text">Pulse</span>
	</div>

	<nav class="sidebar-nav">
		<button
			type="button"
			class="nav-item"
			class:active={currentPath === '/'}
			on:click={() => goto('/')}
		>
			Personal Activity
		</button>

		<div class="section">
			<button type="button" class="section-header" on:click={toggleFilters}>
				<span class="chevron" class:open={filtersOpen}>›</span>
				Личные фильтры
			</button>
			{#if filtersOpen}
				<div class="section-body">
					<button type="button" class="action-btn">+ Создать фильтр</button>
					<ZeroState message="Нет фильтров" />
					<button type="button" class="action-btn">Импортировать задачи</button>
				</div>
			{/if}
		</div>

		<div class="section">
			<button type="button" class="section-header" on:click={toggleTeams}>
				<span class="chevron" class:open={teamsOpen}>›</span>
				Команды
			</button>
			{#if teamsOpen}
				<div class="section-body">
					<button type="button" class="action-btn">+ Создать команду</button>
					<ZeroState message="Нет команд" />
				</div>
			{/if}
		</div>

		<button type="button" class="nav-item">Автоматизации</button>
		<button type="button" class="nav-item">Метки</button>
		<button
			type="button"
			class="nav-item"
			class:active={currentPath.startsWith('/about')}
			on:click={() => goto('/about')}
		>
			О системе
		</button>
		<button
			type="button"
			class="nav-item"
			class:active={currentPath.startsWith('/docs')}
			on:click={() => goto('/docs')}
		>
			Документация
		</button>
	</nav>

	<div class="sidebar-footer">
		{#if account}
			<button
				type="button"
				class="profile-btn"
				on:click={() => (profileOpen = !profileOpen)}
				aria-label="Профиль"
			>
				<span class="avatar">{account.username[0]?.toUpperCase()}</span>
				<span class="username">{account.username}</span>
			</button>
		{/if}
		<ProfilePopup bind:open={profileOpen} />
	</div>
</aside>

<style>
	.sidebar {
		width: 240px;
		background: #1e1e2e;
		color: #cdd6f4;
		display: flex;
		flex-direction: column;
		flex-shrink: 0;
		position: relative;
		height: 100vh;
	}

	.sidebar-header {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 20px 16px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.06);
		flex-shrink: 0;
	}

	.logo-icon {
		font-size: 14px;
		color: #cba6f7;
	}

	.logo-text {
		font-size: 18px;
		font-weight: 700;
		color: #cba6f7;
	}

	.sidebar-nav {
		flex: 1;
		overflow-y: auto;
		padding: 8px 0;
	}

	.nav-item {
		display: flex;
		align-items: center;
		width: 100%;
		padding: 8px 16px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 13px;
		color: #cdd6f4;
		text-align: left;
		transition: background 0.1s;
		border-radius: 0;
	}

	.nav-item:hover {
		background: rgba(255, 255, 255, 0.06);
	}

	.nav-item.active {
		background: rgba(99, 102, 241, 0.15);
		color: #a5b4fc;
	}

	.section {
		margin: 2px 0;
	}

	.section-header {
		display: flex;
		align-items: center;
		gap: 6px;
		width: 100%;
		padding: 7px 16px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 11px;
		font-weight: 600;
		color: rgba(205, 214, 244, 0.5);
		text-transform: uppercase;
		letter-spacing: 0.05em;
		text-align: left;
		transition: background 0.1s;
	}

	.section-header:hover {
		background: rgba(255, 255, 255, 0.04);
	}

	.chevron {
		font-size: 14px;
		line-height: 1;
		display: inline-block;
		transition: transform 0.15s;
		transform: rotate(0deg);
	}

	.chevron.open {
		transform: rotate(90deg);
	}

	.section-body {
		padding: 0 0 4px 12px;
	}

	.action-btn {
		display: block;
		width: 100%;
		padding: 6px 16px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 12px;
		color: rgba(205, 214, 244, 0.45);
		text-align: left;
		transition: color 0.1s;
	}

	.action-btn:hover {
		color: #cdd6f4;
	}

	.sidebar-footer {
		margin-top: auto;
		padding: 12px 8px;
		border-top: 1px solid rgba(255, 255, 255, 0.06);
		flex-shrink: 0;
	}

	.profile-btn {
		display: flex;
		align-items: center;
		gap: 10px;
		width: 100%;
		background: none;
		border: none;
		cursor: pointer;
		padding: 8px;
		border-radius: 6px;
		color: #cdd6f4;
		transition: background 0.1s;
	}

	.profile-btn:hover {
		background: rgba(255, 255, 255, 0.06);
	}

	.avatar {
		width: 28px;
		height: 28px;
		background: #6366f1;
		color: #fff;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 12px;
		font-weight: 600;
		flex-shrink: 0;
	}

	.username {
		font-size: 13px;
		font-weight: 500;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
