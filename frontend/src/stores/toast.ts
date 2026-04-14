import { writable } from 'svelte/store';

export const toastError = writable<string | null>(null);

export function showError(message: string) {
	toastError.set(message);
}
