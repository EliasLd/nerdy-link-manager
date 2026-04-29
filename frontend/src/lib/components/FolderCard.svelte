<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { FolderItem } from '$lib/types';

	let { folder }: { folder: FolderItem } = $props();
	const dispatch = createEventDispatcher<{ open: FolderItem }>();
</script>

<div
	class="relative border border-cyan-500/25 rounded-lg bg-black/30 p-1.5 hover:border-cyan-300/60 transition overflow-hidden"
	role="button"
	tabindex="0"
	onclick={() => dispatch('open', folder)}
>
	<div class="w-full aspect-square grid place-items-center">
		<div class="w-8 h-8 rounded-lg border border-cyan-500/30 text-cyan-300 grid place-items-center">
			<svg viewBox="0 0 24 24" class="w-4 h-4" fill="currentColor" aria-hidden="true">
				<path d="M3 6.75A2.75 2.75 0 015.75 4h4.19c.73 0 1.43.29 1.95.81l.8.8c.19.19.45.3.72.3h4.84A2.75 2.75 0 0121 8.66v7.59A2.75 2.75 0 0118.25 19H5.75A2.75 2.75 0 013 16.25V6.75z"/>
			</svg>
		</div>
	</div>

	<!-- label overlay -->
	<p class="absolute bottom-1 left-1 right-1 text-center text-[11px] text-cyan-300 truncate px-1">
		{folder.name}
	</p>
</div>
