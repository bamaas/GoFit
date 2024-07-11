import type { PageLoad } from "./$types.js";
import { getAuthToken } from "$lib/functions/auth";
import { authenticated } from "$lib/stores/auth";

export const ssr = false;
export const prerender = true;

export const load: PageLoad = async ({}) => {

    if (getAuthToken() != ""){
        authenticated.set(true)
    }

};