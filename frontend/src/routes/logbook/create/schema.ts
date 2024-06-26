import { z } from "zod";

export const formSchema = z.object({
    uuid: z.string(),
    date: z.string().min(1, "Please pick a date").default(new Date().toISOString().split('T')[0]),
    weight: z.coerce.number().min(20).max(200).multipleOf(0.1).default('' as unknown as number),
    notes: z.string().max(100, "Must be less than 100 characters").optional(),
});

export type FormSchema = typeof formSchema;