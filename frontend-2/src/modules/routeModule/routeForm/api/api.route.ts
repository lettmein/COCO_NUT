import { instance } from "@/shared/api/instance";
import type { RouteFormData } from "../lib/validation";

export const routeCreate = async (data: RouteFormData) => {
    const response = await instance.post('/api/v1/routes', data)
    return response.data
}