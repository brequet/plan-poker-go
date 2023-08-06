<script lang="ts">
	export let roomCode: string;
</script>

<div class="max-w-md mx-auto">
	<div class="bg-white p-6 rounded-lg shadow-md text-center">
		<h2 class="text-2xl font-bold mb-4">Room Not Found</h2>
		<p class="text-gray-700">This room does not exist anymore or may have been deleted.</p>
        <p class="text-gray-700">
        <a href="/room/{roomCode}" class="text-blue-500 no-underline hover:underline">Reload this page</a>
         or go back to the <a href="/" class="text-blue-500 no-underline hover:underline">Home page</a>.
        </p>
	</div>
</div>
