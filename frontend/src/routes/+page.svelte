<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getActiveAccount, session } from '../stores/session';
	import ProfilePopup from '../components/ui/ProfilePopup.svelte';

	let profileOpen = false;
	let account = getActiveAccount();

	onMount(() => {
		if (!getActiveAccount()) {
			goto('/auth');
		}
	});

	session.subscribe(() => {
		account = getActiveAccount();
	});
</script>

{#if account}
	<div class="app">
		<aside class="sidebar">
			<div class="sidebar-top">
				<span class="logo">Pulse</span>
			</div>

			<div class="sidebar-bottom">
				<button
					type="button"
					class="profile-btn"
					on:click={() => (profileOpen = !profileOpen)}
					aria-label="Профиль"
				>
					<span class="avatar">{account.username[0]?.toUpperCase()}</span>
					<span class="username">{account.username}</span>
				</button>
			</div>

			<ProfilePopup bind:open={profileOpen} />
		</aside>

		<main class="main">
			<p class="placeholder">Personal Activity Screen — coming soon</p>
		</main>
	</div>
{/if}

<style>
	.app {
		display: flex;
		height: 100vh;
		overflow: hidden;
	}

	.sidebar {
		width: 240px;
		background: #1e1e2e;
		color: #cdd6f4;
		display: flex;
		flex-direction: column;
		flex-shrink: 0;
		position: relative;
	}

	.sidebar-top {
		padding: 20px 16px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.06);
	}

	.logo {
		font-size: 18px;
		font-weight: 700;
		color: #cba6f7;
	}

	.sidebar-bottom {
		margin-top: auto;
		padding: 12px 8px;
		border-top: 1px solid rgba(255, 255, 255, 0.06);
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

	.main {
		flex: 1;
		overflow: auto;
		padding: 32px;
		background: #f9fafb;
	}

	.placeholder {
		color: #9ca3af;
		font-size: 14px;
	}
</style>
