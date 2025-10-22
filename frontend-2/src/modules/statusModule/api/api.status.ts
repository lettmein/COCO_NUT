import { instance } from "@/shared/api/instance"
import type { Route, RoutesResponse } from "../model/types.status";

export const getStatuses = async (): Promise<RoutesResponse> => {
    const response = await instance.get('/api/v1/routes');
        return response.data;
}

export const getStatusById = async (id: string): Promise<Route> => {
    const response = await instance.get(`/api/v1/routes/${id}`)
    return response.data
}