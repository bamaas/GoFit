import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./schema.js";
import { zod } from "sveltekit-superforms/adapters";
import { goto } from "$app/navigation";
 
export const load: PageLoad = async () => {

    // Redirect to logbook if already logged in
    if (document.cookie != "") {
        goto("/logbook")
    }

    return {
        title: "Login",
        form: await superValidate(zod(formSchema))
    };
};