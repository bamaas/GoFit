import { redirectIfNoAuthCookie } from "$lib/functions/auth.js";
import type { PageLoad } from "./$types.js";

export const load: PageLoad = async ({}) => {
    
    redirectIfNoAuthCookie();

    return { 
        title: "Profile"
    };
};