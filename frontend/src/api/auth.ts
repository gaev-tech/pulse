import { api } from './client';

interface AuthTokens {
	access_token: string;
	refresh_token: string;
}

interface User {
	id: string;
	email: string;
	username: string;
}

interface VerifyResponse extends AuthTokens {
	user: User;
}

export const auth = {
	sendMagicLink: (email: string) =>
		api.post<void>('/v1/auth/magic-link', { email }),

	verifyMagicLink: (token: string) =>
		api.post<VerifyResponse>('/v1/auth/magic-link/verify', { token }),

	refresh: (refresh_token: string) =>
		api.post<AuthTokens>('/v1/auth/refresh', { refresh_token }),

	logout: (refresh_token: string) =>
		api.post<void>('/v1/auth/logout', { refresh_token }),

	me: () =>
		api.get<User>('/v1/auth/me')
};
