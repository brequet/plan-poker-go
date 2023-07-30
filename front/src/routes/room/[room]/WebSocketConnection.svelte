<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import { MessageType } from './message';
	import type { JoinRoomMessage, Message } from './message';

	export let user: { nickname: string };
	export let room: { code: string };

	const dispatch = createEventDispatcher();

	let socket: WebSocket;

	onMount(() => {
		const url = 'ws://localhost:8080/ws'; // Replace with the URL of your Golang backend WebSocket server
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log('WebSocket connected!');

			const joinRoomMessage: JoinRoomMessage = {
				type: MessageType.JOIN_ROOM,
				payload: {
					roomCode: room.code, // Replace with the actual room code
					nickname: user.nickname // Replace with the actual user nickname
				}
			};

			// Send the join_room message to the backend
			socket.send(JSON.stringify(joinRoomMessage));

			// socket.send(`hello from ${user.nickname}`);

		};

		socket.onmessage = (event) => {
			const message: Message = event.data;
			console.log('Message received :', message);
			// Handle incoming messages from the server
			// (e.g., updating the user interface based on received data)
			dispatch(message.type, message.payload);
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
