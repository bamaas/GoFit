import type { PageLoad } from "./$types.js";
import { superValidate } from "sveltekit-superforms";
import { formSchema } from "./schema.js";
import { zod } from "sveltekit-superforms/adapters";

 
export const load: PageLoad = async () => {
  return {
    title: "Add check-in",
    form: await superValidate(zod(formSchema)),
  };
};