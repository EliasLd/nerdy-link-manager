<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { isAuthenticated, logout } from '$lib/auth';

	onMount(async () => {
		if (!isAuthenticated()) {
			await goto('/auth');
		}
	});

	async function doLogout() {
		logout();
		await goto('/auth');
	}
</script>

<main class="min-h-screen flex items-center justify-center px-4">
	<div class="w-full max-w-xl border border-cyan-500/30 bg-black/30 rounded-2xl p-8 text-center">
		<h1 class="text-3xl font-bold text-cyan-300">[ dashboard ]</h1>
		<p class="text-gray-400 mt-2">Logged in!</p>

		<div class="mt-6">
			<button class="btn-ghost" on:click={doLogout}>Logout</button>
		</div>
	</div>
</main>
