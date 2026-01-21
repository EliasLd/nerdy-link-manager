<script lang="ts">
    import { onMount } from "svelte";
    import { PUBLIC_API_URL } from "$env/static/public";

    let healthStatus: string = "Loadind...";

    onMount(async () => {
        try {
            const res = await fetch(`${PUBLIC_API_URL}/health`);
            healthStatus = await res.text();
        } catch (err) {
            healthStatus = "Server unreachable";
        }
    });
</script>

<main class="min-h-screen flex items-center justify-center bg-gray-950 text-gray-100">
  <div class="text-center space-y-4">
    <h1 class="text-4xl font-bold">Nerdy Link Manager</h1>
    <p class="text-lg text-gray-400">Backend status:</p>
    <p class="text-2xl font-mono">{healthStatus}</p>
  </div>
</main>
