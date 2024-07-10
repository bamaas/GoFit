import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
import { request } from "$lib/functions/request.js";
import { profileStore, type UserProfile } from "$lib/stores/profile";

export function fetchUserProfile(): void {
    request(`${PUBLIC_BACKEND_BASE_URL}/v1/users/me`).then(response => {
        const user: UserProfile = {
            id: response.data.id,
            email: response.data.email,
            goal: response.data.goal,
            activated: response.data.activated,
            createdAt: response.data.created_at,
        }
        profileStore.set(user); 
    });
}