import { getActiveAccount, refreshActiveTokens, removeAccount } from '../stores/session';

const BASE_URL = '/api';

async function request<T>(path: string, options: RequestInit = {}, isRetry = false): Promise<T> {
	const account = getActiveAccount();

	const response = await fetch(`${BASE_URL}${path}`, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...(account ? { Authorization: `Bearer ${account.accessToken}` } : {}),
			...options.headers
		}
	});

	if (response.status === 401 && !isRetry) {
		try {
			await refreshActiveTokens();
			return request<T>(path, options, true);
		} catch {
			const current = getActiveAccount();
			if (current) removeAccount(current.id);
			throw new Error('Session expired');
		}
	}

	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: { message: 'Unknown error' } }));
		throw new Error(error?.error?.message ?? error?.message ?? response.statusText);
	}

	if (response.status === 204) {
		return undefined as T;
	}

	return response.json();
}

export const api = {
	get: <T>(path: string) => request<T>(path),
	post: <T>(path: string, body: unknown) =>
		request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
	patch: <T>(path: string, body: unknown) =>
		request<T>(path, { method: 'PATCH', body: JSON.stringify(body) }),
	delete: <T>(path: string) => request<T>(path, { method: 'DELETE' })
};
