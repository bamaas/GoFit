import { z } from "zod";

export const formSchema = z.object({
    uuid: z.string(),
    date: z.coerce.date().default(new Date()),
    weight: z.coerce.number().min(20).max(200).default('' as unknown as number)
});


 
export type FormSchema = typeof formSchema;