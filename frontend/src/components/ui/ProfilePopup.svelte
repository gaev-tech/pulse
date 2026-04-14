<script lang="ts">
	import { goto } from '$app/navigation';
	import { Sun, Moon, Monitor, LogOut, UserPlus } from 'lucide-svelte';
	import { session, getActiveAccount, removeAccount, switchAccount } from '../../stores/session';
	import { auth } from '../../api/auth';
	import { theme, setTheme } from '../../stores/theme';
	import type { ThemeMode } from '../../stores/theme';

	export let open = false;

	let loggingOut = false;

	async function handleLogout() {
		const account = getActiveAccount();
		if (!account) return;

		loggingOut = true;
		try {
			await auth.logout(account.refreshToken);
		} catch {
			// продолжаем logout даже при сетевой ошибке
		} finally {
			removeAccount(account.id);
			loggingOut = false;
			open = false;
			const remaining = getActiveAccount();
			goto(remaining ? '/' : '/auth');
		}
	}

	function handleSwitch(id: string) {
		switchAccount(id);
		open = false;
		goto('/');
	}

	function handleAddUser() {
		open = false;
		goto('/auth');
	}

	function handleBackdropClick() {
		open = false;
	}

	const THEME_OPTIONS: Array<{ mode: ThemeMode; icon: typeof Sun; label: string }> = [
		{ mode: 'light', icon: Sun, label: 'Светлая' },
		{ mode: 'dark', icon: Moon, label: 'Тёмная' },
		{ mode: 'system', icon: Monitor, label: 'Системная' }
	];
</script>

{#if open}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div class="backdrop" on:click={handleBackdropClick} />

	<div class="popup" role="dialog" aria-label="Профиль">
		<ul class="accounts-list">
			{#each $session.accounts as account, index (account.id)}
				<li>
					<button
						type="button"
						class="account-item"
						class:active={index === $session.activeIndex}
						on:click={() => handleSwitch(account.id)}
					>
						<span class="avatar">{account.username[0]?.toUpperCase()}</span>
						<span class="account-info">
							<span class="username">{account.username}</span>
							<span class="email">{account.email}</span>
						</span>
					</button>
				</li>
			{/each}
		</ul>

		<div class="theme-section">
			<span class="theme-label">Тема</span>
			<div class="theme-buttons">
				{#each THEME_OPTIONS as opt}
					<button
						type="button"
						class="theme-btn"
						class:active={$theme === opt.mode}
						on:click={() => setTheme(opt.mode)}
						aria-label={opt.label}
						title={opt.label}
					>
						<svelte:component this={opt.icon} size={14} />
					</button>
				{/each}
			</div>
		</div>

		<div class="actions">
			<button type="button" class="action-btn" on:click={handleAddUser}>
				<UserPlus size={14} />
				Добавить пользователя
			</button>
			<button
				type="button"
				class="action-btn logout-btn"
				on:click={handleLogout}
				disabled={loggingOut}
			>
				{#if loggingOut}
					<span class="loader" />
				{:else}
					<LogOut size={14} />
					Выйти
				{/if}
			</button>
		</div>
	</div>
{/if}

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		z-index: 99;
	}

	.popup {
		position: fixed;
		bottom: 64px;
		left: 16px;
		z-index: 100;
		background: hsl(var(--popover));
		color: hsl(var(--popover-foreground));
		border: 1px solid hsl(var(--border));
		border-radius: 10px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
		min-width: 240px;
		overflow: hidden;
	}

	.accounts-list {
		list-style: none;
		margin: 0;
		padding: 8px 0;
		border-bottom: 1px solid hsl(var(--border));
	}

	.account-item {
		display: flex;
		align-items: center;
		gap: 10px;
		width: 100%;
		padding: 10px 16px;
		background: none;
		border: none;
		cursor: pointer;
		text-align: left;
		transition: background 0.15s;
		font-family: inherit;
	}

	.account-item:hover {
		background: hsl(var(--muted));
	}

	.account-item.active {
		background: hsl(var(--accent));
	}

	.avatar {
		width: 32px;
		height: 32px;
		background: hsl(var(--primary));
		color: hsl(var(--primary-foreground));
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 13px;
		font-weight: 600;
		flex-shrink: 0;
	}

	.account-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
		overflow: hidden;
	}

	.username {
		font-size: 13px;
		font-weight: 500;
		color: hsl(var(--foreground));
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.email {
		font-size: 12px;
		color: hsl(var(--muted-foreground));
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.theme-section {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 10px 16px;
		border-bottom: 1px solid hsl(var(--border));
	}

	.theme-label {
		font-size: 12px;
		color: hsl(var(--muted-foreground));
	}

	.theme-buttons {
		display: flex;
		gap: 2px;
		background: hsl(var(--muted));
		border-radius: 6px;
		padding: 2px;
	}

	.theme-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		background: none;
		border: none;
		cursor: pointer;
		border-radius: 4px;
		color: hsl(var(--muted-foreground));
		transition: background 0.15s, color 0.15s;
	}

	.theme-btn:hover {
		background: hsl(var(--secondary));
		color: hsl(var(--foreground));
	}

	.theme-btn.active {
		background: hsl(var(--card));
		color: hsl(var(--foreground));
	}

	.actions {
		padding: 8px 0;
		display: flex;
		flex-direction: column;
	}

	.action-btn {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 10px 16px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 13px;
		color: hsl(var(--foreground));
		text-align: left;
		transition: background 0.15s;
		font-family: inherit;
	}

	.action-btn:hover:not(:disabled) {
		background: hsl(var(--muted));
	}

	.logout-btn {
		color: hsl(var(--destructive));
	}

	.logout-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.loader {
		width: 14px;
		height: 14px;
		border: 2px solid hsl(var(--destructive) / 0.3);
		border-top-color: hsl(var(--destructive));
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
		display: inline-block;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
