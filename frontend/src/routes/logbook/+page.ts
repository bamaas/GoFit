import { goto } from "$app/navigation";
import type { PageLoad } from "./$types.js";
import { redirectIfNoAuthCookie } from "$lib/functions/auth";

export const load: PageLoad = async ({}) => {

    redirectIfNoAuthCookie();

    return { 
        title: "Logbook",
        renderToolbar: true,
    };
};