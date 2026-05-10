<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { LinkItem } from '$lib/types';
	import { getFaviconCandidates } from '$lib/favicon';

	type MatchSegment = { text: string; match: boolean };

	type SearchResult = {
		link: LinkItem;
		nameSegments: MatchSegment[];
		descSegments: MatchSegment[];
		score: number;
	};

	let { links = [] as LinkItem[] } = $props();

	const dispatch = createEventDispatcher<{ select: LinkItem }>();

	let query = $state('');
	let open = $state(false);
	let inputRef: HTMLInputElement | null = null;

	onMount(() => {
		const onKeydown = (e: KeyboardEvent) => {
			const target = e.target as HTMLElement | null;

			const isEditable =
				target &&
				(target.tagName === 'INPUT' ||
					target.tagName === 'TEXTAREA' ||
					target.isContentEditable);

			if (e.key === '/' && !isEditable) {
				e.preventDefault();
				openSearch();
			}

			if (e.key === 'Escape' && open) {
				e.preventDefault();
				closeSearch();
			}
		};

		document.addEventListener('keydown', onKeydown);

		return () => {
			document.removeEventListener('keydown', onKeydown);
		};
	});

	function openSearch() {
		open = true;
		requestAnimationFrame(() => inputRef?.focus());
	}

    function closeSearch() {
        open = false;
        query = '';
        inputRef?.blur();
    }

	function normalizeText(value: string) {
		return value.toLowerCase().trim();
	}

	function fuzzyMatch(text: string, q: string) {
		const t = normalizeText(text);
		const queryText = normalizeText(q);

		let ti = 0;
		let qi = 0;

		const indices: number[] = [];

		while (ti < t.length && qi < queryText.length) {
			if (t[ti] === queryText[qi]) {
				indices.push(ti);
				qi += 1;
			}

			ti += 1;
		}

		if (qi !== queryText.length) return null;

		const span =
			indices.length > 0
				? indices[indices.length - 1] - indices[0]
				: 0;

		const score = span + t.length * 0.1;

		return { indices, score };
	}

	function buildSegments(text: string, indices: number[]) {
		if (!indices.length) {
			return [{ text, match: false }];
		}

		const segments: MatchSegment[] = [];

		let lastIndex = 0;

		for (const idx of indices) {
			if (idx > lastIndex) {
				segments.push({
					text: text.slice(lastIndex, idx),
					match: false
				});
			}

			segments.push({
				text: text[idx],
				match: true
			});

			lastIndex = idx + 1;
		}

		if (lastIndex < text.length) {
			segments.push({
				text: text.slice(lastIndex),
				match: false
			});
		}

		return segments;
	}

	function getLinkIcon(link: LinkItem) {
		return (
			link.customIcon ||
			link.faviconUrl ||
			getFaviconCandidates(link.url)[0] ||
			''
		);
	}

	const results = $derived.by(() => {
		const q = query.trim();

		if (!q) return [];

		const list: SearchResult[] = [];

		for (const link of links) {
			const nameMatch = fuzzyMatch(link.name, q);

			const descMatch = link.description
				? fuzzyMatch(link.description, q)
				: null;

			if (nameMatch || descMatch) {
				const score = Math.min(
					nameMatch?.score ?? Infinity,
					descMatch?.score ?? Infinity
				);

				list.push({
					link,
					nameSegments: buildSegments(
						link.name,
						nameMatch?.indices ?? []
					),
					descSegments: link.description
						? buildSegments(
								link.description,
								descMatch?.indices ?? []
							)
						: [],
					score
				});
			}
		}

		return list
			.sort((a, b) => a.score - b.score)
			.slice(0, 20);
	});
</script>

<div class="relative z-40 mb-5 max-w-xl">
	<div class={`relative transition ${open ? 'scale-[1.02]' : ''}`}>
		<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
			<svg
				viewBox="0 0 24 24"
				class="w-4 h-4"
				fill="currentColor"
				aria-hidden="true"
			>
				<path
					d="M10 2a8 8 0 105.293 14.293l4.207 4.207 1.414-1.414-4.207-4.207A8 8 0 0010 2zm0 2a6 6 0 110 12 6 6 0 010-12z"
				/>
			</svg>
		</span>

		<input
			class="w-full input-nerd !pl-9 !pr-12 transition"
			placeholder="search"
			bind:value={query}
			onfocus={openSearch}
			bind:this={inputRef}
		/>

		<span
			class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-gray-400 border border-gray-700 rounded px-1.5 py-0.5"
		>
			/
		</span>
	</div>

	{#if open && query.trim()}
        <div
            class="absolute mt-2 w-full rounded-xl border border-cyan-500/30 bg-gray-950/95 shadow-xl p-2 max-h-80 overflow-auto"
        >
			{#if query && results.length === 0}
				<p class="text-gray-500 text-sm px-2 py-2">
					No results.
				</p>
			{/if}

			{#each results as result (result.link.id)}
				<button
					class="w-full text-left px-2 py-2 rounded-lg hover:bg-cyan-500/10 transition flex items-center gap-3"
					onclick={() => {
						closeSearch();
						dispatch('select', result.link);
					}}
				>
					{#if getLinkIcon(result.link)}
						<img
							src={getLinkIcon(result.link)}
							alt=""
							class="w-5 h-5 rounded-sm"
							draggable="false"
						/>
					{:else}
						<div
							class="w-5 h-5 rounded-sm bg-cyan-500/20 border border-cyan-500/30"
						/>
					{/if}

					<div class="min-w-0 flex-1">
						<div class="truncate">
							{#each result.nameSegments as segment}
								<span
									class={segment.match
										? 'text-cyan-200 font-semibold'
										: ''}
								>
									{segment.text}
								</span>
							{/each}
						</div>

						{#if result.link.description}
							<div class="text-xs text-gray-400 truncate">
								{#each result.descSegments as segment}
									<span
										class={segment.match
											? 'text-cyan-200 font-semibold'
											: ''}
									>
										{segment.text}
									</span>
								{/each}
							</div>
						{/if}
					</div>
				</button>
			{/each}
		</div>
	{/if}
</div>

{#if open}
	<div
		class="fixed inset-0 bg-black/40 backdrop-blur-sm z-30"
		onclick={closeSearch}
	></div>
{/if}
