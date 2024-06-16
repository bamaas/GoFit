import { writable } from 'svelte/store';

export const goal = writable<string>("lose");
