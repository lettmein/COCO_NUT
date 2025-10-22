import {z} from "zod/index";

export type ShipmentForm = z.infer<typeof shipmentSchema>;