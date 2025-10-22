import { shipmentSchema } from '../pages/ShipmentPage.tsx';
import {z} from 'zod';

export type ShipmentForm = z.infer<typeof shipmentSchema>;
