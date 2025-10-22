import { shipmentSchema } from '../pages/ShipmentPage.tsx'

export type ShipmentForm = z.infer<typeof shipmentSchema>;
