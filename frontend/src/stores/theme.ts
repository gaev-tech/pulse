import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type ThemeMode = 'light' | 'dark' | 'system';

const STORAGE_KEY = 'pulse_theme';

function resolveTheme(mode: ThemeMode): 'light' | 'dark' {
	if (mode === 'system') {
		return browser && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
	}
	return mode;
}

function applyTheme(mode: ThemeMode) {
	if (!browser) return;
	const resolved = resolveTheme(mode);
	document.documentElement.classList.toggle('dark', resolved === 'dark');
}

export function getTheme(): ThemeMode {
	if (!browser) return 'system';
	return (localStorage.getItem(STORAGE_KEY) as ThemeMode) ?? 'system';
}

export function setTheme(mode: ThemeMode) {
	if (!browser) return;
	localStorage.setItem(STORAGE_KEY, mode);
	applyTheme(mode);
	theme.set(mode);
}

export const theme = writable<ThemeMode>(getTheme());

export function initTheme() {
	const mode = getTheme();
	applyTheme(mode);
	if (browser) {
		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
			if (getTheme() === 'system') applyTheme('system');
		});
	}
}
