<script lang="ts">
	import { Button } from '$lib/components/ui/button';

	export let open: boolean = false;
	export let onConfirm: () => void;
	export let onCancel: () => void;
</script>

{#if open}
	<div class="backdrop" on:click={onCancel} role="presentation" />
	<div class="popup" role="dialog" aria-modal="true">
		<p class="message">Вы уверены? Несохранённые данные будут потеряны.</p>
		<div class="actions">
			<Button variant="outline" on:click={onCancel}>Отмена</Button>
			<Button variant="destructive" on:click={onConfirm}>Да, закрыть</Button>
		</div>
	</div>
{/if}

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.4);
		z-index: 99;
	}

	.popup {
		position: fixed;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		z-index: 100;
		background: hsl(var(--card));
		color: hsl(var(--card-foreground));
		border: 1px solid hsl(var(--border));
		border-radius: 10px;
		padding: 24px;
		min-width: 320px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.16);
	}

	.message {
		font-size: 14px;
		color: hsl(var(--foreground));
		margin: 0 0 20px;
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
	}
</style>
