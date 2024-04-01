import { z } from "zod";

export const formSchema = z.object({
    uuid: z.string(),
    date: z.coerce.date().default(new Date()),
    weight: z.coerce.number().min(20).max(200).default(0),
});

const crudSchema = formSchema.extend({
    uuid: formSchema.shape.uuid.optional()
})
 
export type FormSchema = typeof formSchema;