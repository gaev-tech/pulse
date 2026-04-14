<script lang="ts">
	export let open: boolean = false;
	export let onConfirm: () => void;
	export let onCancel: () => void;
</script>

{#if open}
	<div class="backdrop" on:click={onCancel} role="presentation" />
	<div class="popup" role="dialog" aria-modal="true">
		<p class="message">Вы уверены? Несохранённые данные будут потеряны.</p>
		<div class="actions">
			<button class="btn-cancel" on:click={onCancel}>Отмена</button>
			<button class="btn-confirm" on:click={onConfirm}>Да, закрыть</button>
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
		background: #fff;
		border-radius: 10px;
		padding: 24px;
		min-width: 320px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.16);
	}

	.message {
		font-size: 14px;
		color: #374151;
		margin: 0 0 20px;
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
	}

	.btn-cancel {
		padding: 8px 16px;
		background: none;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 14px;
		color: #374151;
		cursor: pointer;
		transition: background 0.1s;
	}

	.btn-cancel:hover {
		background: #f9fafb;
	}

	.btn-confirm {
		padding: 8px 16px;
		background: #ef4444;
		border: none;
		border-radius: 6px;
		font-size: 14px;
		color: #fff;
		cursor: pointer;
		transition: background 0.1s;
	}

	.btn-confirm:hover {
		background: #dc2626;
	}
</style>
