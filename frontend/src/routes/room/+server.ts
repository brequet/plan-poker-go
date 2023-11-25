import { json } from '@sveltejs/kit';
import { GO_SERVER_ADDRESS } from '$env/static/private';

// TODO: move to api route

export async function POST({ request }): Promise<any> {
	const { roomName } = await request.json();
	let responseData: any; //TODO: type response here on room creation
	if (!roomName || roomName.length === 0) return json(responseData);
	console.log(`Creating room with name '${roomName}'`);
	try {
		const response = await fetch(`${GO_SERVER_ADDRESS}/api/room`, {
			method: 'POST',
			body: JSON.stringify({ roomName })
		});

		if (!response.ok) {
			throw new Error(`POST response not ok: ${response.status}`);
		} else {
			responseData = await response.json();
		}
	} catch (error) {
		console.error('Something went wrong while creating the room.', error);
	}

	return json(responseData);
}
