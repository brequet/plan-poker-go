<script lang="ts">
	import { page } from '$app/stores';
	import NicknameChoice from './NicknameChoice.svelte';
	import PlanningPokerRoom from './PlanningPokerRoom.svelte';
	import RoomNotFound from './RoomNotFound.svelte';
	import WebSocketConnection from './WebSocketConnection.svelte';

	export let data;

	let isConnectionConfirmedByUser = false;

	let roomData = {
		user: {
			// TODO: type, with id..
			nickname: data.nickname ?? ''
		},
		room: {
			code: $page.params.room,
            name: 'todo',
			exist: !false // TODO: fetch and update room status at page arrival
		}
	};

	async function onNicknameChoice(nickname: string) {
		roomData.user.nickname = nickname;
        localStorage.setItem('nickname', nickname)
		// TODO: check with server if room exist
        const response = await fetch('/room', {
				method: 'POST',
				body: JSON.stringify( {roomName: roomData.room.name} ),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			if (!response.ok) {
				console.log('response not ok');
			} else {
				const { roomCode } = await response.json();
				console.log('eheeheh', roomCode);
			}

		isConnectionConfirmedByUser = true;
	}

    // TODO: validate route room name https://learn.svelte.dev/tutorial/param-matchers
</script>

<div class="container mx-auto">
	{#if !isConnectionConfirmedByUser}
		<NicknameChoice nickname={roomData.user.nickname} on:nicknameChoice={(event) => onNicknameChoice(event.detail.nickname)} />
	{:else if !roomData.room.exist}
		<RoomNotFound roomCode={roomData.room.code} />
	{:else}
		<WebSocketConnection {...roomData} />
		<PlanningPokerRoom {...roomData} />
	{/if}
</div>
