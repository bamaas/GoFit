import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { formSchema } from "./(components)/schema.js";
import { redirectIfNoAuthCookie } from "$lib/functions/auth.js";
import type { PageLoad } from "./$types.js";

export const load: PageLoad = async ({}) => {
    
    redirectIfNoAuthCookie();

    return { 
        title: "Profile",
        renderToolbar: true,
        form: await superValidate(zod(formSchema)),
    };
};