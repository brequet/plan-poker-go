<script lang="ts">
	import { CheckCircleIcon, CircleDashedIcon } from 'lucide-svelte';
	import { minidenticon } from 'minidenticons';
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

	function resetPlanning(): void {
		console.log('resetPlanning');
		const resetPlanningMessage: ResetPlanningMessage = {
			type: MessageType.RESET_PLANNING
		};

		socket.send(JSON.stringify(resetPlanningMessage));
	}

	function computeEstimateAverage(estimates: string[]): number {
		let validEstimatesCount = 0;
		let totalEstimate = 0;

		for (const estimate of estimates) {
			// Try parsing the estimate as a number
			const numericEstimate = parseInt(estimate);

			// Check if the parsed value is a valid number
			if (!isNaN(numericEstimate)) {
				totalEstimate += numericEstimate;
				validEstimatesCount++;
			}
		}

		if (validEstimatesCount === 0) {
			return 0;
		}

		return totalEstimate / validEstimatesCount;
	}

	onDestroy(() => {
		unsubscribeFromSocketWritable();
		unsubscribeFromRoomStore();
		unsubscribeFromCurrentUserStore();
		unsubscribeFromConnectedUsersStore();
	});
</script>

<div class="bg-white p-4 rounded-lg shadow mb-4 flex-1">
	<h2 class="text-2xl font-bold mb-2">Room: {room.name}</h2>

	<h3 class="text-xl font-bold mb-2">Users in the Room</h3>
	<!-- TODO cut in half, left for users, right for stats (average..) -->
	<ul class="list-disc pl-6">
		{#if connectedUsers.length === 0}
			Feeling lonely TODO ?
		{:else}
			{#each connectedUsers as connectedUser}
				<li class="flex items-center py-2 border-b border-gray-200">
					<div class="w-12 h-12 flex-shrink-0">
						<img
							class="w-full h-full rounded-full object-cover"
							src={`data:image/svg+xml;utf8,${encodeURIComponent(
								minidenticon(connectedUser.nickname, undefined, undefined)
							)}`}
							alt={connectedUser.nickname}
						/>
					</div>
					<div class="ml-4">
						<p class="font-semibold">{connectedUser.nickname}</p>
						<p class="text-gray-500">
							{#if room.isEstimateRevealed}
								<!-- SHOULD SHOW USER ESTIMATE -->
								{#if connectedUser.estimate}
									<span>{connectedUser.estimate}</span>
								{:else}
									<span>Didn't vote yet...</span>
								{/if}
							{:else}
								<!-- SHOULD HIDE USER ESTIMATE -->
								{#if connectedUser.estimate}
									<span>Is ready !</span>
								{:else}
									<span>Is still voting...</span>
								{/if}
							{/if}
						</p>
					</div>
					{#if connectedUser.estimate}
						<CheckCircleIcon class="ml-auto" color="green" />
					{:else}
						<CircleDashedIcon class="ml-auto" color="orange" />
					{/if}
				</li>
			{/each}
		{/if}
	</ul>
</div>

<!-- Poker Planning Interface -->
<div class="p-4">
	<p class="text-gray-600">Select your estimate:</p>
	<div class="grid grid-cols-10 gap-4 mt-4">
		{#each votingOptions as votingOption}
			<button
				class="border border-blue-500 text-center py-2 h-16 rounded-lg cursor-pointer
					{votingOption === selectedEstimate
					? 'bg-blue-500 text-white -translate-y-2 hover:bg-blue-200 hover:text-black'
					: 'bg-white translate-y-0 hover:bg-blue-100'}
			  		hover:-translate-y-2
					transition-transform duration-300 transform"
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
			Reset estimates and estimate next issue
		</button>
	{:else}
		<button
			class="{allUsersVoted
				? 'bg-green-500 hover:bg-green-600'
				: 'bg-orange-500 hover:bg-orange-600'}
		 		text-white py-2 px-4 rounded-lg cursor-pointer mt-4 w-full disabled:bg-gray-400 disabled:cursor-default"
			on:click={() => toggleVotedEstimate()}
			disabled={countNumberOfVote === 0}
		>
			Reveal Voted Estimates
			{allUsersVoted
				? ''
				: `(${allUsers.length - countNumberOfVote} user(s) didn't submit their estimate yet)`}
		</button>
	{/if}
</div>
