import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./(components)/schema.js";
import { zod } from "sveltekit-superforms/adapters";
import { goto } from "$app/navigation";
import { hasAuthCookie } from "$lib/functions/cookie.js";
 
export const load: PageLoad = async () => {

    if (hasAuthCookie()) {
        goto("/logbook")
    }

    return {
        title: "Login",
        form: await superValidate(zod(formSchema)),
        renderHeader: false,
        renderToolbar: false
    };
};