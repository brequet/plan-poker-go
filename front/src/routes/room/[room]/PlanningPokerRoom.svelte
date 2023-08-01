<script lang="ts">
	import { onDestroy } from 'svelte';
	import { webSocketConnection } from './webSocketStore';
	import { MessageType, type SubmitEstimateMessage } from './message';

	export let user: User;
	export let room: {
		code: string;
		name?: string;
		exist: boolean;
	};
	type User = {
		// TODO: refactor out
		nickname: string;
		uuid?: string;
		estimate?: number;
	};

	export let users: User[];

	const votingOptions = ['1', '2', '3', '5', '8', '13', '20', '40', '?', '☕'];

	let socket: WebSocket;
	const unsubscribeFromSocketWritable = webSocketConnection.subscribe((value) => (socket = value));

	$: connectedUsers = [
		{
			nickname: user.nickname,
			uuid: null
		},
		...users
	];

	function submitEstimate(estimate: string) {
		const submitEstimateMessage: SubmitEstimateMessage = {
			type: MessageType.SUBMIT_ESTIMATE,
			payload: {
				estimate
			}
		};
		console.log('Submitting estimate', submitEstimateMessage)
		socket.send(JSON.stringify(submitEstimateMessage));
	}

	onDestroy(() => unsubscribeFromSocketWritable());
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
			<button
				class="bg-blue-500 hover:bg-blue-600 text-white text-center py-2 rounded-lg cursor-pointer"
				on:click={() => submitEstimate(votingOption)}
			>
				{votingOption}
			</button>
		{/each}
	</div>
</div>
