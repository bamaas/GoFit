import { z } from "zod";
 
export const formSchema = z.object({
    date: z.date(),
    weight: z.coerce.number().min(20).max(200),
});
 
export type FormSchema = typeof formSchema;
