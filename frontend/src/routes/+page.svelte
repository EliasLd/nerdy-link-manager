<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/auth';

	onMount(() => {
		goto(isAuthenticated() ? '/dashboard' : '/auth');
	});
</script>

<p class="text-gray-400 p-6">Redirecting...</p>
