<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher<{ close: void; submit: { current: string; next: string } }>();

	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let error = $state('');

	function close() { dispatch('close'); }

	function submit() {
		error = '';
		if (!currentPassword.trim() || !newPassword.trim() || !confirmPassword.trim()) {
			error = 'All fields are required';
			return;
		}
		if (newPassword !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}
		dispatch('submit', { current: currentPassword, next: newPassword });
	}

	function onBackdropClick(e: MouseEvent) { if (e.target === e.currentTarget) close(); }
	function onEscape(e: KeyboardEvent) { if (e.key === 'Escape') close(); }

	onMount(() => {
		document.addEventListener('keydown', onEscape);
		return () => document.removeEventListener('keydown', onEscape);
	});
</script>

<div class="fixed inset-0 z-50 bg-black/55 backdrop-blur-[1px] grid place-items-center px-4" onclick={onBackdropClick}>
	<div class="w-full max-w-md border border-cyan-500/30 bg-gray-950/95 rounded-2xl p-6 animate-[modalIn_.18s_ease-out]">
		<h2 class="text-xl font-bold text-cyan-300 mb-1">Change password</h2>
		<p class="text-sm text-gray-400 mb-5">[auth] update://password</p>

		<div class="space-y-3">
			<input class="w-full input-nerd" placeholder="current password" type="password" bind:value={currentPassword} />
			<input class="w-full input-nerd" placeholder="new password" type="password" bind:value={newPassword} />
			<input class="w-full input-nerd" placeholder="confirm new password" type="password" bind:value={confirmPassword} />
		</div>

		{#if error}<p class="text-red-400 text-sm mt-3">{error}</p>{/if}

		<div class="mt-5 flex items-center justify-end gap-2">
			<button class="btn-ghost" type="button" onclick={close}>Cancel</button>
			<button class="btn-primary" type="button" onclick={submit}>Update</button>
		</div>
	</div>
</div>
