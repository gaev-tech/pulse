<script lang="ts">
	import { goto } from '$app/navigation';
	import { session, getActiveAccount, removeAccount, switchAccount } from '../../stores/session';
	import { auth } from '../../api/auth';
	import { showError } from '../../stores/toast';

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

		<div class="actions">
			<button type="button" class="action-btn" on:click={handleAddUser}>
				+ Добавить пользователя
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
		background: #fff;
		border-radius: 10px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
		min-width: 240px;
		overflow: hidden;
	}

	.accounts-list {
		list-style: none;
		margin: 0;
		padding: 8px 0;
		border-bottom: 1px solid #f3f4f6;
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
		transition: background 0.1s;
	}

	.account-item:hover {
		background: #f9fafb;
	}

	.account-item.active {
		background: #eef2ff;
	}

	.avatar {
		width: 32px;
		height: 32px;
		background: #6366f1;
		color: #fff;
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
		color: #111827;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.email {
		font-size: 12px;
		color: #6b7280;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.actions {
		padding: 8px 0;
		display: flex;
		flex-direction: column;
	}

	.action-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		padding: 10px 16px;
		background: none;
		border: none;
		cursor: pointer;
		font-size: 13px;
		color: #374151;
		text-align: left;
		justify-content: flex-start;
		transition: background 0.1s;
	}

	.action-btn:hover:not(:disabled) {
		background: #f9fafb;
	}

	.logout-btn {
		color: #ef4444;
	}

	.logout-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.loader {
		width: 14px;
		height: 14px;
		border: 2px solid rgba(239, 68, 68, 0.3);
		border-top-color: #ef4444;
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
