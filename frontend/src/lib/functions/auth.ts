import { goto } from "$app/navigation";
import { getCookieValue, hasAuthCookie } from "./cookie";

export function getAuthToken(): string {
    const token = getCookieValue("token");
    if (token == null) {
        return "";
    }
    return token;
}

export function redirectIfNoAuthCookie() {
    if (!hasAuthCookie()) {
        goto("/login");
    }
}