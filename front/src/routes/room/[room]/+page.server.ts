import type { Actions } from '@sveltejs/kit';

 export  const load = async ({ cookies, params }) => {
	const nickname = cookies.get('nickname');
	// TODO: fetch room info (code, name, exist ?)
	const roomCode = params.room;
	let roomExist = false;
	try {
		const response = await fetch(`http://127.0.0.1:8080/api/room/${roomCode}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
            console.log('response not ok');
		} else {
			// The room was created successfully, you can handle the response here
			const room = await response.json() ;
            console.log('room found', room);
			roomExist = room !== null
		}
	} catch (error) {
		console.error('Error fetching room:', error);
	}

	return {
		nickname,
		roomExist
	};
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const formData = await request.formData();
		const nickname = formData.get('nickname') as string;
		if (nickname !== '') {
			cookies.set('nickname', nickname);
		}
	}
};
