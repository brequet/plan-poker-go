import type { Actions } from '@sveltejs/kit';
import { GO_SERVER_ADDRESS } from '$env/static/private'

export const load = async ({ cookies, params }) => {
	console.log('BEGIN room/[room]/+page.server.ts load')
	// TODO: fetch room info (code, name, exist ?)
	const roomCode = params.room;
	let room: {
		name: string;
		code: string;
	} | undefined;
	try {
		console.log('reqqquuueest url :', `${GO_SERVER_ADDRESS}/api/room/${roomCode}`);
		const response = await fetch(`${GO_SERVER_ADDRESS}/api/room/${roomCode}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
			console.log('response not ok load');
		} else {
			// The room was created successfully, you can handle the response here
			room = await response.json();
			console.log('room found', room);
		}
	} catch (error) {
		console.error('Error fetching room:', error);
	}

	console.log('END room/[room]/+page.server.ts load')
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
