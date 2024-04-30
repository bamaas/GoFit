import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./schema.js";
import { zod } from "sveltekit-superforms/adapters";
import { goto } from "$app/navigation";

 
export const load: PageLoad = async () => {
  
  // Get auth token
  if (document.cookie == "") {
    goto("/login")
  }
  const authToken: string = document.cookie.split('=')[1]; 

  return {
    title: "Add check-in",
    form: await superValidate(zod(formSchema)),
    authToken: authToken
  };
};