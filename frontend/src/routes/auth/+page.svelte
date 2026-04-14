<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '../../api/auth';
	import { getActiveAccount } from '../../stores/session';

	let email = '';
	let emailError = '';
	let loading = false;
	let sent = false;
	let touched = false;

	onMount(() => {
		if (getActiveAccount()) {
			goto('/');
		}
	});

	function validateEmail(value: string): string {
		if (!value.trim()) return 'Введите email';
		if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return 'Введите корректный email';
		return '';
	}

	function handleBlur() {
		touched = true;
		emailError = validateEmail(email);
	}

	async function handleSubmit() {
		touched = true;
		emailError = validateEmail(email);
		if (emailError) return;

		loading = true;
		try {
			await auth.sendMagicLink(email);
			sent = true;
		} catch (error) {
			emailError = error instanceof Error ? error.message : 'Ошибка отправки';
		} finally {
			loading = false;
		}
	}

	async function handleResend() {
		sent = false;
		loading = true;
		try {
			await auth.sendMagicLink(email);
			sent = true;
		} catch (error) {
			emailError = error instanceof Error ? error.message : 'Ошибка отправки';
			sent = false;
		} finally {
			loading = false;
		}
	}
</script>

<div class="auth-screen">
	<div class="auth-card">
		<h1 class="auth-title">Pulse</h1>

		{#if sent}
			<div class="sent-state">
				<p class="sent-message">
					Проверьте почту — ссылка отправлена на <strong>{email}</strong>
				</p>
				<button type="button" class="btn-secondary" on:click={handleResend} disabled={loading}>
					{#if loading}
						<span class="loader" />
					{:else}
						Отправить снова
					{/if}
				</button>
			</div>
		{:else}
			<form on:submit|preventDefault={handleSubmit} novalidate>
				<div class="field" class:has-error={touched && emailError}>
					<label for="email">Email</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						on:blur={handleBlur}
						placeholder="you@example.com"
						autocomplete="email"
						disabled={loading}
					/>
					{#if touched && emailError}
						<span class="field-error">{emailError}</span>
					{/if}
				</div>

				<button type="submit" class="btn-primary" disabled={loading}>
					{#if loading}
						<span class="loader" />
					{:else}
						Войти
					{/if}
				</button>
			</form>
		{/if}
	</div>
</div>

<style>
	.auth-screen {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: #f5f5f5;
	}

	.auth-card {
		background: #fff;
		border-radius: 12px;
		padding: 40px;
		width: 100%;
		max-width: 400px;
		box-shadow: 0 2px 16px rgba(0, 0, 0, 0.08);
	}

	.auth-title {
		font-size: 28px;
		font-weight: 700;
		margin: 0 0 32px;
		text-align: center;
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 6px;
		margin-bottom: 20px;
	}

	label {
		font-size: 14px;
		font-weight: 500;
		color: #374151;
	}

	input {
		border: 1px solid #d1d5db;
		border-radius: 6px;
		padding: 10px 12px;
		font-size: 14px;
		outline: none;
		transition: border-color 0.15s;
	}

	input:focus {
		border-color: #6366f1;
	}

	.has-error input {
		border-color: #ef4444;
	}

	.field-error {
		font-size: 12px;
		color: #ef4444;
	}

	.btn-primary {
		width: 100%;
		padding: 11px;
		background: #6366f1;
		color: #fff;
		border: none;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 40px;
		transition: background 0.15s;
	}

	.btn-primary:hover:not(:disabled) {
		background: #4f46e5;
	}

	.btn-primary:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.btn-secondary {
		width: 100%;
		padding: 11px;
		background: transparent;
		color: #6366f1;
		border: 1px solid #6366f1;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 40px;
		transition: background 0.15s;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #eef2ff;
	}

	.sent-state {
		display: flex;
		flex-direction: column;
		gap: 20px;
		text-align: center;
	}

	.sent-message {
		font-size: 14px;
		color: #374151;
		line-height: 1.5;
		margin: 0;
	}

	.loader {
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.4);
		border-top-color: #fff;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
		display: inline-block;
	}

	.btn-secondary .loader {
		border-color: rgba(99, 102, 241, 0.3);
		border-top-color: #6366f1;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
