<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { FolderItem } from '$lib/types';

	let { folders = [] as FolderItem[] } = $props();

	const dispatch = createEventDispatcher<{ close: void; choose: { folderId: string | null } }>();

	function close() { dispatch('close'); }
	function onBackdropClick(e: MouseEvent) { if (e.target === e.currentTarget) close(); }
	function onEscape(e: KeyboardEvent) { if (e.key === 'Escape') close(); }

	onMount(() => {
		document.addEventListener('keydown', onEscape);
		return () => document.removeEventListener('keydown', onEscape);
	});
</script>

<div class="fixed inset-0 z-50 bg-black/55 backdrop-blur-[1px] grid place-items-center px-4" onclick={onBackdropClick}>
	<div class="w-full max-w-sm border border-cyan-500/30 bg-gray-950/95 rounded-2xl p-5 animate-[modalIn_.18s_ease-out]">
		<h2 class="text-lg font-bold text-cyan-300 mb-1">Add to folder</h2>
		<p class="text-xs text-gray-400 mb-4">select target folder</p>

		<div class="space-y-2 max-h-72 overflow-auto">
			<button class="btn-ghost w-full text-left" onclick={() => dispatch('choose', { folderId: null })}>
				No folder
			</button>

			{#each folders as folder (folder.id)}
				<button class="btn-ghost w-full text-left" onclick={() => dispatch('choose', { folderId: folder.id })}>
					{folder.name}
				</button>
			{/each}
		</div>
	</div>
</div>
