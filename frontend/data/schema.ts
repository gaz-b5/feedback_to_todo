import { z } from "zod";

export const taskSchema = z.object({
  id: z.string(),
  title: z.string(),
  status: z.string(),
  nature: z.string(),
  priority: z.coerce.string(),
  occurrence: z
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
  created: z.string().transform((val) => {
    try {
      const date = new Date(val);
      if (isNaN(date.getTime())) return val; // fallback to original if invalid

      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const hours = String(date.getHours()).padStart(2, "0");
      const minutes = String(date.getMinutes()).padStart(2, "0");

      return `${day}-${month}-${year} ${hours}:${minutes}`;
    } catch {
      return val;
    }
  }),
});

export type Task = z.infer<typeof taskSchema>;
