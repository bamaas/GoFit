import { writable } from 'svelte/store';

export const authenticated = writable<boolean>(false);
