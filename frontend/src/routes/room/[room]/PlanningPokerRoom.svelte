<script lang="ts">
	import 'chart.js/auto';
	import { CheckCircleIcon, CircleDashedIcon, Share2Icon } from 'lucide-svelte';
	import { minidenticon } from 'minidenticons';
	import { onDestroy } from 'svelte';
	import ShareLinkModal from './components/ShareLinkModal.svelte';
	import StatsView from './components/StatsView.svelte';
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
			// Put the current user first of the list
			connectedUsers = connectedUsersStore.sort((a, b) => {
				if (a.uuid === currentUser.uuid) return -1;
				if (b.uuid === currentUser.uuid) return 1;
				return 0;
			});
		}
	);

	const votingOptions = ['1', '2', '3', '5', '8', '13', '20', '40', '?', '☕'];

	let socket: WebSocket;
	const unsubscribeFromSocketWritable = webSocketConnection.subscribe((ws) => {
		if (ws) {
			socket = ws;
		}
	});

	let isShareLinkModalOpen = false;

	$: userTabWidth = room.isEstimateRevealed ? 50 : 100;

	$: selectedEstimate = currentUser.estimate;
	$: allUsersVoted = connectedUsers.every((user) => hasVoted(user));
	$: countNumberOfVote = connectedUsers.reduce((count, user) => count + (hasVoted(user) ? 1 : 0), 0);
	$: average = computeEstimateAverage(
		connectedUsers.map((user) => user.estimate).filter((estimate) => estimate !== undefined) as string[]
	);

	let estimatesGraphData: {
		labels: string[];
		datasets: {
			label: string;
			data: number[];
			backgroundColor: string[];
			borderColor: string[];
			borderWidth: number;
		}[];
	} = {
		labels: [],
		datasets: [
			{
				label: '# of Votes',
				data: [],
				backgroundColor: [],
				borderColor: [],
				borderWidth: 1
			}
		]
	};

	$: {
		let estimateCounts = connectedUsers.reduce((acc: { [key: string]: number }, user) => {
			if (user.estimate) {
				if (!acc[user.estimate]) {
					acc[user.estimate] = 1;
				} else {
					acc[user.estimate]++;
				}
			}
			return acc;
		}, {});

		estimatesGraphData.labels = [];
		estimatesGraphData.datasets[0].data = [];
		for (let estimate in estimateCounts) {
			estimatesGraphData.labels.push(estimate);
			estimatesGraphData.datasets[0].data.push(estimateCounts[estimate]);
			estimatesGraphData.datasets[0].backgroundColor.push('#3B94CB33');
			estimatesGraphData.datasets[0].borderColor.push('#3B94CB');
		}
	}

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
		console.debug('resetPlanning');
		const resetPlanningMessage: ResetPlanningMessage = {
			type: MessageType.RESET_PLANNING
		};

		socket.send(JSON.stringify(resetPlanningMessage));
	}

	function computeEstimateAverage(estimates: string[]): string {
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
			return '0';
		}

		return (totalEstimate / validEstimatesCount).toFixed(1);
	}

	onDestroy(() => {
		unsubscribeFromSocketWritable();
		unsubscribeFromRoomStore();
		unsubscribeFromCurrentUserStore();
		unsubscribeFromConnectedUsersStore();
	});
</script>

<ShareLinkModal bind:isShareLinkModalOpen />

<div class="bg-white p-4 rounded-lg shadow mb-4 flex flex-1 flex-col">
	<div class="flex">
		<h2 class="text-2xl font-bold mb-2">Room: {room.name}</h2>

		<button
			class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 ml-auto flex items-center"
			on:click={() => (isShareLinkModalOpen = true)}
		>
			<Share2Icon class="w-5 h-5 mr-2" />
			Share
		</button>
	</div>

	<div class="flex overflow-hidden flex-1">
		<!-- Users half -->
		<div style:width={userTabWidth + '%'} style:transition={'width 0.5s ease'}>
			<h3 class="text-xl font-bold mb-2">Users in the Room</h3>
			<ul class="list-disc pl-6">
				{#if connectedUsers.length === 0}
					No one else is here ! Feeling lonely ?
					<br />
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
								<p class="font-semibold">
									{connectedUser.nickname}
									{connectedUser.uuid === currentUser.uuid ? '(me)' : ''}
								</p>
								<p class="text-gray-500">
									{#if room.isEstimateRevealed}
										<!-- SHOULD SHOW USER ESTIMATE -->
										{#if connectedUser.estimate}
											<span class="font-bold text-neutral-700">{connectedUser.estimate}</span>
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

		<div class="ml-4" style:width={100 - userTabWidth + '%'} style:transition={'width 0.5s ease'}>
			<!-- Stats half -->
			{#if room.isEstimateRevealed}
				<StatsView bind:average bind:estimatesGraphData />
			{/if}
		</div>
	</div>
</div>

<!-- Poker Planning Interface -->
<div class="p-4">
	<p class="text-gray-600">Select your estimate:</p>
	<div class="grid grid-cols-5 md:grid-cols-10 gap-x-1 gap-y-2 md:gap-3 mt-2 sm:mt-4">
		{#each votingOptions as votingOption}
			<button
				class="border border-blue-500 text-center py-2 h-16 rounded-lg cursor-pointer
					{selectedEstimate === votingOption
					? 'bg-blue-500 text-white -translate-y-1 sm:-translate-y-2 hover:bg-blue-200 hover:text-black'
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
				? 'bg-blue-500 hover:bg-blue-600'
				: 'bg-orange-500 hover:bg-orange-600'}
		 		text-white py-2 px-4 rounded-lg cursor-pointer mt-4 w-full disabled:bg-gray-400 disabled:cursor-default"
			on:click={() => toggleVotedEstimate()}
			disabled={countNumberOfVote === 0}
		>
			Reveal Voted Estimates
			{allUsersVoted
				? ''
				: `(${connectedUsers.length - countNumberOfVote} user(s) didn't submit their estimate yet)`}
		</button>
	{/if}
</div>
