<script lang="ts">
	export let user: {nickname: string};
	export let room: {code: string};
	type User = { // TODO: refactor out
		nickname: string;
		uuid: string;
	};
	export let users: User[]

	const votingOptions = ['1', '2', '3', '5', '8', '13', '20', '40', '?', '☕'];

    $: connectedUsers = [{
		nickname: user.nickname,
		uuid: null
	}, ...users]
</script>

	<!-- Room Details -->
	<div class="bg-white p-4 rounded-lg shadow mb-4">
		<h2 class="text-2xl font-bold mb-2">Room Details</h2>
		<p class="text-gray-600">Room Name: {room.code}</p>
	</div>

	<!-- Users in the Room -->
	<div class="bg-white p-4 rounded-lg shadow mb-4">
		<h2 class="text-2xl font-bold mb-2">Users in the Room</h2>
		<ul class="list-disc pl-6">
			{#each connectedUsers as connectedUser}
                <li>{connectedUser.nickname}</li>
            {/each}
		</ul>
	</div>

	<!-- Poker Planning Interface -->
	<div class="bg-white p-4 rounded-lg shadow">
		<h2 class="text-2xl font-bold mb-2">Poker Planning</h2>
		<p class="text-gray-600">Select your estimate:</p>
		<div class="grid grid-cols-10 gap-4 mt-4">
			{#each votingOptions as votingOption}
				<div class="bg-blue-500 text-white text-center py-2 rounded-lg cursor-pointer">
					{votingOption}
				</div>
			{/each}
		</div>
	</div>
