export interface RoutePoint {
  latitude: number;
  longitude: number;
  address: string;
  arrival_time: string; 
}

export interface Route {
  id: string; 
  max_volume: number;
  max_weight: number;
  current_volume: number;
  current_weight: number;
  departure_date: string;
  route_points: RoutePoint[];
  status: string; 
  request_ids: string[];
  created_at: string;
  updated_at: string;
}

export type RoutesResponse = Route[];
