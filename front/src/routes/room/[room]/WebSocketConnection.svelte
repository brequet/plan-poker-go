<script lang="ts">
	import { onMount } from 'svelte';

	export let room: string;

	let socket: WebSocket;

	$: connState = socket && socket.readyState;

	onMount(() => {
		const url = 'ws://localhost:8080/ws'; // Replace with the URL of your Golang backend WebSocket server
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log('WebSocket connected!');
            socket.send("salut")
		};

		socket.onmessage = (event) => {
			const message = event.data;
			console.log('Message received :', message);
			// Handle incoming messages from the server
			// (e.g., updating the user interface based on received data)
		};

		socket.onerror = (error) => {
			console.error('WebSocket error:', error);
		};

		socket.onclose = (event) => {
			console.log('WebSocket connection closed:', event.code, event.reason);
			console.log('->', event);
		};
	});
</script>

<p>Status : {connState}</p>
