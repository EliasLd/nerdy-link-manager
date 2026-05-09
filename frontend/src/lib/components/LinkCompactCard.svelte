<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import type { LinkItem } from '$lib/types';
	import { getFaviconCandidates } from '$lib/favicon';
	import { api } from '$lib/api';

	let { link }: { link: LinkItem } = $props();

    const dispatch = createEventDispatcher<{
        open: LinkItem;
        details: LinkItem;
        edit: LinkItem;
        delete: LinkItem;
        addToFolder: LinkItem;
        removeFromFolder: LinkItem;
        iconError: string;
        iconUpdated: void;
    }>();

	const MAX_ICON_SIZE = 64;
	const MAX_ICON_BYTES = 200 * 1024;

	let menuOpen = $state(false);
	let rootRef: HTMLDivElement | null = null;
	let menuRef: HTMLDivElement | null = null;
	let fileInput: HTMLInputElement | null = null;

	let menuStyle = $state('right: 0.375rem; top: 1.75rem;');

	const baseCandidates = getFaviconCandidates(link.url);
	const candidates = link.customIcon
		? [link.customIcon]
		: [
				...(link.faviconUrl ? [link.faviconUrl] : []),
				...baseCandidates.filter((c) => c !== link.faviconUrl)
		  ];

	let faviconSrc = $state(candidates[0] ?? '');
	let faviconIndex = $state(0);

	function onFaviconError() {
		if (link.customIcon) return;
		faviconIndex += 1;
		faviconSrc = candidates[faviconIndex] ?? '';
	}

    async function onFaviconLoad() {
        if (link.customIcon || !faviconSrc) return;

        if (faviconSrc !== link.faviconUrl) {
            try {
                await api.updateLink(link.id, {
                    name: link.name,
                    url: link.url,
                    description: link.description ?? undefined,
                    folderId: link.folderId ? Number(link.folderId) : null,
                    faviconUrl: faviconSrc
                });
            } catch {}
        }
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

	function handleEscape(e: KeyboardEvent) {
		if (e.key === 'Escape') closeMenu();
	}

	function handleResizeOrScroll() {
		if (menuOpen) positionMenu();
	}

	function fileToDataUrl(file: File): Promise<string> {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = () => resolve(String(reader.result));
			reader.onerror = () => reject(reader.error);
			reader.readAsDataURL(file);
		});
	}

	async function resizeImage(file: File): Promise<string> {
		const dataUrl = await fileToDataUrl(file);
		const img = new Image();

		await new Promise<void>((resolve, reject) => {
			img.onload = () => resolve();
			img.onerror = () => reject(new Error('Invalid image'));
			img.src = dataUrl;
		});

		const scale = Math.min(1, MAX_ICON_SIZE / Math.max(img.width, img.height));
		const w = Math.round(img.width * scale);
		const h = Math.round(img.height * scale);

		const canvas = document.createElement('canvas');
		canvas.width = w;
		canvas.height = h;

		const ctx = canvas.getContext('2d');
		if (!ctx) throw new Error('Canvas error');
		ctx.drawImage(img, 0, 0, w, h);

		return canvas.toDataURL('image/png', 0.92);
	}

	function onDragStart(e: DragEvent) {
		if (!e.dataTransfer) return;
		e.dataTransfer.setData('text/plain', link.id);
		e.dataTransfer.effectAllowed = 'move';
	}

    async function onIconFile(e: Event) {
        const input = e.currentTarget as HTMLInputElement;
        const file = input.files?.[0];
        input.value = '';
        if (!file) return;

        if (!file.type.startsWith('image/')) {
            dispatch('iconError', 'Selected file is not an image.');
            return;
        }
        if (file.size > MAX_ICON_BYTES) {
            dispatch('iconError', 'Image is too large. Max 200KB.');
            return;
        }

        try {
            const dataUrl = await resizeImage(file);
            await api.updateLink(link.id, {
                name: link.name,
                url: link.url,
                description: link.description ?? undefined,
                folderId: link.folderId ? Number(link.folderId) : null,
                customIcon: dataUrl
            });
        } catch {
            dispatch('iconError', 'Failed to process image.');
        }
        dispatch('iconUpdated');
    }

	async function removeIcon() {
        try {
            await api.updateLink(link.id, {
                name: link.name,
                url: link.url,
                description: link.description ?? undefined,
                folderId: link.folderId ? Number(link.folderId) : null,
                customIcon: null
            });
        } catch {
            dispatch('iconError', 'Failed to remove custom icon.');
        }
        dispatch('iconUpdated');
    }

	function triggerIconPicker() {
		fileInput?.click();
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

<input type="file" accept="image/*" class="hidden" bind:this={fileInput} onchange={onIconFile} />

<div
	class="relative border border-cyan-500/25 rounded-lg bg-black/30 p-1.5 hover:border-cyan-300/60 transition overflow-visible"
	bind:this={rootRef}
	draggable="true"
	ondragstart={onDragStart}
>	
    <button
		class="absolute top-1.5 right-1.5 text-gray-400 hover:text-cyan-300 transition text-lg leading-none"
		onclick={toggleMenu}
		aria-label="Open actions menu"
	>
		⋯
	</button>

	{#if menuOpen}
		<div
			bind:this={menuRef}
			class="absolute z-[60] w-44 rounded-lg border border-cyan-500/30 bg-gray-950/95 p-1 shadow-xl"
			style={menuStyle}
		>
			<button class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded" onclick={() => { closeMenu(); dispatch('details', link); }}>Open details</button>
			<button class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded" onclick={() => { closeMenu(); dispatch('edit', link); }}>Edit</button>
			<button class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded" onclick={() => { closeMenu(); dispatch('addToFolder', link); }}>Add to folder</button>

			<button class="w-full text-left px-2 py-1.5 hover:bg-cyan-500/15 rounded" onclick={() => { closeMenu(); triggerIconPicker(); }}>
				Set custom icon
			</button>

			{#if link.customIcon}
				<button class="w-full text-left px-2 py-1.5 hover:bg-red-500/15 text-red-300 rounded" onclick={() => { closeMenu(); removeIcon(); }}>
					Remove custom icon
				</button>
			{/if}

			{#if link.folderId}
				<button class="w-full text-left px-2 py-1.5 hover:bg-red-500/15 text-red-300 rounded" onclick={() => { closeMenu(); dispatch('removeFromFolder', link); }}>
					Remove from folder
				</button>
			{/if}

			<button class="w-full text-left px-2 py-1.5 hover:bg-red-500/15 text-red-300 rounded" onclick={() => { closeMenu(); dispatch('delete', link); }}>
				Delete
			</button>
		</div>
	{/if}

	<button class="w-full aspect-square grid place-items-center" onclick={() => dispatch('open', link)}>
		{#if faviconSrc}
			<img
				src={faviconSrc}
				alt=""
				class="w-10 h-10 rounded-sm"
				onerror={onFaviconError}
				onload={onFaviconLoad}
				draggable="false"
			/>
		{:else}
			<div class="w-7 h-7 rounded-sm bg-cyan-500/20 border border-cyan-500/30" draggable="false" />
		{/if}
	</button>

	<p class="absolute bottom-1 left-1 right-1 text-center text-[11px] text-cyan-300 truncate px-1 bg-black/60 rounded">
		{link.name}
	</p>
</div>
