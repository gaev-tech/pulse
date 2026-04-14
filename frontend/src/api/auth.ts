import { api } from './client';

export const auth = {
	sendMagicLink: (email: string) =>
		api.postV1authmagicLink({ email }),

	verifyMagicLink: (token: string) =>
		api.postV1authmagicLinkverify({ token }),

	refresh: (refresh_token: string) =>
		api.postV1authrefresh({ refresh_token }),

	logout: (refresh_token: string) =>
		api.postV1authlogout({ refresh_token }),

	me: () =>
		api.getV1authme()
};
