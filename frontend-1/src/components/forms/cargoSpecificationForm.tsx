import {useMemo} from "react";

import { Card, CardHeader, CardTitle, CardContent, Label, Input, Button } from "../ui/index.ts"
import { useFieldArray } from "react-hook-form";

export function CargoSpecificationForm({ control, register, errors }: any) {
    const { fields, append, remove } = useFieldArray({ control, name: 'cargo' });

    const totals = useMemo(() => {
        let weight = 0;
        let volume = 0;
        let quantity = 0;
        fields.forEach((f: any) => {
            quantity += Number(f.quantity || 0);
            weight += Number(f.weight || 0);
            volume += Number(f.volume || 0);
        });
        return { weight, volume, quantity };
    }, [fields]);

    return (
        <Card>
            <CardHeader>
                <CardTitle>Спецификация груза</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                {fields.map((field: any, idx: number) => (
                    <div key={field.id} className="grid grid-cols-1 gap-2 md:grid-cols-6 items-end">
                        <div className="md:col-span-2">
                            <Label>Наименование</Label>
                            <Input {...register(`cargo.${idx}.name` as const)} defaultValue={field.name} />
                            {errors?.cargo?.[idx]?.name && <p className="text-red-600">{errors.cargo[idx].name.message}</p>}
                        </div>
                        <div>
                            <Label>Количество</Label>
                            <Input type="number" {...register(`cargo.${idx}.quantity` as const, { valueAsNumber: true })} defaultValue={field.quantity} />
                            {errors?.cargo?.[idx]?.quantity && <p className="text-red-600">{errors.cargo[idx].quantity.message}</p>}
                        </div>
                        <div>
                            <Label>Вес (кг)</Label>
                            <Input type="number" {...register(`cargo.${idx}.weight` as const, { valueAsNumber: true })} defaultValue={field.weight} />
                            {errors?.cargo?.[idx]?.weight && <p className="text-red-600">{errors.cargo[idx].weight.message}</p>}
                        </div>
                        <div>
                            <Label>Объём (м³)</Label>
                            <Input type="number" {...register(`cargo.${idx}.volume` as const, { valueAsNumber: true })} defaultValue={field.volume} />
                            {errors?.cargo?.[idx]?.volume && <p className="text-red-600">{errors.cargo[idx].volume.message}</p>}
                        </div>
                        <div className="md:col-span-1">
                            <Label>Особые требования</Label>
                            <Input {...register(`cargo.${idx}.specialRequirements` as const)} defaultValue={field.specialRequirements} />
                        </div>
                        <div className="md:col-span-1">
                            <Button variant="destructive" onClick={() => remove(idx)}>Удалить</Button>
                        </div>
                    </div>
                ))}

                <div className="flex gap-2">
                    <Button onClick={() => append({ name: '', quantity: 1, weight: 0, volume: 0, specialRequirements: '' })}>Добавить позицию</Button>
                    <div className="ml-auto text-sm">
                        <div>Всего: {totals.quantity} шт</div>
                        <div>Общий вес: {totals.weight} кг</div>
                        <div>Общий объём: {totals.volume} м³</div>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}