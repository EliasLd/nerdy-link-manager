<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import type { FolderItem } from '$lib/types';

	let { folder }: { folder: FolderItem } = $props();

	const dispatch = createEventDispatcher<{
		open: FolderItem;
		rename: FolderItem;
		delete: FolderItem;
	}>();

	let menuOpen = $state(false);
	let rootRef: HTMLDivElement | null = null;
	let menuRef: HTMLDivElement | null = null; let menuStyle = $state('right: 0.375rem; top: 1.75rem;');

	async function openMenu() {
		menuOpen = true;
		await tick();
		positionMenu();
	}
	function closeMenu() { menuOpen = false; }
	function toggleMenu() { menuOpen ? closeMenu() : openMenu(); }

	function positionMenu() {
		if (!rootRef || !menuRef) return;
		const root = rootRef.getBoundingClientRect();
		const menu = menuRef.getBoundingClientRect();
		const pad = 8;

		let left = root.width - menu.width - 6;
		if (root.left + left < pad) left = 6;
		if (root.left + left + menu.width > window.innerWidth - pad) {
			left = Math.max(pad - root.left, window.innerWidth - pad - root.left - menu.width);
		}

		let top = 28;
		if (root.top + top + menu.height > window.innerHeight - pad) {
			top = Math.max(6, root.height - menu.height - 6);
		}

		menuStyle = `left:${left}px; top:${top}px;`;
	}

	function handleClickOutside(e: MouseEvent) {
		if (!rootRef) return;
		if (!rootRef.contains(e.target as Node)) closeMenu();
	}
	function handleEscape(e: KeyboardEvent) { if (e.key === 'Escape') closeMenu(); }
	function handleResizeOrScroll() { if (menuOpen) positionMenu(); }

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		document.addEventListener('keydown', handleEscape);
		window.addEventListener('resize', handleResizeOrScroll);
		window.addEventListener('scroll', handleResizeOrScroll, true);
		return () => {
			document.removeEventListener('click', handleClickOutside);
			document.removeEventListener('keydown', handleEscape);
			window.removeEventListener('resize', handleResizeOrScroll);
			window.removeEventListener('scroll', handleResizeOrScroll, true);
		};
	});
</script>

<div
	class="relative border border-cyan-500/25 rounded-lg bg-black/30 p-1.5 hover:border-cyan-300/60 transition overflow-visible"
	role="button"
	tabindex="0"
	bind:this={rootRef}
	onclick={() => dispatch('open', folder)}
>
	<button
		class="absolute top-1.5 right-1.5 text-gray-400 hover:text-cyan-300 transition text-sm leading-none"
		onclick={(e) => { e.stopPropagation(); toggleMenu(); }}
		aria-label="Open folder menu"
	>
		⋯
	</button>

	{#if menuOpen}
		<div
			bind:this={menuRef}
			class="absolute z-[60] w-36 rounded-lg border border-cyan-500/30 bg-gray-950/95 p-1 shadow-xl"
			style={menuStyle}
		>
			<button
				class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded"
				onclick={() => { closeMenu(); dispatch('rename', folder); }}
			>
				Rename
			</button>

			<button
				class="w-full text-left px-2 py-1.5 hover:bg-red-500/15 text-red-300 rounded"
				onclick={() => { closeMenu(); dispatch('delete', folder); }}
			>
				Delete
			</button>
		</div>
	{/if}

    <div class="w-full aspect-square grid place-items-center">
        <div class="w-12 h-12 rounded-lg border border-cyan-500/30 text-cyan-300 grid place-items-center">
            <svg viewBox="0 0 24 24" class="w-9 h-9" fill="currentColor" aria-hidden="true">
                <path d="M3 6.75A2.75 2.75 0 015.75 4h4.19c.73 0 1.43.29 1.95.81l.8.8c.19.19.45.3.72.3h4.84A2.75 2.75 0 0121 8.66v7.59A2.75 2.75 0 0118.25 19H5.75A2.75 2.75 0 013 16.25V6.75z"/>
            </svg>
        </div>
    </div>

	<p class="absolute bottom-1 left-1 right-1 text-center text-[11px] text-cyan-300 truncate px-1 bg-black/60 rounded">
	    {folder.name}
    </p>
</div>
