import type { ZodiosPlugin } from '@zodios/core';
import { createApiClient } from './generated';
import { getActiveAccount, refreshActiveTokens, removeAccount } from '../stores/session';

const authPlugin: ZodiosPlugin = {
	request: async (_api, config) => {
		const account = getActiveAccount();
		if (!account) return config;
		return {
			...config,
			headers: {
				...config.headers,
				Authorization: `Bearer ${account.accessToken}`
			}
		};
	},
	error: async (_api, _config, error) => {
		// @ts-expect-error axios error shape
		if (error?.response?.status === 401) {
			const account = getActiveAccount();
			if (!account) throw error;
			try {
				await refreshActiveTokens();
			} catch {
				removeAccount(account.id);
			}
		}
		throw error;
	}
};

export const api = createApiClient('/api');
api.use(authPlugin);
