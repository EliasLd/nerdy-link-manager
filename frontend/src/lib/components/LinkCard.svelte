<script lang="ts">
	import type { LinkItem } from '$lib/types';

	let { link, onOpen, onEdit, onDelete } = $props<{
		link: LinkItem;
		onOpen: (link: LinkItem) => void;
		onEdit?: (link: LinkItem) => void;
		onDelete?: (link: LinkItem) => void;
	}>();

	function handleCardClick() {
		onOpen(link);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			onOpen(link);
		}
	}
</script>

<div
	role="button"
	tabindex="0"
	class="group border border-cyan-500/25 rounded-xl px-4 py-3 bg-black/30 cursor-pointer
	transition-all duration-300 ease-out hover:border-cyan-300/70 hover:bg-black/40
	grid grid-rows-[1fr_auto] gap-3 min-h-[120px]"
	onclick={handleCardClick}
	onkeydown={handleKeydown}
>
	<div class="flex items-start justify-between gap-3 min-h-0">
		<div class="min-w-0 flex-1">
			<h3 class="text-cyan-300 font-semibold truncate">{link.name}</h3>
			<p class="text-sm text-cyan-500/90 truncate mt-1">{link.url}</p>
			{#if link.description}
				<p class="text-gray-400 text-sm mt-2 line-clamp-2">{link.description}</p>
			{/if}
		</div>

		{#if onEdit || onDelete}
			<div class="flex items-center gap-2 shrink-0">
				{#if onEdit}
					<button
						type="button"
						class="px-2 py-1 text-xs rounded border border-cyan-500/30 text-cyan-300 hover:bg-cyan-500/15 transition"
						onclick={(e) => { e.stopPropagation(); onEdit?.(link); }}
					>
						Edit
					</button>
				{/if}
				{#if onDelete}
					<button
						type="button"
						class="px-2 py-1 text-xs rounded border border-red-500/35 text-red-300 hover:bg-red-500/15 transition"
						onclick={(e) => { e.stopPropagation(); onDelete?.(link); }}
					>
						Delete
					</button>
				{/if}
			</div>
		{/if}
	</div>

	<div class="flex items-center justify-between text-xs font-mono text-gray-400">
		<span>clicks: <span class="text-cyan-300">{link.clicks ?? 0}</span></span>
		<span class="text-cyan-500/70 group-hover:text-cyan-300 transition-colors">open ↗</span>
	</div>
</div>
