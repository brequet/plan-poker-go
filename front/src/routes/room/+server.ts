import { json } from '@sveltejs/kit';

export async function POST({ request }): Promise<any> {
	const { roomName } = await request.json();
	let responseData: any; //TODO: type response on room creation
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
			responseData = await response.json();
		}
	} catch (error) {
		console.error('Error creating room:', error);
	}

    console.log('returning:',  responseData );
	return json(responseData);
}
