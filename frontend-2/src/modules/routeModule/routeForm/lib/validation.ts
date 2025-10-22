import { z } from 'zod';

export const routeFormSchema = z.object({
  maxVolume: z.number().min(0, 'Введите допустимый объем'),
  maxWeight: z.number().min(0, 'Введите допустимый вес'),
  departureDate: z.string().min(1, 'Выберите дату отправления'),
});

export type RouteFormData = z.infer<typeof routeFormSchema>;
