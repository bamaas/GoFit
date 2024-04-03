import { z } from "zod";

export const formSchema = z.object({
    uuid: z.string(),
    date: z.string().default(new Date().toISOString().split('T')[0]),
    weight: z.coerce.number().min(20).max(200).default('' as unknown as number)
});

export type FormSchema = typeof formSchema;