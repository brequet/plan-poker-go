import { json } from '@sveltejs/kit';

export async function POST({ request }): Promise<any> {
	const { roomName } = await request.json();
	let roomId: any;
	console.log('room name:', roomName);
	try {
		console.log('room name body post:', JSON.stringify({ roomName }));
		const response = await fetch('http://127.0.0.1:8080/api/room', {
			method: 'POST',
			body: JSON.stringify({ roomName }),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
            console.log('response not ok');
		} else {
			// The room was created successfully, you can handle the response here
			const responseData = await response.json();
            roomId = responseData.roomId;
			console.log('Room ID:', roomId);
		}
	} catch (error) {
		console.error('Error creating room:', error);
	}

    console.log('returning:',  {roomId} );
	return new Response(JSON.stringify({roomId}));
}
