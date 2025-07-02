<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<main class="w-full min-h-screen bg-blue-900">
	{@render children()}
</main>
