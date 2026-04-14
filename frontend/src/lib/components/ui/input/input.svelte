<script lang="ts">
	export let value: string = '';
	export let placeholder: string = '';
	export let disabled: boolean = false;
	export let id: string | undefined = undefined;
	export let autocomplete: string | undefined = undefined;

	let className = '';
	export { className as class };
</script>

<input
	{id}
	{placeholder}
	{disabled}
	{autocomplete}
	bind:value
	class="input {className}"
	on:blur
	on:focus
	on:input
	{...$$restProps}
/>

<style>
	.input {
		display: block;
		width: 100%;
		padding: 9px 12px;
		background: hsl(var(--card));
		color: hsl(var(--foreground));
		border: 1px solid hsl(var(--border));
		border-radius: calc(var(--radius) - 2px);
		font-size: 14px;
		font-family: inherit;
		outline: none;
		transition: border-color 0.15s, box-shadow 0.15s;
	}

	.input:focus {
		border-color: hsl(var(--primary));
		box-shadow: 0 0 0 2px hsl(var(--primary) / 0.15);
	}

	.input:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.input::placeholder {
		color: hsl(var(--muted-foreground));
	}
</style>
