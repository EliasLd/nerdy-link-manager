<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { LinkItem } from '$lib/types';

	let { link }: { link: LinkItem } = $props();

	const dispatch = createEventDispatcher<{
		close: void;
		open: LinkItem;
		edit: LinkItem;
		delete: LinkItem;
	}>();

	function close() { dispatch('close'); }
	function onBackdropClick(e: MouseEvent) { if (e.target === e.currentTarget) close(); }
	function onEscape(e: KeyboardEvent) { if (e.key === 'Escape') close(); }

	onMount(() => {
		document.addEventListener('keydown', onEscape);
		return () => document.removeEventListener('keydown', onEscape);
	});
</script>

<div class="fixed inset-0 z-50 bg-black/55 backdrop-blur-[1px] grid place-items-center px-4" onclick={onBackdropClick}>
	<div class="w-full max-w-lg border border-cyan-500/30 bg-gray-950/95 rounded-2xl p-5 animate-[modalIn_.18s_ease-out]">
		<div class="flex items-start justify-between gap-3">
			<div class="min-w-0">
				<h3 class="text-cyan-300 font-semibold truncate">{link.name}</h3>
				<p class="text-sm text-cyan-500/90 break-all mt-1">{link.url}</p>
				{#if link.description}
					<p class="text-gray-400 text-sm mt-2">{link.description}</p>
				{/if}
			</div>
			<button class="btn-ghost !px-2 !py-1" onclick={close}>✕</button>
		</div>

		<div class="flex items-center justify-between text-xs font-mono text-gray-400 mt-4">
			<span>clicks: <span class="text-cyan-300">{link.clicks ?? 0}</span></span>
		</div>

		<div class="mt-5 flex items-center justify-end gap-2">
			<button class="btn-ghost" onclick={() => dispatch('open', link)}>Open</button>
			<button class="btn-ghost" onclick={() => dispatch('edit', link)}>Edit</button>
			<button class="px-4 py-2 rounded-lg border border-red-500/35 text-red-300 hover:bg-red-500/15 transition" onclick={() => dispatch('delete', link)}>Delete</button>
		</div>
	</div>
</div>

<style>
	@keyframes modalIn {
		from { opacity: 0; transform: translateY(8px) scale(0.98); }
		to { opacity: 1; transform: translateY(0) scale(1); }
	}
</style>
