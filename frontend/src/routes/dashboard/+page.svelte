<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { isAuthenticated, logout } from '$lib/auth';
	import { api } from '$lib/api';
	import type { LinkItem, CreateLinkPayload, UpdateLinkPayload, FolderItem, CreateFolderPayload } from '$lib/types';
	import LinkCompactCard from '$lib/components/LinkCompactCard.svelte';
	import LinkDetailsDialog from '$lib/components/LinkDetailsDialog.svelte';
	import LinkModal from '$lib/components/LinkModal.svelte';
	import FolderModal from '$lib/components/FolderModal.svelte';
	import AddMenuDialog from '$lib/components/AddMenuDialog.svelte';
	import AddToFolderDialog from '$lib/components/AddToFolderDialog.svelte';
	import FolderCard from '$lib/components/FolderCard.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import AuthMenu from '$lib/components/AuthMenu.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import ChangePasswordModal from '$lib/components/ChangePasswordModal.svelte';
    import Clock from '$lib/components/Clock.svelte';

	let links = $state<LinkItem[]>([]);
    let allLinks = $state<LinkItem[]>([]);
	let folders = $state<FolderItem[]>([]);
	let loading = $state(true);
	let error = $state('');

	let showAddMenu = $state(false);
	let showCreateModal = $state(false);
	let showFolderModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let showDetailsModal = $state(false);
	let showAddToFolder = $state(false);
	let showChangePassword = $state(false);

	let selectedLink = $state<LinkItem | null>(null);
	let selectedFolder = $state<FolderItem | null>(null);

	let showRenameFolderModal = $state(false);
	let showDeleteFolderModal = $state(false);
	let folderActionTarget = $state<FolderItem | null>(null);

	let visibleLinks = $state<LinkItem[]>([]);

    onMount(async () => {
        if (!isAuthenticated()) {
            await goto('/auth');
            return;
        }
        await loadAll();

        const onKeydown = (e: KeyboardEvent) => {
            // Don't trigger shortcuts when typing in an input/textarea/contenteditable
            if (isEditableTarget(e.target)) return;

            // If another modal/dialog is open, ignore shortcuts (optional but safer)
            if (
                showCreateModal ||
                showFolderModal ||
                showEditModal ||
                showDetailsModal ||
                showDeleteModal ||
                showAddToFolder ||
                showRenameFolderModal ||
                showDeleteFolderModal ||
                showChangePassword
            ) {
                return;
            }

            // If AddMenu is open: handle L/F/Escape
            if (showAddMenu) {
                const k = e.key.toLowerCase();

                if (k === 'escape') {
                    e.preventDefault();
                    showAddMenu = false;
                    return;
                }
                if (k === 'l') {
                    e.preventDefault();
                    showAddMenu = false;
                    showCreateModal = true;
                    return;
                }
                if (k === 'f') {
                    e.preventDefault();
                    showAddMenu = false;
                    showFolderModal = true;
                    return;
                }
                return;
            }

            // Open AddMenu with "a"
            if (e.key.toLowerCase() === 'a') {
                e.preventDefault();
                showAddMenu = true;
            }
        };

        document.addEventListener('keydown', onKeydown);
        return () => document.removeEventListener('keydown', onKeydown);
    });

    async function loadAll(folderId?: string | null) {
        loading = true;
        error = '';

        try {
            const [folderLinks, allFolders, globalLinks] = await Promise.all([
                api.getLinks(true, folderId ?? undefined),
                api.getFolders(),
                api.getLinks(true)
            ]);

            folders = allFolders;
            links = folderLinks;
            allLinks = globalLinks;

            visibleLinks = folderId
                ? folderLinks
                : folderLinks.filter((l) => !l.folderId);

        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to load data';
        } finally {
            loading = false;
        }
    }

	async function openLink(link: LinkItem) {
		try {
			await api.registerClick(link.id);
			await loadAll(selectedFolder?.id ?? null);
		} catch {}
		window.open(link.url, '_blank', 'noopener,noreferrer');
	}

	function openAddMenu() { showAddMenu = true; }
	function openCreate() { showCreateModal = true; }
	function openFolderCreate() { showFolderModal = true; }

    function isEditableTarget(target: EventTarget | null) {
        const el = target as HTMLElement | null;
        return !!(
            el &&
            (el.tagName === 'INPUT' ||
                el.tagName === 'TEXTAREA' ||
                (el as HTMLElement).isContentEditable)
        );
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

	function openAddToFolder(link: LinkItem) {
		selectedLink = link;
		showAddToFolder = true;
	}

	function openFolder(folder: FolderItem) {
		selectedFolder = folder;
		loadAll(folder.id);
	}

	function clearFolderFilter() {
		selectedFolder = null;
		loadAll();
	}

	async function createLink(e: CustomEvent<CreateLinkPayload>) {
		try {
			await api.createLink(e.detail);
			showCreateModal = false;
			await loadAll(selectedFolder?.id ?? null);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create link';
            setTimeout(() => (error = ''), 4000);
		}
	}

	async function createFolder(e: CustomEvent<CreateFolderPayload>) {
		try {
			await api.createFolder(e.detail);
			showFolderModal = false;
			await loadAll(selectedFolder?.id ?? null);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create folder';
            setTimeout(() => (error = ''), 4000);
		}
	}

	async function editLink(e: CustomEvent<UpdateLinkPayload>) {
        if (!selectedLink) return;
        try {
            await api.updateLink(selectedLink.id, {
                ...e.detail,
                folderId: selectedLink.folderId ? Number(selectedLink.folderId) : null,
                customIcon: selectedLink.customIcon ?? null,
                faviconUrl: selectedLink.faviconUrl ?? null
            });
            showEditModal = false;
            showDetailsModal = false;
            selectedLink = null;
            await loadAll(selectedFolder?.id ?? null);
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to update link';
            setTimeout(() => (error = ''), 4000);
        }
    }

	async function assignFolder(folderId: string | null) {
        if (!selectedLink) return;
        try {
            await api.updateLink(selectedLink.id, {
                name: selectedLink.name,
                url: selectedLink.url,
                description: selectedLink.description ?? undefined,
                folderId: folderId ? Number(folderId) : null,
                customIcon: selectedLink.customIcon ?? null,
                faviconUrl: selectedLink.faviconUrl ?? null
            });
            showAddToFolder = false;
            selectedLink = null;
            await loadAll(selectedFolder?.id ?? null);
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to update link folder';
            setTimeout(() => (error = ''), 4000);
        }
    }

	async function deleteLink() {
		if (!selectedLink) return;
		try {
			await api.deleteLink(selectedLink.id);
			showDeleteModal = false;
			showDetailsModal = false;
			selectedLink = null;
			await loadAll(selectedFolder?.id ?? null);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete link';
            setTimeout(() => (error = ''), 4000);
		}
	}

	function openRenameFolder(folder: FolderItem) {
		folderActionTarget = folder;
		showRenameFolderModal = true;
	}

	function openDeleteFolder(folder: FolderItem) {
		folderActionTarget = folder;
		showDeleteFolderModal = true;
	}

	async function renameFolder(e: CustomEvent<CreateFolderPayload>) {
		if (!folderActionTarget) return;
		try {
			await api.updateFolder(folderActionTarget.id, { name: e.detail.name });
			showRenameFolderModal = false;
			folderActionTarget = null;
			await loadAll(selectedFolder?.id ?? null);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to rename folder';
            setTimeout(() => (error = ''), 4000);
		}
	}

	async function deleteFolder() {
		if (!folderActionTarget) return;
		try {
			await api.deleteFolder(folderActionTarget.id);

			// si on était dans ce folder, revenir au dashboard normal
			if (selectedFolder?.id === folderActionTarget.id) {
				selectedFolder = null;
			}

			showDeleteFolderModal = false;
			folderActionTarget = null;
			await loadAll(selectedFolder?.id ?? null);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete folder';
            setTimeout(() => (error = ''), 4000);
		}
	}

	async function removeFromFolder(link: LinkItem) {
		try {
			await api.updateLink(link.id, {
				name: link.name,
				url: link.url,
				description: link.description ?? undefined,
				folderId: null
			});
			await loadAll(selectedFolder?.id ?? null);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to remove from folder';
            setTimeout(() => (error = ''), 4000);
		}
	}

	async function doLogout() {
		logout();
		await goto('/auth');
	}

	async function goToRegister() {
		await goto('/register');
	}

	async function handleChangePassword(e: CustomEvent<{ current: string; next: string }>) {
		try {
			await api.changePassword(e.detail.current, e.detail.next);
			logout();
			await goto('/auth');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to change password';
            setTimeout(() => (error = ''), 4000);
		}
	}

    async function dropLinkOnFolder(e: CustomEvent<{ linkId: string; folder: FolderItem }>) {
        const { linkId, folder } = e.detail;
        const link = links.find((l) => l.id === linkId);
        if (!link) return;

        try {
            await api.updateLink(link.id, {
                name: link.name,
                url: link.url,
                description: link.description ?? undefined,
                folderId: Number(folder.id),
                customIcon: link.customIcon ?? null,
                faviconUrl: link.faviconUrl ?? null
            });            
            await loadAll(selectedFolder?.id ?? null);
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to move link to folder';
            setTimeout(() => (error = ''), 4000);
        }
    }

    let allLinksDragOver = $state(false);

    function onAllLinksDragOver(e: DragEvent) {
        e.preventDefault();
        allLinksDragOver = true;
        if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
    }

    function onAllLinksDragLeave() {
        allLinksDragOver = false;
    }

    async function onAllLinksDrop(e: DragEvent) {
        e.preventDefault();
        allLinksDragOver = false;

        const linkId = e.dataTransfer?.getData('text/plain');
        if (!linkId) return;

        const link = links.find((l) => l.id === linkId);
        if (!link) return;

        try {
            await api.updateLink(link.id, {
            name: link.name,
            url: link.url,
            description: link.description ?? undefined,
            folderId: null,
            customIcon: link.customIcon ?? null,
            faviconUrl: link.faviconUrl ?? null
        });
        await loadAll(selectedFolder?.id ?? null);
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to move link to all links';
            setTimeout(() => (error = ''), 4000);
        }
    }
</script>

<main class="min-h-screen px-4 py-8 max-w-5xl mx-auto">
    <header class="flex items-center justify-between gap-3 mb-6">
        <div class="min-w-0">
            <p class="text-gray-400 text-md md:text-xl font-bold">
                {selectedFolder ? `folder://${selectedFolder.name}` : 'read://links'}
            </p>
        </div>

        <div class="hidden md:flex flex-1 justify-center">
            <Clock />
        </div>

        <div class="flex items-center gap-2">
            <button class="btn-primary !px-3 !py-2 font-bold leading-none" onclick={openAddMenu}>+</button>
            <AuthMenu on:logout={doLogout} on:register={goToRegister} on:changePassword={() => (showChangePassword = true)} />
        </div>
    </header>

    <div class="md:hidden mb-3 text-center">
        <Clock />
    </div>

    {#if selectedFolder}
        <div class="mb-5 grid items-center gap-3 md:grid-cols-[auto_1fr_auto]">
            <button
                class={`btn-ghost transition ${allLinksDragOver ? 'border-cyan-300/80 bg-cyan-500/10 text-cyan-200' : ''}`}
                onclick={clearFolderFilter}
                ondragover={onAllLinksDragOver}
                ondragleave={onAllLinksDragLeave}
                ondrop={onAllLinksDrop}
            >
                ← All links
            </button>

            <div class="flex justify-center">
                <SearchBar
                    links={allLinks}
                    folders={folders}
                    on:select={(e) => openLink(e.detail)}
                />
            </div>

            <div class="hidden md:block min-w-[110px]"></div>
        </div>
    {:else}
        <div class="mb-5 flex justify-center">
            <SearchBar
                links={allLinks}
                folders={folders}
                on:select={(e) => openLink(e.detail)}
            />
        </div>
    {/if}

    {#if error}
        <div class="mb-4 rounded-lg border border-red-500/40 bg-red-500/10 text-red-300 px-4 py-3 flex items-start gap-3 w-sm">
            <div class="flex-1 text-sm">{error}</div>
            <button
                type="button"
                class="text-red-300 hover:text-red-200 transition"
                onclick={() => (error = '')}
                aria-label="Dismiss error"
            >
                ✕
            </button>
        </div>
    {/if}

	{#if loading}
		<p class="text-gray-400">Loading...</p>
	{:else}
		<div class="grid grid-cols-4 sm:grid-cols-6 lg:grid-cols-8 gap-2">
			{#if selectedFolder}
				{#each folders as folder (folder.id)}
					{#if folder.id === selectedFolder.id}
                        <FolderCard
                            folder={folder}
                            on:open={(e) => openFolder(e.detail)}
                            on:rename={(e) => openRenameFolder(e.detail)}
                            on:delete={(e) => openDeleteFolder(e.detail)}
                            on:dropLink={dropLinkOnFolder}
                        />
						{#each visibleLinks as link (link.id)}
							<LinkCompactCard
                                {link}
                                on:open={(e) => openLink(e.detail)}
                                on:details={(e) => openDetails(e.detail)}
                                on:edit={(e) => openEdit(e.detail)}
                                on:addToFolder={(e) => openAddToFolder(e.detail)}
                                on:delete={(e) => openDelete(e.detail)}
                                on:removeFromFolder={(e) => removeFromFolder(e.detail)}
                                on:iconError={(e) => (error = e.detail)}
                                on:iconUpdated={() => loadAll(selectedFolder?.id ?? null)}
                            />						
                        {/each}
					{:else}
						<FolderCard
							folder={folder}
							on:open={(e) => openFolder(e.detail)}
							on:rename={(e) => openRenameFolder(e.detail)}
							on:delete={(e) => openDeleteFolder(e.detail)}
                            on:dropLink={dropLinkOnFolder}
						/>
					{/if}
				{/each}
			{:else}
				{#each folders as folder (folder.id)}
					<FolderCard
						folder={folder}
						on:open={(e) => openFolder(e.detail)}
						on:rename={(e) => openRenameFolder(e.detail)}
						on:delete={(e) => openDeleteFolder(e.detail)}
                        on:dropLink={dropLinkOnFolder}
					/>
				{/each}

				{#each visibleLinks as link (link.id)}
					<LinkCompactCard
                        {link}
                        on:open={(e) => openLink(e.detail)}
                        on:details={(e) => openDetails(e.detail)}
                        on:edit={(e) => openEdit(e.detail)}
                        on:addToFolder={(e) => openAddToFolder(e.detail)}
                        on:delete={(e) => openDelete(e.detail)}
                        on:removeFromFolder={(e) => removeFromFolder(e.detail)}
                        on:iconError={(e) => (error = e.detail)}
                        on:iconUpdated={() => loadAll(selectedFolder?.id ?? null)}
                    />
				{/each}
			{/if}
		</div>
	{/if}
</main>

{#if showAddMenu}
	<AddMenuDialog
		on:close={() => (showAddMenu = false)}
		on:addLink={() => { showAddMenu = false; showCreateModal = true; }}
		on:addFolder={() => { showAddMenu = false; showFolderModal = true; }}
	/>
{/if}

{#if showCreateModal}
	<LinkModal mode="create" on:close={() => (showCreateModal = false)} on:save={createLink} />
{/if}

{#if showFolderModal}
	<FolderModal on:close={() => (showFolderModal = false)} on:save={createFolder} />
{/if}

{#if showEditModal && selectedLink}
	<LinkModal mode="edit" initial={selectedLink} on:close={() => { showEditModal = false; selectedLink = null; }} on:save={editLink} />
{/if}

{#if showDeleteModal && selectedLink}
	<ConfirmDialog
		title="Delete link"
		message={`Delete "${selectedLink.name}"? This action cannot be undone.`}
		on:close={() => { showDeleteModal = false; selectedLink = null; }}
		on:confirm={deleteLink}
	/>
{/if}

{#if showDetailsModal && selectedLink}
	<LinkDetailsDialog
		link={selectedLink}
		on:close={() => { showDetailsModal = false; selectedLink = null; }}
		on:open={(e) => openLink(e.detail)}
		on:edit={(e) => { showDetailsModal = false; openEdit(e.detail); }}
		on:delete={(e) => { showDetailsModal = false; openDelete(e.detail); }}
	/>
{/if}

{#if showRenameFolderModal && folderActionTarget}
	<FolderModal
		mode="edit"
		initial={folderActionTarget}
		on:close={() => { showRenameFolderModal = false; folderActionTarget = null; }}
		on:save={renameFolder}
	/>
{/if}

{#if showDeleteFolderModal && folderActionTarget}
	<ConfirmDialog
		title="Delete folder"
		message={`Delete "${folderActionTarget.name}"? Links inside will be moved back to the main dashboard.`}
		on:close={() => { showDeleteFolderModal = false; folderActionTarget = null; }}
		on:confirm={deleteFolder}
	/>
{/if}

{#if showAddToFolder && selectedLink}
	<AddToFolderDialog
		folders={folders}
		on:close={() => { showAddToFolder = false; selectedLink = null; }}
		on:choose={(e) => assignFolder(e.detail.folderId)}
	/>
{/if}
