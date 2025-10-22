import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';

import type { AxiosError } from 'axios';
import { routeFormSchema, type RouteFormData } from '../lib/validation';
import { routeCreate } from '../api/api.route';

interface Props {
  points: { latitude: number; longitude: number; address?: string }[];
}

export const useRouteCreate = ({ points }: Props) => {
  const { register, handleSubmit, formState: { errors }, reset , control} = useForm<RouteFormData>({
    resolver: zodResolver(routeFormSchema),
    defaultValues: {
      maxVolume: 0,
      maxWeight: 0,
      departureDate: '',
    }
  });

  const mutation = useMutation({
    mutationFn: (payload: RouteFormData) => routeCreate(payload),
    onSuccess: () => {
      alert('Маршрут успешно создан');
      reset();
    },
    onError: (err: AxiosError) => {
      console.error(err);
      alert('Ошибка при создании маршрута');
    }
  });

  const onSubmit = (data: RouteFormData) => {
    if (points.length < 2) {
      alert('Выберите хотя бы две точки на карте');
      return;
    }

    const payload = {
      max_volume: data.maxVolume,
      max_weight: data.maxWeight,
      departure_date: new Date(data.departureDate).toISOString(),
      route_points: points
    };

    mutation.mutate(payload);
  };

  return {
    register,
    handleSubmit: handleSubmit(onSubmit),
    errors,
    mutation
    ,control
  };
};
