import { useEffect, useRef } from 'react';
import { YMaps, Map, useYMaps } from '@pbe/react-yandex-maps';

interface Props {
  points: [number, number][];
  setPoints: (points: [number, number][]) => void;
}

export default function YandexRouteMap({ points, setPoints }: Props) {
  const mapRef = useRef<any>(null);
  const ymaps = useYMaps(['multiRouter.MultiRoute']);

  const handleMapClick = (e: any) => {
    const coords = e.get('coords');
    setPoints(prev => {
      if (prev.length < 2) return [...prev, coords];
      return [coords]; 
    });
  };

  useEffect(() => {
    if (points.length === 2 && mapRef.current && ymaps) {
      const route = new ymaps.multiRouter.MultiRoute(
        {
          referencePoints: points,
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
    <YMaps query={{ apikey: 'YOUR_API_KEY', lang: 'ru_RU' }}>
      <Map
        instanceRef={mapRef}
        defaultState={{ center: [55.75, 37.57], zoom: 6 }}
        width="100%"
        height="400px"
        onClick={handleMapClick}
      />
    </YMaps>
  );
}
