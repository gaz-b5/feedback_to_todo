import { z } from "zod";

export const taskSchema = z.object({
  id: z.string(),
  title: z.string(),
  status: z.string(),
  nature: z.string(),
  priority: z.string(),
  occurence: z
    .union([z.number(), z.string()])
    .optional()
    .transform((val) => {
      if (typeof val === "string") {
        const n = Number(val);
        if (Number.isNaN(n)) return undefined;
        return n;
      }
      return val;
    }),
});

export type Task = z.infer<typeof taskSchema>;
