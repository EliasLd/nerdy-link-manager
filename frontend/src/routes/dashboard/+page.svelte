<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { isAuthenticated, logout } from '$lib/auth';
	import { api } from '$lib/api';
	import type { LinkItem, CreateLinkPayload, UpdateLinkPayload } from '$lib/types';
	import LinkCompactCard from '$lib/components/LinkCompactCard.svelte';
	import LinkDetailsDialog from '$lib/components/LinkDetailsDialog.svelte';
	import LinkModal from '$lib/components/LinkModal.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import AuthMenu from '$lib/components/AuthMenu.svelte';

	let links = $state<LinkItem[]>([]);
	let loading = $state(true);
	let error = $state('');

	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let showDetailsModal = $state(false);

	let selectedLink = $state<LinkItem | null>(null);

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
		} catch {}
		window.open(link.url, '_blank', 'noopener,noreferrer');
	}

	function openCreate() {
		showCreateModal = true;
	}

	function openDetails(link: LinkItem) {
		selectedLink = link;
		showDetailsModal = true;
	}

	function openEdit(link: LinkItem) {
		selectedLink = link;
		showEditModal = true;
	}

	function openDelete(link: LinkItem) {
		selectedLink = link;
		showDeleteModal = true;
	}

	function closeDetails() {
		showDetailsModal = false;
		selectedLink = null;
	}

	async function createLink(e: CustomEvent<CreateLinkPayload>) {
		try {
			await api.createLink(e.detail);
			showCreateModal = false;
			await loadLinks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create link';
		}
	}

	async function editLink(e: CustomEvent<UpdateLinkPayload>) {
		if (!selectedLink) return;
		try {
			await api.updateLink(selectedLink.id, e.detail);
			showEditModal = false;
			showDetailsModal = false;
			selectedLink = null;
			await loadLinks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to update link';
		}
	}

	async function deleteLink() {
		if (!selectedLink) return;
		try {
			await api.deleteLink(selectedLink.id);
			showDeleteModal = false;
			showDetailsModal = false;
			selectedLink = null;
			await loadLinks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete link';
		}
	}

	async function doLogout() {
		logout();
		await goto('/auth');
	}

	async function goToRegister() {
		await goto('/register');
	}
</script>

<main class="min-h-screen px-4 py-8 max-w-5xl mx-auto">
	<header class="flex items-start justify-between gap-3 mb-6">
		<div class="min-w-0">
			<h1 class="text-3xl font-bold text-cyan-300">[ dashboard ]</h1>
			<p class="text-gray-400 text-sm">read://links</p>
		</div>

		<!-- Desktop actions -->
		<div class="hidden sm:flex items-center gap-2">
			<button class="btn-primary" onclick={openCreate}>+ Add link</button>
			<AuthMenu on:logout={doLogout} on:register={goToRegister} />
		</div>

		<!-- Mobile actions -->
		<div class="sm:hidden flex flex-col items-end gap-2">
			<button
				class="btn-primary !px-3 !py-2 font-bold leading-none"
				onclick={openCreate}
				aria-label="Add link"
				title="Add link"
			>
				+
			</button>
			<AuthMenu on:logout={doLogout} on:register={goToRegister} />
		</div>
	</header>

	<!-- Mobile sticky add button -->
	<button
		class="sm:hidden fixed bottom-5 right-4 z-40 btn-primary !px-4 !py-3 shadow-lg"
		onclick={openCreate}
		aria-label="Add link"
		title="Add link"
	>
		+
	</button>

	{#if loading}
		<p class="text-gray-400">Loading links...</p>
	{:else if error}
		<p class="text-red-400">{error}</p>
	{:else if links.length === 0}
		<p class="text-gray-400">No links yet.</p>
	{:else}
		<!-- Compact favicon grid -->
		<div class="grid grid-cols-5 sm:grid-cols-7 lg:grid-cols-10 gap-1.5">
        {#each links as link (link.id)}
				<LinkCompactCard
					{link}
					on:open={(e) => openLink(e.detail)}
					on:details={(e) => openDetails(e.detail)}
					on:edit={(e) => openEdit(e.detail)}
					on:delete={(e) => openDelete(e.detail)}
				/>
			{/each}
		</div>
	{/if}
</main>

{#if showCreateModal}
	<LinkModal mode="create" on:close={() => (showCreateModal = false)} on:save={createLink} />
{/if}

{#if showEditModal && selectedLink}
	<LinkModal
		mode="edit"
		initial={selectedLink}
		on:close={() => {
			showEditModal = false;
			selectedLink = null;
		}}
		on:save={editLink}
	/>
{/if}

{#if showDeleteModal && selectedLink}
	<ConfirmDialog
		title="Delete link"
		message={`Delete "${selectedLink.name}"? This action cannot be undone.`}
		on:close={() => {
			showDeleteModal = false;
			selectedLink = null;
		}}
		on:confirm={deleteLink}
	/>
{/if}

{#if showDetailsModal && selectedLink}
	<LinkDetailsDialog
		link={selectedLink}
		on:close={closeDetails}
		on:open={(e) => openLink(e.detail)}
		on:edit={(e) => {
			showDetailsModal = false;
			openEdit(e.detail);
		}}
		on:delete={(e) => {
			showDetailsModal = false;
			openDelete(e.detail);
		}}
	/>
{/if}
