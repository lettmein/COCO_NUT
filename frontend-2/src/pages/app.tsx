import RouteForm from '@/modules/routeForm';
import YandexRouteMap from '@/modules/yandexRouteMap';
import { useState } from 'react';

export default function App() {
  const [points, setPoints] = useState<[number, number][]>([]);

  return (
    <div className="p-4 space-y-6">
      <h1>A7 Логистика: оформление маршрута</h1>
      <YandexRouteMap points={points} setPoints={setPoints} />
      <RouteForm points={points} />
    </div>
  );
}
