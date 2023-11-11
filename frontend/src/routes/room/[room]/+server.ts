import { json } from '@sveltejs/kit';
import { GO_SERVER_ADDRESS } from '$env/static/private'

// TODO: move to api route

export async function GET({ params }): Promise<any> {
	const roomCode = await params.room;
	console.log('In server.ts roomcode', roomCode)
	let responseData: any; //TODO: type response here on room creation
	try {
		const response = await fetch(`${GO_SERVER_ADDRESS}/api/room/${roomCode}`, {
			method: 'GET'
		});

		if (!response.ok) {
			console.log('response not ok');
		} else {
			// The room was created successfully, you can handle the response here
			responseData = await response.json();
		}
	} catch (error) {
		console.error('Error :', error);
	}

	return json(responseData);
}
