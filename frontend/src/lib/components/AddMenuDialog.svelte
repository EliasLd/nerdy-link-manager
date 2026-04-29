<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher<{ close: void; addLink: void; addFolder: void }>();

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
		<h2 class="text-xl font-bold text-cyan-300 mb-1">Create</h2>
		<p class="text-sm text-gray-400 mb-5">[add] choose target</p>

		<div class="space-y-2">
			<button class="btn-primary w-full" onclick={() => dispatch('addLink')}>+ Add link</button>
			<button class="btn-ghost w-full" onclick={() => dispatch('addFolder')}>+ Add folder</button>
		</div>
	</div>
</div>
