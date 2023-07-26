<script lang="ts">
	import { createEventDispatcher } from "svelte";

	let nickname = '';

    const dispatch = createEventDispatcher();

	function confirmNickname() {
		dispatch('nicknameChoice', {nickname})
	}
</script>

<div class="max-w-md mx-auto">
	<form class="bg-white p-6 rounded-lg shadow-md">
		<h2 class="text-2xl font-bold mb-4">Enter Your Nickname</h2>
		<div class="mb-4">
			<label for="nicknameInput" class="block text-gray-700 font-semibold">Nickname:</label>
			<input
				type="text"
				id="nicknameInput"
				class="border rounded-md py-2 px-4 w-full"
				bind:value={nickname}
				required
				placeholder="Enter your nickname"
			/>
		</div>
		<button
			class="bg-blue-500 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-md w-full"
			disabled={nickname.length === 0}
			on:click={confirmNickname}
		>
			Confirm
		</button>
	</form>
</div>

