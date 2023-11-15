import { json } from '@sveltejs/kit';
import { GO_SERVER_ADDRESS } from '$env/static/private'


// TODO: move to api route

export async function POST({ request }): Promise<any> {
	const { roomName } = await request.json();
	let responseData: any; //TODO: type response here on room creation
	if (!roomName || roomName.length === 0) return json(responseData)
	console.log('room name:', roomName);
	try {
		console.log('room name body post:', JSON.stringify({ roomName }));
		console.log('Response here:', `${GO_SERVER_ADDRESS}/api/room`)
		const response = await fetch(`${GO_SERVER_ADDRESS}/api/room`, {
			method: 'POST',
			body: JSON.stringify({ roomName })
		});

		if (!response.ok) {
            console.log('response not ok');
		} else {
			responseData = await response.json();
		}
	} catch (error) {
		console.error('Error creating room:', error);
	}

    console.log('returning:',  responseData );
	return json(responseData);
}
