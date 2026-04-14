<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { auth } from '../../api/auth';

	let email = '';
	let emailError = '';
	let loading = false;
	let sent = false;
	let touched = false;

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
				<Button variant="outline" class="w-full" on:click={handleResend} disabled={loading}>
					{#if loading}
						<span class="loader" />
					{:else}
						Отправить снова
					{/if}
				</Button>
			</div>
		{:else}
			<form on:submit|preventDefault={handleSubmit} novalidate>
				<div class="field">
					<label for="email">Email</label>
					<Input
						id="email"
						type="email"
						bind:value={email}
						on:blur={handleBlur}
						placeholder="you@example.com"
						autocomplete="email"
						disabled={loading}
						class={touched && emailError ? 'border-red-500' : ''}
					/>
					{#if touched && emailError}
						<span class="field-error">{emailError}</span>
					{/if}
				</div>

				<Button type="submit" class="w-full" disabled={loading}>
					{#if loading}
						<span class="loader" />
					{:else}
						Войти
					{/if}
				</Button>
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
		background: hsl(var(--background));
	}

	.auth-card {
		background: hsl(var(--card));
		color: hsl(var(--card-foreground));
		border: 1px solid hsl(var(--border));
		border-radius: 12px;
		padding: 40px;
		width: 100%;
		max-width: 400px;
		box-shadow: 0 2px 16px rgba(0, 0, 0, 0.08);
	}

	.auth-title {
		font-size: 28px;
		font-weight: 700;
		color: hsl(var(--foreground));
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
		color: hsl(var(--foreground));
	}

	.field-error {
		font-size: 12px;
		color: hsl(var(--destructive));
	}

	.sent-state {
		display: flex;
		flex-direction: column;
		gap: 20px;
		text-align: center;
	}

	.sent-message {
		font-size: 14px;
		color: hsl(var(--foreground));
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

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
