<script lang="ts">
	export let message: string = 'Нет данных';
</script>

<div class="zero-state">
	<p class="message">{message}</p>
</div>

<style>
	.zero-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 16px 0;
	}

	.message {
		font-size: 12px;
		color: hsl(var(--muted-foreground));
		margin: 0;
	}
</style>
