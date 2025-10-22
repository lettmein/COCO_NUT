import { Card, Input, Button } from '@/shared/ui/index';
import {useRouteCreate} from './model/useRouteCreate'

interface Props {
  points: { latitude: number; longitude: number; address?: string }[];
}

export default function RouteForm({ points }: Props) {
  const {register, handleSubmit, errors, mutation} = useRouteCreate({points})
  
  return (
    <Card className="p-4 space-y-4">
      <h2>Форма маршрута</h2>

      <Input
        type="number"
        {...register('maxVolume',  { valueAsNumber: true })}
        placeholder="Допустимый объем"
      />
      {errors.maxVolume && <span className="text-red-500">{errors.maxVolume.message}</span>}

      <Input
        type="number"
        {...register('maxWeight',  { valueAsNumber: true })}
        placeholder="Допустимый вес"
      />
      {errors.maxWeight && <span className="text-red-500">{errors.maxWeight.message}</span>}

      <Input
        type="datetime-local"
        {...register('departureDate')}
        placeholder="Дата отправления"
      />
      {errors.departureDate && <span className="text-red-500">{errors.departureDate.message}</span>}

      <div className="space-y-1">
        <h3>Точки маршрута:</h3>
        {points.map((p, i) => (
          <div key={i} className="text-sm">
            {i + 1}. {p.latitude.toFixed(5)}, {p.longitude.toFixed(5)} {p.address && `– ${p.address}`}
          </div>
        ))}
      </div>

      <Button onClick={handleSubmit} disabled={points.length < 2 || mutation.isPending}>
        {mutation.isPending ? 'Отправка...' : 'Сохранить маршрут'}
      </Button>
    </Card>
  );
}
