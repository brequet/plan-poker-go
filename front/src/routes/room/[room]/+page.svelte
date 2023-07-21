<script lang="ts">
	import WebSocketConnection from './WebSocketConnection.svelte';


    export let data;

    $: room = data.room;
</script>

Room <span style="font-weight: bold;">{data.room}</span>

<WebSocketConnection room={data.room}/>
