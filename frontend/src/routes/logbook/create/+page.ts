import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./schema.js";
import { zod } from "sveltekit-superforms/adapters";
import { redirectIfNoAuthCookie } from "$lib/functions/auth.js";

 
export const load: PageLoad = async () => {
  
  redirectIfNoAuthCookie();

  return {
    title: "Add check-in",
    form: await superValidate(zod(formSchema))
  };
};