<script lang="ts">
	import WebSocketConnection from './WebSocketConnection.svelte';


    export let data;

    $: room = data.room;
</script>

<h1>Welcome in room : {data.room}</h1>

<WebSocketConnection room={data.room}/>
