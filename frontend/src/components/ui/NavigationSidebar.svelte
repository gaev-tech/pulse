<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getActiveAccount, session } from '../../stores/session';
	import { Search, ChevronRight, ChevronDown, Plus, Activity } from 'lucide-svelte';
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
		<div class="logo-icon">
			<Activity size={16} />
		</div>
		<span class="logo-text">Pulse</span>
	</div>

	<div class="search-row">
		<button type="button" class="search-btn" aria-label="Поиск">
			<Search size={13} />
			<span>Поиск...</span>
		</button>
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
				{#if filtersOpen}
					<ChevronDown size={12} />
				{:else}
					<ChevronRight size={12} />
				{/if}
				Личные фильтры
			</button>
			{#if filtersOpen}
				<div class="section-body">
					<button type="button" class="action-btn">
						<Plus size={11} />
						Создать фильтр
					</button>
					<ZeroState message="Нет фильтров" />
					<button type="button" class="action-btn">Импортировать задачи</button>
				</div>
			{/if}
		</div>

		<div class="section">
			<button type="button" class="section-header" on:click={toggleTeams}>
				{#if teamsOpen}
					<ChevronDown size={12} />
				{:else}
					<ChevronRight size={12} />
				{/if}
				Команды
			</button>
			{#if teamsOpen}
				<div class="section-body">
					<button type="button" class="action-btn">
						<Plus size={11} />
						Создать команду
					</button>
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
		background: hsl(var(--card));
		color: hsl(var(--foreground));
		display: flex;
		flex-direction: column;
		flex-shrink: 0;
		position: relative;
		height: 100vh;
		border-right: 1px solid hsl(var(--border));
	}

	.sidebar-header {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 16px;
		border-bottom: 1px solid hsl(var(--border));
		flex-shrink: 0;
	}

	.logo-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		background: hsl(var(--primary));
		color: hsl(var(--primary-foreground));
		border-radius: 6px;
		flex-shrink: 0;
	}

	.logo-text {
		font-size: 15px;
		font-weight: 700;
		color: hsl(var(--foreground));
	}

	.search-row {
		padding: 8px 10px;
		flex-shrink: 0;
	}

	.search-btn {
		display: flex;
		align-items: center;
		gap: 6px;
		width: 100%;
		padding: 6px 8px;
		background: hsl(var(--muted));
		border: 1px solid hsl(var(--border));
		border-radius: 6px;
		cursor: pointer;
		font-size: 12px;
		color: hsl(var(--muted-foreground));
		transition: background 0.15s;
		text-align: left;
		font-family: inherit;
	}

	.search-btn:hover {
		background: hsl(var(--secondary));
	}

	.sidebar-nav {
		flex: 1;
		overflow-y: auto;
		padding: 4px 0;
	}

	.nav-item {
		display: flex;
		align-items: center;
		width: 100%;
		padding: 7px 12px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 13px;
		color: hsl(var(--foreground));
		text-align: left;
		transition: background 0.15s;
		font-family: inherit;
	}

	.nav-item:hover {
		background: hsl(var(--muted));
	}

	.nav-item.active {
		background: hsl(var(--accent));
		color: hsl(var(--accent-foreground));
		font-weight: 500;
	}

	.section {
		margin: 2px 0;
	}

	.section-header {
		display: flex;
		align-items: center;
		gap: 4px;
		width: 100%;
		padding: 6px 12px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 11px;
		font-weight: 600;
		color: hsl(var(--muted-foreground));
		text-transform: uppercase;
		letter-spacing: 0.05em;
		text-align: left;
		transition: background 0.15s;
		font-family: inherit;
	}

	.section-header:hover {
		background: hsl(var(--muted));
	}

	.section-body {
		padding: 0 0 4px 16px;
	}

	.action-btn {
		display: flex;
		align-items: center;
		gap: 5px;
		width: 100%;
		padding: 5px 12px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 12px;
		color: hsl(var(--muted-foreground));
		text-align: left;
		transition: color 0.15s;
		font-family: inherit;
	}

	.action-btn:hover {
		color: hsl(var(--foreground));
	}

	.sidebar-footer {
		margin-top: auto;
		padding: 10px 8px;
		border-top: 1px solid hsl(var(--border));
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
		color: hsl(var(--foreground));
		transition: background 0.15s;
		text-align: left;
		font-family: inherit;
	}

	.profile-btn:hover {
		background: hsl(var(--muted));
	}

	.avatar {
		width: 28px;
		height: 28px;
		background: hsl(var(--primary));
		color: hsl(var(--primary-foreground));
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
