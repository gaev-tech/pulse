<script lang="ts">
	import { X } from 'lucide-svelte';
	import { toastError } from '../../stores/toast';

	let timer: ReturnType<typeof setTimeout>;

	$: if ($toastError) {
		clearTimeout(timer);
		timer = setTimeout(() => toastError.set(null), 4000);
	}

	function dismiss() {
		clearTimeout(timer);
		toastError.set(null);
	}
</script>

{#if $toastError}
	<div class="error-toast" role="alert">
		<span>{$toastError}</span>
		<button type="button" aria-label="Закрыть" on:click={dismiss}>
			<X size={14} />
		</button>
	</div>
{/if}

<style>
	.error-toast {
		position: fixed;
		bottom: 24px;
		right: 24px;
		z-index: 9999;
		display: flex;
		align-items: center;
		gap: 12px;
		background: hsl(var(--destructive));
		color: hsl(var(--destructive-foreground));
		padding: 12px 16px;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
		cursor: pointer;
		max-width: 360px;
		font-size: 14px;
	}

	button {
		display: flex;
		align-items: center;
		background: none;
		border: none;
		color: hsl(var(--destructive-foreground));
		cursor: pointer;
		padding: 0;
		flex-shrink: 0;
	}
</style>
