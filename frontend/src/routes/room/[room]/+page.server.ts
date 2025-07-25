import { GO_SERVER_ADDRESS } from '$env/static/private';
import type { Actions } from '@sveltejs/kit';

export const load = async ({ params }) => {
	// TODO: fetch room info (code, name, exist ?)
	const roomCode = params.room;
	let room: {
		name: string;
		code: string;
	} | undefined;

	try {
		console.log(`User trying to join room at address: ${GO_SERVER_ADDRESS}/api/room/${roomCode}`);
		const response = await fetch(`${GO_SERVER_ADDRESS}/api/room/${roomCode}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
			throw new Error(`GET response not ok: ${response.status}`);
		} else {
			room = await response.json();
		}
	} catch (error) {
		console.error('Something went wrong while fetching the room.', error);
	}
	console.log("LOADED DATA BACK", room)
	return {
		room,
	};
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		console.log('room default action')
		const formData = await request.formData();
		const nickname = formData.get('nickname') as string;
		if (nickname !== '') {
			cookies.set('nickname', nickname);
		}
	}
};
