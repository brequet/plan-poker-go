<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import { webSocketConnection } from './webSocketStore';
	import { MessageType } from './message';
	import type { JoinRoomMessage, Message } from './message';

	export let user: { nickname: string };
	export let room: { code: string };

	const dispatch = createEventDispatcher();

	let socket: WebSocket;

	onMount(() => {
		const url = 'ws://localhost:8080/ws'; // TODO: from conf file || env file
		socket = new WebSocket(url);
		webSocketConnection.set(socket)

		socket.onopen = () => {
			console.log('WebSocket connected!');

			const joinRoomMessage: JoinRoomMessage = {
				type: MessageType.JOIN_ROOM,
				payload: {
					roomCode: room.code,
					nickname: user.nickname
				}
			};

			socket.send(JSON.stringify(joinRoomMessage));
		};

		socket.onmessage = (event) => {
			const message: Message = JSON.parse(event.data);
			dispatch('message', message);
		};

		socket.onerror = (error) => {
			console.error('WebSocket error:', error);
		};

		socket.onclose = (event) => {
			console.log('WebSocket connection closed:', event.code, event.reason);
			console.log('->', event);
		};
	});

	onDestroy(() => {
		if (socket !== null && socket !== undefined) {
			socket.close();
		}
	});
</script>
