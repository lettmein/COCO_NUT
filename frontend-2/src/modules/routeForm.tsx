import { useState } from 'react';
import { Card, Input, Textarea, Button } from '@/shared/ui/index';

interface Props {
  points: [number, number][];
}

export default function RouteForm({ points }: Props) {
  const [form, setForm] = useState({
    senderName: '',
    senderINN: '',
    senderPhone: '',
    cargoDescription: '',
    cargoWeight: '',
    receiverName: '',
    receiverAddress: '',
    receiverPhone: ''
  });

  const handleSubmit = () => {
    console.log('Маршрут:', form, points);
    alert('Данные отправлены в консоль');
  };

  return (
    <Card className="p-4 space-y-4">
      <h2>Форма маршрута</h2>
      <Input placeholder="Отправитель: Наименование" value={form.senderName} onChange={(e) => setForm({ ...form, senderName: e.target.value })} />
      <Input placeholder="ИНН" value={form.senderINN} onChange={(e) => setForm({ ...form, senderINN: e.target.value })} />
      <Input placeholder="Телефон/Email" value={form.senderPhone} onChange={(e) => setForm({ ...form, senderPhone: e.target.value })} />
      <Textarea placeholder="Описание груза" value={form.cargoDescription} onChange={(e) => setForm({ ...form, cargoDescription: e.target.value })} />
      <Input placeholder="Вес (кг)" value={form.cargoWeight} onChange={(e) => setForm({ ...form, cargoWeight: e.target.value })} />
      <Input placeholder="Получатель: Наименование" value={form.receiverName} onChange={(e) => setForm({ ...form, receiverName: e.target.value })} />
      <Input placeholder="Адрес" value={form.receiverAddress} onChange={(e) => setForm({ ...form, receiverAddress: e.target.value })} />
      <Input placeholder="Телефон" value={form.receiverPhone} onChange={(e) => setForm({ ...form, receiverPhone: e.target.value })} />
      <Button onClick={handleSubmit} disabled={points.length < 2}>Сохранить маршрут</Button>
    </Card>
  );
}
