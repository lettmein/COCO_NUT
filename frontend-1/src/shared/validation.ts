import {z} from "zod";

export const cargoItemSchema = z.object({
    name: z.string().min(1, "Введите наименование груза"),
    quantity: z.number().positive("Количество должно быть положительным числом"),
    weight: z.number().nonnegative("Вес не может быть отрицательным"),
    volume: z.number().nonnegative("Объём не может быть отрицательным"),
    specialRequirements: z.string().optional(),
});

export const shipmentSchema = z.object({
    logisticsPointId: z.string().min(1, "Выберите логистическую точку"),
    sender: z.object({
        organization: z.string().min(1, "Введите наименование организации заказчика"),
        inn: z.string().min(10, "ИНН должен содержать не менее 10 символов").max(12),
        contactName: z.string().min(1, "Введите ФИО ответственного"),
        contactPhone: z.string().optional(),
        contactEmail: z.string().email("Неверный формат email").optional(),
    }),
    receiver: z.object({
        organization: z.string().min(1, "Введите наименование организации получателя"),
        address: z.string().min(1, "Введите адрес точки получения"),
        contactName: z.string().min(1, "Введите ФИО ответственного у получателя"),
        contactPhone: z.string().min(5, "Введите контактный телефон"),
    }),
    cargo: z.array(cargoItemSchema).min(1, "Добавьте хотя бы одну позицию груза"),
});