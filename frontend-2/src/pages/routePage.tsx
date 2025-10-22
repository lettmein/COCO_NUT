import { useState } from 'react';


import RouteForm from '@/modules/routeForm/routeForm';
import YandexRouteMap from '@/modules/yandexRouteMap';

export const RoutePage = () =>{ 
     const [points, setPoints] = useState<
    { latitude: number; longitude: number; address?: string }[]
  >([]);
    return (
        <section>

            <h1>A7 Логистика: оформление маршрута</h1>
           <YandexRouteMap points={points} setPoints={setPoints} />
           <RouteForm points={points} />
        </section>
    )
}