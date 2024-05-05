export function getCookieValue(name: string): string | null {
    const regex = new RegExp(`(^| )${name}=([^;]+)`)
    const match = document.cookie.match(regex)
    if (match) {
        return match[2]
    }
    return null
}

export function hasAuthCookie(): boolean {
    const cookie = getCookieValue("token");
    if (cookie == null) {
        return false;
    }
    return true;
}

export function deleteCookie(name: string) {
    document.cookie = `${name}=xxx; expires=Thu, 18 Dec 1970 12:00:00 UTC; path=/;`;
}