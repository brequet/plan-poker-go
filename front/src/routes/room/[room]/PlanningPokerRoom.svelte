<script lang="ts">
	import { onDestroy } from 'svelte';
	import {
		MessageType,
		type ResetPlanningMessage,
		type RevealEstimateMessage,
		type SubmitEstimateMessage
	} from './message';
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

	$: selectedEstimate = currentUser.estimate;
	$: allUsers = [...connectedUsers, currentUser];
	$: allUsersVoted = allUsers.every((user) => hasVoted(user));
	$: countNumberOfVote = allUsers.reduce((count, user) => count + (hasVoted(user) ? 1 : 0), 0);

	function hasVoted(user: User): boolean {
		return user.estimate !== undefined && user.estimate.length > 0;
	}

	function submitEstimate(estimate: string) {
		if (estimate === selectedEstimate) {
			// second click on the same button => unselect
			estimate = '';
		}
		const submitEstimateMessage: SubmitEstimateMessage = {
			type: MessageType.SUBMIT_ESTIMATE,
			payload: {
				estimate
			}
		};
		selectedEstimate = estimate;
		console.log('Submitting estimate', submitEstimateMessage);
		socket.send(JSON.stringify(submitEstimateMessage));
	}

	function toggleVotedEstimate(): void {
		console.log('revealVotedEstimate');

		// TODO count down ! like 3 secondes

		const revealEstimateMessagge: RevealEstimateMessage = {
			type: MessageType.REVEAL_ESTIMATE,
			payload: {
				shouldReveal: !room.isEstimateRevealed
			}
		};

		socket.send(JSON.stringify(revealEstimateMessagge));
	}

	onDestroy(() => {
		unsubscribeFromSocketWritable();
		unsubscribeFromRoomStore();
		unsubscribeFromCurrentUserStore();
		unsubscribeFromConnectedUsersStore();
	});

	function resetPlanning(): void {
		console.log('resetPlanning');
		const resetPlanningMessage: ResetPlanningMessage = {
			type: MessageType.RESET_PLANNING
		};

		socket.send(JSON.stringify(resetPlanningMessage));
	}
</script>

<!-- Room Details -->
<div class="bg-white p-4 rounded-lg shadow mb-4">
	<h2 class="text-2xl font-bold mb-2">Room Details</h2>
	<p class="text-gray-600">Room Name: {room.name} [{room.code}]</p>
</div>

<!-- Users in the Room -->
<!-- TODO rework UI + add average etc. -->
<div class="bg-white p-4 rounded-lg shadow mb-4">
	<h2 class="text-2xl font-bold mb-2">Users in the Room</h2>
	<ul class="list-disc pl-6">
		<li>
			{currentUser.nickname} (me) :
			{#if room.isEstimateRevealed}
				{currentUser.estimate ? currentUser.estimate : "didn't voted"}
			{:else}
				{currentUser.estimate ? 'has voted' : 'is voting..'}
			{/if}
		</li>
		{#each connectedUsers as connectedUser}
			<li>
				{connectedUser.nickname} :
				{#if room.isEstimateRevealed}
					{connectedUser.estimate ? connectedUser.estimate : "didn't voted"}
				{:else}
					{connectedUser.estimate ? 'has voted' : 'is voting..'}
				{/if}
			</li>
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
				class="bg-blue-500 hover:bg-blue-600 text-white text-center py-2 rounded-lg cursor-pointer
				{votingOption === selectedEstimate ? 'bg-blue-700 hover:bg-blue-800' : ''}"
				on:click={() => submitEstimate(votingOption)}
			>
				{votingOption}
			</button>
		{/each}
	</div>

	{#if room.isEstimateRevealed}
		<button
			class="bg-slate-500 hover:bg-slate-600 text-white py-2 px-4 rounded-lg cursor-pointer mt-4 w-full"
			on:click={() => toggleVotedEstimate()}
		>
			Hide estimates
		</button>

		<button
			class="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded-lg cursor-pointer mt-4 w-full"
			on:click={() => resetPlanning()}
		>
			Reset estimates !
		</button>
	{:else}
		<button
			class="{allUsersVoted
				? 'bg-green-500 hover:bg-green-600'
				: 'bg-orange-500 hover:bg-orange-600'} 
		 text-white py-2 px-4 rounded-lg cursor-pointer mt-4 w-full"
			on:click={() => toggleVotedEstimate()}
		>
			Reveal Voted Estimates
			{allUsersVoted
				? ''
				: `(${allUsers.length - countNumberOfVote} user(s) didn't submit their estimate yet)`}
		</button>
	{/if}
</div>
