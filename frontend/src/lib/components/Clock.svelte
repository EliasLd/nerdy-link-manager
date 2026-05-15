<script lang="ts">
	import { onMount } from 'svelte';

	let time = $state('');

	function formatTime(d: Date) {
		return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
	}

	function update() {
		time = formatTime(new Date());
	}

	onMount(() => {
		update();
		const id = setInterval(update, 1000 * 30);
		return () => clearInterval(id);
	});
</script>

<div class="font-mono text-cyan-300/90 text-lg md:text-2xl font-semibold tracking-wide">
	{time}
</div>
