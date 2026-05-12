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
        class="btn-ghost !px-3 !py-2 sm:min-w-20 flex items-center justify-center gap-2"
        onclick={(e) => {
            e.stopPropagation();
            toggle();
        }}
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label="Auth menu"
    >
        <span class="hidden sm:inline">auth</span>
        <span class="sm:hidden inline-flex">
            <svg viewBox="0 0 24 24" class="w-5 h-5" fill="currentColor" aria-hidden="true">
                <path d="M3 3h12a2 2 0 0 1 2 2v5h-2V5H5v14h10v-3h2v3a2 2 0 0 1-2 2H3V3z"/>
                <path d="M14 12l-2 2 2 2v-1h6v-2h-6v-1z"/>
            </svg>
        </span>
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
				class="w-full text-left px-3 py-2 rounded-lg text-gray-200 hover:bg-cyan-500/15 transition"
				role="menuitem"
				onclick={() => {
					close();
					dispatch('changePassword');
				}}
			>
				Change password
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
