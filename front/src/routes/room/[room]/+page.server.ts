import type { Actions } from '@sveltejs/kit';

export const load = ({ cookies }) => {
	const nickname = cookies.get('nickname');

    console.log('fetch cookie nickname', nickname)
	// todo: fetch room info (code, name, exist ?)

	return {
		nickname
	};
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
        console.log('cookie nickname actio,', request)
		const formData = await request.formData();
		const nickname = formData.get('nickname') as string;
		if (nickname !== '') {
            console.log('cookie nickname', nickname)
			cookies.set('nickname', nickname);
		}
	}
};
