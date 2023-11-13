<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';
    
	export let isShareLinkModalOpen = false;

	const shareableLink = window.location.href;
</script>

<Modal isOpen={isShareLinkModalOpen} onClose={() => (isShareLinkModalOpen = false)}>
	<h2 class="text-xl font-semibold mb-2">Share Room Link</h2>
	<p>Copy and share this link:</p>
	<input
		type="text"
		value={shareableLink}
		class="w-full p-2 border rounded mt-2"
		readonly
	/>
	<button
		class="mt-2 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
		on:click={() => navigator.clipboard.writeText(shareableLink)}
	>
		Copy Link
	</button>
	<!-- TODO QR code + copyable url more clean + URL in .env and access using global store-->
</Modal>
