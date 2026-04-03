<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { isAuthenticated, logout } from '$lib/auth';
	import { api } from '$lib/api';
	import type { LinkItem } from '$lib/types';
	import LinkCard from '$lib/components/LinkCard.svelte';

	let links = $state<LinkItem[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		if (!isAuthenticated()) {
			await goto('/auth');
			return;
		}
		await loadLinks();
	});

	async function loadLinks() {
		loading = true;
		error = '';
		try {
			links = await api.getLinks(true);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load links';
		} finally {
			loading = false;
		}
	}

	async function openLink(link: LinkItem) {
		try {
			await api.registerClick(link.id);
            await loadLinks();
		} catch {
		}
		window.open(link.url, '_blank', 'noopener,noreferrer');
	}

	async function doLogout() {
		logout();
		await goto('/auth');
	}
</script>

<main class="min-h-screen px-4 py-8 max-w-5xl mx-auto">
	<header class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-cyan-300">[ dashboard ]</h1>
			<p class="text-gray-400 text-sm">read://links</p>
		</div>
		<button class="btn-ghost" onclick={doLogout}>Logout</button>
	</header>

	{#if loading}
		<p class="text-gray-400">Loading links...</p>
	{:else if error}
		<p class="text-red-400">{error}</p>
	{:else if links.length === 0}
		<p class="text-gray-400">No links yet.</p>
	{:else}
		<div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each links as link (link.id)}
				<LinkCard {link} onOpen={openLink} />
			{/each}
		</div>
	{/if}
</main>
