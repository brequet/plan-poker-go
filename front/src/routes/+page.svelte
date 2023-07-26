<script>
	import { goto } from '$app/navigation';

	let roomName = '';
	let roomCode = '';

	async function createRoom() {
		console.log('Creating room:', roomName);

		try {
			const response = await fetch('/room', {
				method: 'POST',
				body: JSON.stringify({ roomName }),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			if (!response.ok) {
				console.log('response not ok');
			} else {
				const { roomCode } = await response.json();
				goto(`/room/${roomCode}`);
				// TODO: refactor
			}
		} catch (error) {
			console.error('Error creating room:', error);
		}
	}

	function joinRoom() {
		console.log('Joining room:', roomCode);
		goto(`/room/${roomCode}`);
		// TODO: check room exist
	}

</script>

<div class="flex justify-center items-center h-full">
	<!-- Create Room Block -->
	<form class="bg-white p-6 rounded-lg shadow-lg mx-4"> 
		<!-- TODO: could be a component (create/join) -->
		<h2 class="text-2xl font-bold mb-4">Create a Room</h2>
		<div class="flex items-center mb-4">
			<label class="mr-2" for="roomNameInput">
				Room Name:
			</label>
				<input
					class="border rounded-md py-2 px-4 w-full"
					type="text"
					name="roomNameInput"
					bind:value={roomName}
					required
					placeholder="Enter Room Name"
				/>
		</div>
		<button
			class="bg-blue-500 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-md w-full"
			disabled={roomName.length === 0}
			on:click={createRoom}
		>
			Create
		</button>
	</form>

	<!-- Join Room Block -->
	<form class="bg-white p-6 rounded-lg shadow-lg mx-4">
		<h2 class="text-2xl font-bold mb-4">Join a Room</h2>
		<div class="flex items-center mb-4">
			<label class="mr-2" for="roomCodeInput">
				Room Code:
			</label>
				<input
					class="border rounded-md py-2 px-4 w-full"
					type="text"
					name="roomCodeInput"
					bind:value={roomCode}
					required
					placeholder="Enter Room Code"
				/>
		</div>
		<button
			class="bg-blue-500 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-md w-full"
			disabled={roomCode.length === 0}
			on:click={joinRoom}
		>
			Join
		</button>
	</form>
</div>
