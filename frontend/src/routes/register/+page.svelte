<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { isAuthenticated } from '$lib/auth';
	import { onMount } from 'svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let success = $state('');
	let loading = $state(false);

	onMount(async () => {
		if (!isAuthenticated()) {
			await goto('/auth');
		}
	});

	async function submit() {
		error = '';
		success = '';
		loading = true;
		try {
			await api.register(email, password);
			success = 'Account created successfully.';
			email = '';
			password = '';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Register error';
		} finally {
			loading = false;
		}
	}
</script>

<main class="min-h-screen grid place-items-center px-4">
	<div class="w-full max-w-md border border-cyan-500/30 bg-black/30 rounded-2xl p-6">
		<h1 class="text-2xl font-bold text-cyan-300 mb-1">Create account</h1>
		<p class="text-sm text-gray-400 mb-5">[protected] requires-auth-session</p>

		<div class="space-y-3">
			<input class="w-full input-nerd" placeholder="email" bind:value={email} />
			<input class="w-full input-nerd" placeholder="password" type="password" bind:value={password} />
		</div>

		{#if error}<p class="text-red-400 text-sm mt-3">{error}</p>{/if}
		{#if success}<p class="text-green-400 text-sm mt-3">{success}</p>{/if}

		<button class="btn-primary w-full mt-4" onclick={submit} disabled={loading}>
			{loading ? '...' : 'register'}
		</button>

		<button class="btn-ghost w-full mt-3" onclick={() => goto('/dashboard')}>
			back to dashboard
		</button>
	</div>
</main>
