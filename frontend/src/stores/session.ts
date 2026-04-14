import { writable, get } from 'svelte/store';

export interface Account {
	id: string;
	email: string;
	username: string;
	accessToken: string;
	refreshToken: string;
}

interface SessionState {
	accounts: Account[];
	activeIndex: number;
}

const STORAGE_KEY = 'pulse_session';

function loadFromStorage(): SessionState {
	if (typeof localStorage === 'undefined') return { accounts: [], activeIndex: 0 };
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (!raw) return { accounts: [], activeIndex: 0 };
		return JSON.parse(raw);
	} catch {
		return { accounts: [], activeIndex: 0 };
	}
}

function saveToStorage(state: SessionState) {
	if (typeof localStorage === 'undefined') return;
	localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

const { subscribe, set, update } = writable<SessionState>(loadFromStorage());

export const session = { subscribe };

export function getActiveAccount(): Account | null {
	const state = get({ subscribe });
	if (!state.accounts.length) return null;
	return state.accounts[state.activeIndex] ?? null;
}

export function addAccount(
	user: { id: string; email: string; username: string },
	accessToken: string,
	refreshToken: string
) {
	update((state) => {
		const existingIndex = state.accounts.findIndex((account) => account.id === user.id);
		const account: Account = { ...user, accessToken, refreshToken };

		if (existingIndex !== -1) {
			state.accounts[existingIndex] = account;
			return { accounts: state.accounts, activeIndex: existingIndex };
		}

		const newAccounts = [...state.accounts, account];
		const newState = { accounts: newAccounts, activeIndex: newAccounts.length - 1 };
		saveToStorage(newState);
		return newState;
	});
}

export function removeAccount(id: string) {
	update((state) => {
		const index = state.accounts.findIndex((account) => account.id === id);
		if (index === -1) return state;

		const newAccounts = state.accounts.filter((account) => account.id !== id);
		let newActiveIndex = state.activeIndex;

		if (newAccounts.length === 0) {
			newActiveIndex = 0;
		} else if (index <= state.activeIndex) {
			newActiveIndex = Math.max(0, state.activeIndex - 1);
		}

		const newState = { accounts: newAccounts, activeIndex: newActiveIndex };
		saveToStorage(newState);
		return newState;
	});
}

export function switchAccount(id: string) {
	update((state) => {
		const index = state.accounts.findIndex((account) => account.id === id);
		if (index === -1) return state;
		const newState = { ...state, activeIndex: index };
		saveToStorage(newState);
		return newState;
	});
}

export async function refreshActiveTokens(): Promise<void> {
	const account = getActiveAccount();
	if (!account) throw new Error('No active account');

	const response = await fetch('/api/v1/auth/refresh', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ refresh_token: account.refreshToken })
	});

	if (!response.ok) {
		removeAccount(account.id);
		throw new Error('Token refresh failed');
	}

	const data: { access_token: string; refresh_token: string } = await response.json();

	update((state) => {
		const newAccounts = state.accounts.map((acc) =>
			acc.id === account.id
				? { ...acc, accessToken: data.access_token, refreshToken: data.refresh_token }
				: acc
		);
		const newState = { ...state, accounts: newAccounts };
		saveToStorage(newState);
		return newState;
	});
}
