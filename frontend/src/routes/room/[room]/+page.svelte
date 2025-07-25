<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { onDestroy } from 'svelte';
	import NicknameChoice from './NicknameChoice.svelte';
	import PlanningPokerRoom from './PlanningPokerRoom.svelte';
	import RoomNotFound from './RoomNotFound.svelte';
	import WebSocketConnection from './WebSocketConnection.svelte';
	import { MessageType, type Message } from './message';
	import {
		connectedUsersStore,
		currentUserStore,
		roomStore,
		type CurrentUser,
		type Room,
		type User
	} from './room';
	import { webSocketConnection } from './webSocketStore';
	import { invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';

	export let data;
	console.log('LOADED DATA', data);
	let webSocketUrl = '';
	onMount(() => {
		if (browser) {
			const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
			const host = window.location.host;
			webSocketUrl = `${protocol}://${host}/api/ws`;
		}

		initRoom(data);
	});

	function initRoom(data: {
		room: { name: string; code: string } | undefined;
		webSocketUrl: string;
	}) {
		console.log('LOADED DATA', data);
		roomStore.set({
			code: $page.params.room,
			name: data.room?.name,
			exist: data.room !== undefined,
			isEstimateRevealed: false
		});

		currentUserStore.set({
			nickname: browser ? localStorage.getItem('nickname') ?? '' : '',
			isConnected: false
		});
	}

	let socket: WebSocket | undefined;
	const unsubscribeFromSocketWritable = webSocketConnection.subscribe((ws) => {
		console.log(new Date(), 'ws', ws);
		if (socket !== undefined && ws === undefined) {
			console.log('Lost web socket connection');
			invalidateAll().then(() => initRoom(data));
		}

		if (ws) {
			socket = ws;
		}
	});

	let room: Room;
	const unsubscribeFromRoomStore = roomStore.subscribe((roomStore) => {
		room = roomStore;
	});

	let currentUser: CurrentUser;
	const unsubscribeFromCurrentUserStore = currentUserStore.subscribe((currentUserStore) => {
		currentUser = currentUserStore;
	});

	async function onNicknameChoice(nickname: string) {
		localStorage.setItem('nickname', nickname);

		currentUserStore.update((user) => {
			user.isConnected = true;
			user.nickname = nickname;
			return user;
		});
	}

	// TODO: validate route room name https://learn.svelte.dev/tutorial/param-matchers -> simple regex like [AZ]{4} -> 4 from env file / conf / properties

	function handleWsMessage(message: Message) {
		switch (message.type) {
			case MessageType.CONFIRM_CONNECTION:
				console.log(
					'You are succesfuly connected to room',
					message.payload.room.roomCode,
					message.payload.room.roomName
				);

				roomStore.update((room) => {
					return {
						...room,
						code: message.payload.room.roomCode,
						name: message.payload.room.roomName,
						isEstimateRevealed: message.payload.room.isEstimateRevealed
					};
				});

				currentUserStore.update((user) => {
					user.nickname = message.payload.user.userName;
					user.uuid = message.payload.user.uuid;
					return user;
				});

				connectedUsersStore.set([
					...message.payload.connectedUsers.map((user) => {
						const joiningUser: User = {
							nickname: user.userName,
							uuid: user.uuid,
							estimate: user.estimate
						};
						return joiningUser;
					})
				]);

				break;

			case MessageType.USER_DISCONNECTED:
				console.log('pelo disconnected');
				connectedUsersStore.update((connectedUsers) =>
					connectedUsers.filter((connectedUser) => connectedUser.uuid !== message.payload.user.uuid)
				);
				break;

			case MessageType.USER_JOINED:
				console.log('pelo joined');
				const joiningUser: User = {
					nickname: message.payload.user.userName,
					uuid: message.payload.user.uuid,
					estimate: message.payload.user.estimate
				};
				connectedUsersStore.update((connectedUsers) => [...connectedUsers, joiningUser]);
				break;

			case MessageType.CONFIRM_ESTIMATE_SUBMISSION:
				console.log('You submitted succesfully !', message);
				currentUserStore.update((currentUser) => {
					return { ...currentUser, estimate: message.payload.estimate };
				});
				break;

			case MessageType.USER_SUBMITTED_ESTIMATE:
				console.log('user submitted estimate !', message);
				connectedUsersStore.update((connectedUsers) => {
					connectedUsers.map((u) => {
						if (u.uuid === message.payload.user.uuid) {
							u.estimate = message.payload.estimate;
						}
					});
					return connectedUsers;
				});
				break;

			case MessageType.REVEAL_ESTIMATE: // TODO use ESTIMATE_REVEAL (1 direction per message type)
				console.log('Reveal estimate toggled !', message);
				roomStore.update((room) => {
					return { ...room, isEstimateRevealed: message.payload.shouldReveal };
				});
				break;

			case MessageType.PLANNING_RESETED:
				roomStore.update((room) => {
					return { ...room, isEstimateRevealed: false };
				});

				currentUserStore.update((user) => {
					return { ...user, estimate: undefined };
				});

				connectedUsersStore.update((users) => {
					return users.map((user) => {
						return { ...user, estimate: undefined };
					});
				});
				break;

			default:
				console.log('default in switch (message type not handled):', message);
				break;
		}
	}

	onDestroy(() => {
		unsubscribeFromSocketWritable();
		unsubscribeFromRoomStore();
		unsubscribeFromCurrentUserStore();
	});
</script>

<svelte:head>
	<title>{room.name}</title>
</svelte:head>

<div class="container mx-auto h-full flex flex-col">
	{#if !room?.exist}
		<RoomNotFound />
	{:else if !currentUser?.isConnected || currentUser?.nickname === ''}
		<NicknameChoice
			nickname={currentUser.nickname}
			on:nicknameChoice={(event) => onNicknameChoice(event.detail.nickname)}
		/>
	{:else}
		<WebSocketConnection
			{currentUser}
			roomCode={room.code}
			{webSocketUrl}
			on:message={(event) => handleWsMessage(event.detail)}
		/>
		{#if socket !== undefined && room.name}
			<!-- TODO: loading while socket connecting -->
			<PlanningPokerRoom />
		{/if}
	{/if}
</div>
