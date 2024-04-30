import { goto } from "$app/navigation";
import type { PageLoad } from "./$types.js";

export const load: PageLoad = async ({}) => {
    
    // Get auth token
    if (document.cookie == "") {
        goto("/login")
    }
    const authToken: string = document.cookie.split('=')[1];    

    // Return page data
    return { 
        title: "Summary",
        renderToolbar: true,
        authToken: authToken
    };
};