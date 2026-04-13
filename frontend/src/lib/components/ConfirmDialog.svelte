<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher<{ close: void; confirm: void }>();

	let { title = 'Confirm action', message = 'Are you sure?' }: { title?: string; message?: string } = $props();

	function close() {
		dispatch('close');
	}
	function confirm() {
		dispatch('confirm');
	}
	function onBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) close();
	}
	function onEscape(e: KeyboardEvent) {
		if (e.key === 'Escape') close();
	}

	onMount(() => {
		document.addEventListener('keydown', onEscape);
		return () => document.removeEventListener('keydown', onEscape);
	});
</script>

<div class="fixed inset-0 z-50 bg-black/55 backdrop-blur-[1px] grid place-items-center px-4" onclick={onBackdropClick}>
	<div class="w-full max-w-md border border-red-500/30 bg-gray-950/95 rounded-2xl p-6 animate-[modalIn_.18s_ease-out]">
		<h2 class="text-xl font-bold text-red-300 mb-2">{title}</h2>
		<p class="text-sm text-gray-300">{message}</p>

		<div class="mt-5 flex items-center justify-end gap-2">
			<button class="btn-ghost" type="button" onclick={close}>Cancel</button>
			<button
				type="button"
				class="px-4 py-2 rounded-lg border border-red-500/35 text-red-300 hover:bg-red-500/15 transition"
				onclick={confirm}
			>
				Delete
			</button>
		</div>
	</div>
</div>

<style>
	@keyframes modalIn {
		from { opacity: 0; transform: translateY(8px) scale(0.98); }
		to { opacity: 1; transform: translateY(0) scale(1); }
	}
</style>
