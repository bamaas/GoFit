import { z } from "zod";
 
export const formSchema = z.object({
    goal: z.string({ required_error: "Please select a goal" }),
});
 
export type FormSchema = typeof formSchema;