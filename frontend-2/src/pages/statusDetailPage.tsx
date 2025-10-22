import { getStatusById } from '@/modules/statusModule/api/api.status';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

export const RouteDetailPage = () => {
    const { id } = useParams<{ id: string }>();

    const { data: route, isLoading, error } = useQuery({
        queryKey: ['route', id],
        queryFn: () => getStatusById(id!),
        enabled: !!id,
    });

    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div>Ошибка при загрузке данных: {error.message}</div>;
    }

    if (!route) {
        return <div>Маршрут не найден</div>;
    }

    return (
        <div>
            <h1>Детали маршрута</h1>
            <div style={{ marginBottom: '20px', border: '1px solid #ccc', padding: '10px' }}>
                <h3>Маршрут ID: {route.id}</h3>
                <p><strong>Статус:</strong> {route.status}</p>
                <p><strong>Максимальный объем:</strong> {route.max_volume}</p>
                <p><strong>Максимальный вес:</strong> {route.max_weight}</p>
                <p><strong>Текущий объем:</strong> {route.current_volume}</p>
                <p><strong>Текущий вес:</strong> {route.current_weight}</p>
                <p><strong>Дата отправления:</strong> {new Date(route.departure_date).toLocaleString()}</p>
                <p><strong>Создан:</strong> {new Date(route.created_at).toLocaleString()}</p>
                <p><strong>Обновлен:</strong> {new Date(route.updated_at).toLocaleString()}</p>
                <h4>Точки маршрута:</h4>
                <ul>
                    {route.route_points.map((point) => (
                        <li key={point.address}>
                            <p><strong>Адрес:</strong> {point.address}</p>
                            <p><strong>Координаты:</strong> {point.latitude}, {point.longitude}</p>
                            <p><strong>Время прибытия:</strong> {new Date(point.arrival_time).toLocaleString()}</p>
                        </li>
                    ))}
                </ul>
                <p><strong>Request IDs:</strong> {route.request_ids.length > 0 ? route.request_ids.join(', ') : 'Нет'}</p>
            </div>
        </div>
    );
};
