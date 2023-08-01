<script lang="ts">
	import { page } from '$app/stores';
	import { onDestroy } from 'svelte';
	import NicknameChoice from './NicknameChoice.svelte';
	import PlanningPokerRoom from './PlanningPokerRoom.svelte';
	import RoomNotFound from './RoomNotFound.svelte';
	import WebSocketConnection from './WebSocketConnection.svelte';
	import {
		MessageType,
		type Message,
		type UserJoinedMessage,
		type UserDisconnectedMessage,
		type ConfirmConnectionMessage,
	} from './message';
	import { webSocketConnection } from './webSocketStore';

	export let data;

	let isConnectionConfirmedByUser = false;

	type User = {
		// TODO: refactor out
		nickname: string;
		uuid?: string;
		estimate?: number;
	};

	let roomData: {
		// TODO: in store
		user: User;
		room: {
			code: string;
			name?: string;
			exist: boolean;
		};
		users: User[];
	} = {
		user: {
			// TODO: type, with id..
			nickname: data.nickname ?? '',
			uuid: undefined,
			estimate: undefined
		},
		room: {
			code: $page.params.room,
			name: data.room?.name,
			exist: data.room !== undefined // TODO: fetch and update room status at page arrival
		},
		users: []
	};

	let socket: WebSocket; //TODO: do something if socket becomes null (ex: server crashes)
	const unsubscribeFromSocketWritable = webSocketConnection.subscribe((value) => (socket = value));

	async function onNicknameChoice(nickname: string) {
		roomData.user.nickname = nickname;
		localStorage.setItem('nickname', nickname);
		// TODO: check with server if room exist
		const response = await fetch(`http://127.0.0.1:8080/api/room/${roomData.room.code}`, {
			method: 'GET'
		});

		if (!response.ok) {
			console.log('response not ok');
		} else {
			const { roomCode } = await response.json();
			console.log('found room onNicknameChoice', roomCode);
		}

		isConnectionConfirmedByUser = true;
	}

	// TODO: validate route room name https://learn.svelte.dev/tutorial/param-matchers -> simple regex like [AZ]{4} -> 4 from env file / conf / properties

	function handleWsMessage(message: Message) {
		switch (message.type) {
			case MessageType.CONFIRM_CONNECTION:
				console.log('You are succesfuly connected');
				const confirmConnectionMessage: ConfirmConnectionMessage =
					message as ConfirmConnectionMessage;
				roomData.users = [
					...confirmConnectionMessage.payload.ConnectedUsers.map((user) => {
						// TODO: clean + clean message.ts also (do like in message.go, send, receive, etc..)
						const joiningUser: User = {
							nickname: user.userName,
							uuid: user.uuid,
							estimate: user.estimate
						};
						return joiningUser;
					})
				];
				break;
			case MessageType.USER_DISCONNECTED:
				console.log('pelo disconnected');
				const userDisconnectedMessage: UserDisconnectedMessage = message as UserDisconnectedMessage;
				roomData.users = roomData.users.filter(
					(user) => user.uuid !== userDisconnectedMessage.payload.user.uuid
				);
				break;
			case MessageType.USER_JOINED:
				console.log('pelo joined');
				const userJoinedMessage: UserJoinedMessage = message as UserJoinedMessage;
				const joiningUser: User = {
					nickname: userJoinedMessage.payload.user.userName,
					uuid: userJoinedMessage.payload.user.uuid,
					estimate: userJoinedMessage.payload.user.estimate
				};
				roomData.users = [...roomData.users, joiningUser];
				break;
			case MessageType.CONFIRM_ESTIMATE_SUBMISSION:
				console.log('You submitted succesfully !', message);
				break;
			case MessageType.SUBMIT_ESTIMATE:
				break;
			case MessageType.ESTIMATE_REVEALED:
				break;
			case MessageType.RESET_PLANNING:
				break;
			default:
				console.log('default in switch :/', message);

				break;
		}
	}

	// TODO: page title 'Poker Room ABCD'

	onDestroy(() => unsubscribeFromSocketWritable());
</script>

<div class="container mx-auto">
	{#if !roomData.room.exist}
		<RoomNotFound roomCode={roomData.room.code} />
	{:else if !isConnectionConfirmedByUser || roomData.user.nickname === ''}
		<NicknameChoice
			nickname={roomData.user.nickname}
			on:nicknameChoice={(event) => onNicknameChoice(event.detail.nickname)}
		/>
	{:else}
		<WebSocketConnection {...roomData} on:message={(event) => handleWsMessage(event.detail)} />
		{#if socket !== null}
			<PlanningPokerRoom {...roomData} />
		{/if}
	{/if}
</div>
