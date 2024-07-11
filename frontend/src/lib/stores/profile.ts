import { derived, writable } from 'svelte/store';

export type UserProfile = {
    id: number;
    createdAt: string;
    email: string;
    activated: boolean;
    goal: string;
};

export const profileStore = writable<UserProfile | undefined>(undefined);

export const profileReadable = derived(profileStore, ($profileStore) => {
    return $profileStore;
});