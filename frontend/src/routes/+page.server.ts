import { GO_SERVER_ADDRESS } from '$env/static/private';

export async function load({ fetch }) {
	let isServerOk = false;
	try {
		const response = await fetch(`${GO_SERVER_ADDRESS}/api/health`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		isServerOk = response.ok;
	} catch (error) {
		console.log('Something went wrong when checking the health of the server');
	}
	return { isServerOk };
}
