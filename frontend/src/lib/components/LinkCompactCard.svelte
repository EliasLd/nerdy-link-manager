<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import type { LinkItem } from '$lib/types';
	import { getFaviconCandidates } from '$lib/favicon';

	let { link }: { link: LinkItem } = $props();

    const dispatch = createEventDispatcher<{
        open: LinkItem;
        details: LinkItem;
        edit: LinkItem;
        delete: LinkItem;
        addToFolder: LinkItem;
    }>();

	let menuOpen = $state(false);
	let rootRef: HTMLDivElement | null = null;
	let menuRef: HTMLDivElement | null = null;

	let menuStyle = $state('right: 0.375rem; top: 1.75rem;'); // default

	let faviconSrc = $state(getFaviconCandidates(link.url)[0] ?? '');
	let faviconIndex = $state(0);
	const candidates = getFaviconCandidates(link.url);

	function onFaviconError() {
		faviconIndex += 1;
		faviconSrc = candidates[faviconIndex] ?? '';
	}

	async function openMenu() {
		menuOpen = true;
		await tick();
		positionMenu();
	}

	function closeMenu() {
		menuOpen = false;
	}

	function toggleMenu() {
		if (menuOpen) closeMenu();
		else openMenu();
	}

	function positionMenu() {
		if (!rootRef || !menuRef) return;

		const root = rootRef.getBoundingClientRect();
		const menu = menuRef.getBoundingClientRect();
		const pad = 8;

		// Horizontal: prefer right align, flip to left if overflow
		let left = root.width - menu.width - 6; // like right:6px
		if (root.left + left < pad) left = 6;
		if (root.left + left + menu.width > window.innerWidth - pad) {
			left = Math.max(pad - root.left, window.innerWidth - pad - root.left - menu.width);
		}

		// Vertical: prefer below button, flip above if overflow
		let top = 28; // below 3 dots
		if (root.top + top + menu.height > window.innerHeight - pad) {
			top = Math.max(6, root.height - menu.height - 6); // open upward inside card bounds
		}

		menuStyle = `left:${left}px; top:${top}px;`;
	}

	function handleClickOutside(e: MouseEvent) {
		if (!rootRef) return;
		if (!rootRef.contains(e.target as Node)) closeMenu();
	}

	function handleEscape(e: KeyboardEvent) {
		if (e.key === 'Escape') closeMenu();
	}

	function handleResizeOrScroll() {
		if (menuOpen) positionMenu();
	}

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
    bind:this={rootRef}
>
	<button
		class="absolute top-1.5 right-1.5 text-gray-400 hover:text-cyan-300 transition text-sm leading-none"
		onclick={toggleMenu}
		aria-label="Open actions menu"
	>
		⋯
	</button>

	{#if menuOpen}
		<div
			bind:this={menuRef}
			class="absolute z-[60] w-40 rounded-lg border border-cyan-500/30 bg-gray-950/95 p-1 shadow-xl"
			style={menuStyle}
		>
		<button
            class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded"
            onclick={() => {
                closeMenu();
                dispatch('details', link);
            }}
        >
            Open details
        </button>

        <button
            class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded"
            onclick={() => {
                closeMenu();
                dispatch('edit', link);
            }}
        >
            Edit
        </button>

        <button
            class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded"
            onclick={() => {
                closeMenu();
                dispatch('addToFolder', link);
            }}
        >
            Add to folder
</button>

<button
	class="w-full text-left px-2 py-1.5 hover:bg-red-500/15 text-red-300 rounded"
	onclick={() => {
		closeMenu();
		dispatch('delete', link);
	}}
>
	Delete
</button>        </div>
	{/if}

	<button class="w-full aspect-square grid place-items-center" onclick={() => dispatch('open', link)}>
    {#if faviconSrc}
			<img src={faviconSrc} alt="" class="w-7 h-7 rounded-sm" onerror={onFaviconError} />
		{:else}
			<div class="w-7 h-7 rounded-sm bg-cyan-500/20 border border-cyan-500/30" />
		{/if}
	</button>
</div>
