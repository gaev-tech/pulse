<script lang="ts">
	export let message: string = 'Нет данных';
</script>

<div class="zero-state">
	<span class="icon">·</span>
	<p class="message">{message}</p>
</div>

<style>
	.zero-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 16px 0;
		gap: 6px;
	}

	.icon {
		font-size: 20px;
		color: rgba(205, 214, 244, 0.2);
	}

	.message {
		font-size: 12px;
		color: rgba(205, 214, 244, 0.35);
		margin: 0;
	}
</style>
