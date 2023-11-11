<script lang="ts">
	import { onMount } from 'svelte';

	export let isOpen: boolean;
	export let onClose: () => void;

	function handleBackgroundClick(event: any) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	onMount(() => {
		window.addEventListener('keydown', (event) => {
			if (event.key === 'Escape') {
				onClose();
			}
		});
	});
</script>

{#if isOpen}
	<div
		class="fixed top-0 left-0 w-full h-full flex items-center justify-center bg-opacity-50 bg-gray-900 z-50"
		on:click={handleBackgroundClick}
	>
		<div class="bg-white p-4 rounded-lg shadow-lg">
			<slot />
			<button
				class="mt-4 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
				on:click={onClose}
			>
				Close
			</button>
		</div>
	</div>
{/if}
