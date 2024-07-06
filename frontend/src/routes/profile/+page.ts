import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { formSchema } from "./(components)/schema.js";
import { redirectIfNoAuthCookie } from "$lib/functions/auth.js";
import type { PageLoad } from "./$types.js";
import { request } from "$lib/functions/request.js";
import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
import { capitalizeFirstLetter } from "$lib/functions/string";

export const load: PageLoad = async ({}) => {
    
    redirectIfNoAuthCookie();

    const apiData = await request(`${PUBLIC_BACKEND_BASE_URL}/v1/users/me`)

    const formData = {
        goal: capitalizeFirstLetter(apiData.data.goal),
    }

    return { 
        title: "Profile",
        renderToolbar: true,
        form: await superValidate(formData, zod(formSchema)),
    };
};