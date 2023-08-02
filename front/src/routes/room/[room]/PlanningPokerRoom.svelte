<script lang="ts">
	import { onDestroy } from 'svelte';
	import { MessageType, type SubmitEstimateMessage } from './message';
	import {
		connectedUsersStore,
		currentUserStore,
		roomStore,
		type CurrentUser,
		type Room,
		type User
	} from './room';
	import { webSocketConnection } from './webSocketStore';

	let room: Room;
	const unsubscribeFromRoomStore = roomStore.subscribe((roomStore) => {
		room = roomStore;
	});

	let currentUser: CurrentUser;
	const unsubscribeFromCurrentUserStore = currentUserStore.subscribe((currentUserStore) => {
		currentUser = currentUserStore;
	});

	let connectedUsers: User[];
	const unsubscribeFromConnectedUsersStore = connectedUsersStore.subscribe(
		(connectedUsersStore) => {
			connectedUsers = connectedUsersStore;
		}
	);

	const votingOptions = ['1', '2', '3', '5', '8', '13', '20', '40', '?', '☕'];

	let socket: WebSocket;
	const unsubscribeFromSocketWritable = webSocketConnection.subscribe((ws) => {
		if (ws) {
			socket = ws;
		}
	});

	let selected: string;

	function submitEstimate(estimate: string) {
		if (estimate === selected) {
			// second click on same button => unselect
			estimate = ""
		}
		const submitEstimateMessage: SubmitEstimateMessage = {
			type: MessageType.SUBMIT_ESTIMATE,
			payload: {
				estimate
			}
		};
		selected = estimate;
		console.log('Submitting estimate', submitEstimateMessage);
		socket.send(JSON.stringify(submitEstimateMessage));
	}

	onDestroy(() => {
		unsubscribeFromSocketWritable();
		unsubscribeFromRoomStore();
		unsubscribeFromCurrentUserStore();
		unsubscribeFromConnectedUsersStore();
	});
</script>

<!-- Room Details -->
<div class="bg-white p-4 rounded-lg shadow mb-4">
	<h2 class="text-2xl font-bold mb-2">Room Details</h2>
	<p class="text-gray-600">Room Name: {room.name} [{room.code}]</p>
</div>

<!-- Users in the Room -->
<div class="bg-white p-4 rounded-lg shadow mb-4">
	<h2 class="text-2xl font-bold mb-2">Users in the Room</h2>
	<ul class="list-disc pl-6">
		<li>{currentUser.nickname} (me) {currentUser.estimate ? 'has voted' : 'is voting..'}</li>
		{#each connectedUsers as connectedUser}
			<li>{connectedUser.nickname} {connectedUser.estimate ? 'has voted' : 'is voting..'}</li>
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
				class="bg-blue-500 hover:bg-blue-600 text-white text-center py-2 rounded-lg cursor-pointer {votingOption ===
				selected
					? 'bg-blue-700 hover:bg-blue-800'
					: ''}"
				on:click={() => submitEstimate(votingOption)}
			>
				{votingOption}
			</button>
		{/each}
	</div>
</div>
