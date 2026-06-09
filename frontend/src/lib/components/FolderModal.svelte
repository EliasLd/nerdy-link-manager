<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { FolderItem } from '$lib/types';

	const dispatch = createEventDispatcher<{ close: void; save: { name: string } }>();

	let {
		mode = 'create',
		initial
	}: { mode?: 'create' | 'edit'; initial?: FolderItem | null } = $props();

	let name = $state(initial?.name ?? '');
	let error = $state('');
    let nameInput: HTMLInputElement;

	function close() { dispatch('close'); }
	function submit() {
		error = '';
		if (!name.trim()) {
			error = 'Folder name is required';
			return;
		}
		dispatch('save', { name: name.trim() });
	}
	function onBackdropClick(e: MouseEvent) { if (e.target === e.currentTarget) close(); }
	function onEscape(e: KeyboardEvent) { if (e.key === 'Escape') close(); }

	onMount(() => {
        nameInput?.focus();
		document.addEventListener('keydown', onEscape);
		return () => {
            document.removeEventListener('keydown', onEscape);
        };
	});
</script>

<div class="fixed inset-0 z-50 bg-black/55 backdrop-blur-[1px] grid place-items-center px-4" onclick={onBackdropClick}>
	<div class="w-full max-w-lg border border-cyan-500/30 bg-gray-950/95 rounded-2xl p-6 animate-[modalIn_.18s_ease-out]">
		<h2 class="text-xl font-bold text-cyan-300 mb-1">
			{mode === 'create' ? 'Add folder' : 'Rename folder'}
		</h2>
		<p class="text-sm text-gray-400 mb-5">
			{mode === 'create' ? '[create] mkdir://new-folder' : '[edit] mv://folder-name'}
		</p>

		<div class="space-y-3">
			<input bind:this={nameInput} class="w-full input-nerd" placeholder="Folder name" bind:value={name} />
		</div>

		{#if error}<p class="text-red-400 text-sm mt-3">{error}</p>{/if}

		<div class="mt-5 flex items-center justify-end gap-2">
			<button class="btn-ghost" type="button" onclick={close}>Cancel</button>
			<button class="btn-primary" type="button" onclick={submit}>
				{mode === 'create' ? 'Create' : 'Save'}
			</button>
		</div>
	</div>
</div>
