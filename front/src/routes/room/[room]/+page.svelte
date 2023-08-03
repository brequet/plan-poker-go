<script lang="ts">
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
	import Layout from '../../+layout.svelte';

	export let data;

	roomStore.set({
		code: $page.params.room,
		name: data.room?.name,
		exist: data.room !== undefined,
		isEstimateRevealed: false,
	});

	currentUserStore.set({
		nickname: data.nickname ?? '',
		isConnected: false
	});

	let socket: WebSocket; //TODO: do something if socket becomes null (ex: server crashes)
	const unsubscribeFromSocketWritable = webSocketConnection.subscribe((ws) => {
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
		const response = await fetch(`http://127.0.0.1:8080/api/room/${room.code}`, {
			method: 'GET'
		});

		if (!response.ok) {
			console.log('onNicknameChoice response not ok');
		} else {
			const roomInfo = await response.json();
			const { roomCode } = roomInfo;
			console.log('onNicknameChoice found room', roomInfo);
		}

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
				roomStore.update(room => {
					return {...room, isEstimateRevealed: message.payload.shouldReveal}
				})
				break;

			case MessageType.PLANNING_RESETED:
				roomStore.update(room => {
					return {...room, isEstimateRevealed: false}
				});

				currentUserStore.update(user => {
					return {...user, estimate: undefined}
				})

				connectedUsersStore.update(users => {
					return users.map(user => {
						return {...user, estimate: undefined}
					});
				})
				break;

			default:
				console.log('default in switch (message type not handled):', message);
				break;
		}
	}

	// TODO: page title 'Poker Room ABCD'

	onDestroy(() => {
		unsubscribeFromSocketWritable();
		unsubscribeFromRoomStore();
		unsubscribeFromCurrentUserStore();
	});
</script>

<div class="container mx-auto">
	{#if !room?.exist}
		<RoomNotFound roomCode={room.code} />
	{:else if !currentUser?.isConnected || currentUser?.nickname === ''}
		<NicknameChoice 
		nickname={currentUser.nickname}
		on:nicknameChoice={(event) => onNicknameChoice(event.detail.nickname)} 
		/>
	{:else}
		<WebSocketConnection
			{currentUser}
			roomCode={room.code}
			on:message={(event) => handleWsMessage(event.detail)}
		/>
		{#if socket !== null && room.name} 
		<!-- TODO: loading while socket connecting -->
			<PlanningPokerRoom />
		{/if}
	{/if}
</div>
