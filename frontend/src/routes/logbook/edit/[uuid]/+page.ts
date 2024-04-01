import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "../../create/schema.js";
import { zod } from "sveltekit-superforms/adapters";
import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
import { toast } from "svelte-sonner";

export const prerender = false;

export const load: PageLoad = async ({params}) => {
    let response = await fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins/${params.uuid}`)
    let data = await response.json()
    const form = await superValidate(data, zod(formSchema));
    return { form };
};