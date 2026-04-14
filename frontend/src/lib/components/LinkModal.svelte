<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { LinkItem } from '$lib/types';

	const dispatch = createEventDispatcher<{
		close: void;
		save: { name: string; url: string; description?: string };
	}>();

	let {
		mode = 'create',
		initial
	}: { mode?: 'create' | 'edit'; initial?: LinkItem | null } = $props();

	let name = $state(initial?.name ?? '');
	let url = $state(initial?.url ?? '');
	let description = $state(initial?.description ?? '');
	let error = $state('');

	function close() {
		dispatch('close');
	}

	function submit() {
		error = '';
		if (!name.trim()) {
			error = 'Name is required';
			return;
		}
		if (!url.trim()) {
			error = 'URL is required';
			return;
		}

		dispatch('save', {
			name: name.trim(),
			url: url.trim(),
			description: description.trim() || undefined
		});
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
	<div class="w-full max-w-lg border border-cyan-500/30 bg-gray-950/95 rounded-2xl p-6 animate-[modalIn_.18s_ease-out]">
		<h2 class="text-xl font-bold text-cyan-300 mb-1">
			{mode === 'create' ? 'Add link' : 'Edit link'}
		</h2>
		<p class="text-sm text-gray-400 mb-5">
			{mode === 'create' ? '[create] insert://new-link' : '[edit] patch://existing-link'}
		</p>

		<div class="space-y-3">
			<input class="w-full input-nerd" placeholder="name" bind:value={name} />
			<input class="w-full input-nerd" placeholder="https://example.com" bind:value={url} />
			<textarea class="w-full input-nerd min-h-24" placeholder="description (optional)" bind:value={description}></textarea>
		</div>

		{#if error}<p class="text-red-400 text-sm mt-3">{error}</p>{/if}

		<div class="mt-5 flex items-center justify-end gap-2">
			<button class="btn-ghost" type="button" onclick={close}>Cancel</button>
			<button class="btn-primary" type="button" onclick={submit}>
				{mode === 'create' ? 'Create' : 'Save changes'}
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
