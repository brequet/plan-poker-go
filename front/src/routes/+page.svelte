<script>
	let roomName = '';
	let roomCode = '';

	// Function to create a new room (HTTP POST request here)
	async function createRoom() {
		// Implement your logic here to create a room
		// Send a POST request to your Golang backend to create the room
		// Use the 'roomName' variable to get the entered room name
		// Handle the response from the server accordingly
		console.log('Creating room:', roomName);

		try {
			const response = await fetch('/room', {
				method: 'POST',
				body: JSON.stringify({ roomName }),
				headers: {
					'Content-Type': 'application/json'
				}
			});
			// const response = await fetch('http://127.0.0.1/api/room', {
			// 	method: 'POST',
			// 	body: JSON.stringify({ roomName }),
			// 	headers: {
			// 		'Content-Type': 'application/json'
			// 	}
			// });

			if (!response.ok) {
				console.log('response not ok');
			} else {
				// The room was created successfully, you can handle the response here
				const responseData = await response.json();
				console.log('Room ID:', responseData.roomId);
				alert('ROOM CODE : ' + responseData.roomId);
				// TODO:change all room id to room code
                // todo: refactor
			}
		} catch (error) {
			console.error('Error creating room:', error);
		}
	}

	// Function to join an existing room (WebSocket connection here)
	function joinRoom() {
		// Implement your logic here to join a room
		// Use the 'roomName' variable to get the entered room name
		// Establish a WebSocket connection to the Golang backend with the room name
		console.log('Joining room:', roomName);
	}
</script>

<h1 class="text-3xl font-bold underline">
	Plan Poker Go
</h1>

<div class="flex justify-center items-center h-screen">
	<!-- Create Room Block -->
	<div class="bg-white p-6 rounded-lg shadow-lg mx-4">
	  <h2 class="text-2xl font-bold mb-4">Create a Room</h2>
	  <div class="flex items-center mb-4">
		<label class="mr-2">Room Name:</label>
		<input
		  class="border rounded-md py-2 px-4 w-full"
		  type="text"
		  bind:value={roomName}
		  placeholder="Enter Room Name"
		/>
	  </div>
	  <button
		class="bg-blue-500 text-white font-semibold px-4 py-2 rounded-md"
		on:click={createRoom}
	  >
		Create
	  </button>
	</div>
  
	<!-- Join Room Block -->
	<div class="bg-white p-6 rounded-lg shadow-lg mx-4">
	  <h2 class="text-2xl font-bold mb-4">Join a Room</h2>
	  <div class="flex items-center mb-4">
		<label class="mr-2">Room Code:</label>
		<input
		  class="border rounded-md py-2 px-4 w-full"
		  type="text"
		  bind:value={roomCode}
		  placeholder="Enter Room Code"
		/>
	  </div>
	  <button
		class="bg-blue-500 text-white font-semibold px-4 py-2 rounded-md"
		on:click={joinRoom}
	  >
		Join
	  </button>
	</div>
  </div>

<style>
	.container {
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
	}

	.block {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 20px;
		border: 1px solid #ccc;
		border-radius: 5px;
		margin: 10px;
		max-width: 400px;
	}

	.input-group {
		display: flex;
		flex-direction: row;
		align-items: center;
		margin: 10px;
	}

	.input-label {
		min-width: 100px;
		padding-right: 10px;
		font-weight: bold;
	}

	.input-field {
		flex-grow: 1;
		padding: 5px;
		border: 1px solid #ccc;
		border-radius: 3px;
	}

	.button {
		padding: 10px 20px;
		background-color: #007bff;
		color: #fff;
		border: none;
		border-radius: 3px;
		cursor: pointer;
	}

	.button:disabled {
		opacity: 0.3;
	}

	.button:disabled:hover {
		cursor: not-allowed;
	}
</style>
