<script lang="ts">
	import { enhance } from "$app/forms";
	import { createEventDispatcher } from "svelte";

	export let nickname = '';

    const dispatch = createEventDispatcher();

	function confirmNickname() {
		dispatch('nicknameChoice', {nickname})
	}
</script>

<div class="max-w-md mx-auto">
	<div class="bg-white p-6 rounded-lg shadow-md" method="post"  >
		<h2 class="text-2xl font-bold mb-4">Enter Your Nickname</h2>
		<div class="mb-4">
			<label class="block text-gray-700 font-semibold">Nickname:
			<input
				type="text"
				class="border rounded-md py-2 px-4 w-full"
                name="nickname"
				bind:value={nickname}
				required
				placeholder="Enter your nickname"
			/>
        </label>
		</div>
		<button
			class="bg-blue-500 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-md w-full"
			disabled={nickname.length === 0}
			on:click={confirmNickname}
		>
			Confirm
		</button>
	</div>
</div>

