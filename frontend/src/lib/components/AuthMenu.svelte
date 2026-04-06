<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher<{
		logout: void;
		register: void;
	}>();

	let open = $state(false);
	let menuRef: HTMLDivElement | null = null;

	function toggle() {
		open = !open;
	}

	function close() {
		open = false;
	}

	function handleClickOutside(e: MouseEvent) {
		if (!menuRef) return;
		if (!menuRef.contains(e.target as Node)) close();
	}

	function handleEscape(e: KeyboardEvent) {
		if (e.key === 'Escape') close();
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		document.addEventListener('keydown', handleEscape);
		return () => {
			document.removeEventListener('click', handleClickOutside);
			document.removeEventListener('keydown', handleEscape);
		};
	});
</script>

<div class="relative" bind:this={menuRef}>
	<button
		type="button"
		class="btn-ghost min-w-20"
		onclick={(e) => {
			e.stopPropagation();
			toggle();
		}}
		aria-haspopup="menu"
		aria-expanded={open}
	>
		auth
	</button>

	{#if open}
		<div
			class="absolute right-0 mt-2 w-56 rounded-xl border border-cyan-500/30 bg-gray-950/95 backdrop-blur p-2 z-50
			       animate-[fadeIn_.14s_ease-out]"
			role="menu"
		>
			<button
				type="button"
				class="w-full text-left px-3 py-2 rounded-lg text-gray-200 hover:bg-cyan-500/15 transition"
				role="menuitem"
				onclick={() => {
					close();
					dispatch('register');
				}}
			>
				Register new account
			</button>

			<button
				type="button"
				class="w-full text-left px-3 py-2 rounded-lg text-red-300 hover:bg-red-500/15 transition"
				role="menuitem"
				onclick={() => {
					close();
					dispatch('logout');
				}}
			>
				Logout
			</button>
		</div>
	{/if}
</div>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-4px) scale(0.98);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}
</style>
