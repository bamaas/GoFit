import { z } from "zod";

export const formSchema = z.object({
    date: z.coerce.date().default(new Date()),
    weight: z.coerce.number().min(20).max(200).default(0),
});
 
export type FormSchema = typeof formSchema;