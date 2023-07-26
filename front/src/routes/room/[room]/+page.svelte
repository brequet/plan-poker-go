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
			// todo type, with id..
			nickname: data.nickname ?? ''
		},
		room: {
			code: $page.params.room,
            name: 'todo',
			exist: false
		}
	};

	function onNicknameChoice(nickname: string) {
		roomData.user.nickname = nickname;

		// TODO: check with server if room exist

		isConnectionConfirmedByUser = true;
	}

    // TODO: valid route room name https://learn.svelte.dev/tutorial/param-matchers
</script>

<div class="container mx-auto">
	{#if !isConnectionConfirmedByUser}
		<NicknameChoice on:nicknameChoice={(event) => onNicknameChoice(event.detail.nickname)} />
	{:else if !roomData.room.exist}
		<RoomNotFound roomCode={roomData.room.code} />
	{:else}
		<WebSocketConnection {...roomData} />
		<PlanningPokerRoom {...roomData} />
	{/if}
</div>
