import { useEffect, useRef } from 'react';
import { YMaps, Map, useYMaps } from '@pbe/react-yandex-maps';

interface Props {
  points: { latitude: number; longitude: number; address?: string }[];
  setPoints: (points: { latitude: number; longitude: number; address?: string }[]) => void;
}

export default function YandexRouteMap({ points, setPoints }: Props) {
  const mapRef = useRef<any>(null);

  const handleMapClick = (e: any) => {
    const coords = e.get('coords');
    setPoints(prev => {
      if (prev.length < 2)
        return [...prev, { latitude: coords[0], longitude: coords[1] }];
      return [{ latitude: coords[0], longitude: coords[1] }]; 
    });
  };

  return (
    <YMaps query={{ apikey: 'YOUR_API_KEY', lang: 'ru_RU' }}>
      <InnerMap points={points} setPoints={setPoints} mapRef={mapRef} handleMapClick={handleMapClick} />
    </YMaps>
  );
}

function InnerMap({ points, setPoints, mapRef, handleMapClick }: any) {
  const ymaps = useYMaps(['multiRouter.MultiRoute']);

  useEffect(() => {
    if (points.length >= 2 && mapRef.current && ymaps) {
      const route = new ymaps.multiRouter.MultiRoute(
        {
          referencePoints: points.map(p => [p.latitude, p.longitude]),
          params: { routingMode: 'auto' },
        },
        {
          boundsAutoApply: true,
          wayPointStartIconColor: 'green',
          wayPointFinishIconColor: 'red',
          routeStrokeWidth: 5,
          routeStrokeColor: '#3b82f6',
        }
      );

      mapRef.current.geoObjects.removeAll();
      mapRef.current.geoObjects.add(route);
    }
  }, [points, ymaps]);

  return (
    <Map
      instanceRef={mapRef}
      defaultState={{ center: [55.75, 37.57], zoom: 6 }}
      width="100%"
      height="400px"
      onClick={handleMapClick}
    />
  );
}
