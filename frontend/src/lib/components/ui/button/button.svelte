<script lang="ts">
	type Variant = 'default' | 'destructive' | 'outline' | 'ghost';
	type Size = 'default' | 'sm' | 'lg';

	export let variant: Variant = 'default';
	export let size: Size = 'default';
	export let disabled: boolean = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';

	let className = '';
	export { className as class };
</script>

<button
	{type}
	{disabled}
	class="btn btn-{variant} btn-size-{size} {className}"
	on:click
	{...$$restProps}
>
	<slot />
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		border-radius: calc(var(--radius) - 2px);
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		border: 1px solid transparent;
		transition: background 0.15s, color 0.15s, border-color 0.15s, opacity 0.15s;
		white-space: nowrap;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	/* Sizes */
	.btn-size-default { padding: 8px 16px; min-height: 36px; }
	.btn-size-sm      { padding: 5px 12px; min-height: 30px; font-size: 13px; }
	.btn-size-lg      { padding: 11px 24px; min-height: 44px; }

	/* Variants */
	.btn-default {
		background: hsl(var(--primary));
		color: hsl(var(--primary-foreground));
	}
	.btn-default:hover:not(:disabled) {
		background: hsl(var(--primary) / 0.9);
	}

	.btn-destructive {
		background: hsl(var(--destructive));
		color: hsl(var(--destructive-foreground));
	}
	.btn-destructive:hover:not(:disabled) {
		background: hsl(var(--destructive) / 0.9);
	}

	.btn-outline {
		background: transparent;
		color: hsl(var(--foreground));
		border-color: hsl(var(--border));
	}
	.btn-outline:hover:not(:disabled) {
		background: hsl(var(--muted));
	}

	.btn-ghost {
		background: transparent;
		color: hsl(var(--foreground));
	}
	.btn-ghost:hover:not(:disabled) {
		background: hsl(var(--muted));
	}
</style>
