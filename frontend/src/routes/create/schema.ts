import { z } from "zod";

export const formSchema = z.object({
    date: z.string().default(new Date().toISOString().split('T')[0]),
    weight: z.coerce.number().min(20).max(200).default(0),
});
 
export type FormSchema = typeof formSchema;
