import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "../../create/schema.js";
import { zod } from "sveltekit-superforms/adapters";
import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";
import type { CheckIn } from "../../store.js";
import { goto } from "$app/navigation";

export const prerender = false;

export const load: PageLoad = async ({fetch, params}) => {

    // Get auth token
    if (document.cookie == "") {
        goto("/login")
    }
    const authToken: string = document.cookie.split('=')[1]; 

    // Get check-in data
    const response = await fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/check-ins/${params.uuid}`,
        {
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }
    )
    const apiData: CheckIn = await response.json()
    // Need to transform the datetime to a format that the form data field 'date' can understand
    const formData = {
        uuid: apiData.uuid,
        date: new Date(apiData.datetime*1000).toISOString().split('T')[0],
        weight: apiData.weight,
        notes: apiData.notes
    }
    const form = await superValidate(formData, zod(formSchema));

    return {
        title: "Edit check-in", 
        form: form,
        authToken: authToken
    };
};