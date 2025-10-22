import {useMutation, useQueryClient} from "@tanstack/react-query";
import {Controller, useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {Card, CardContent, CardHeader, CardTitle} from "../components/ui/card.tsx";
import {Button} from "../components/ui/button.tsx";
import {Separator} from "../components/ui/separator.tsx";

import {  }

export default function CreateShipmentPage() {
    const queryClient = useQueryClient();

    const form = useForm<ShipmentForm>({
        resolver: zodResolver(shipmentSchema),
        defaultValues: {
            logisticsPointId: "",
            sender: { organization: '', inn: '', contactName: '', contactPhone: '', contactEmail: '' },
            receiver: { organization: '', address: '', contactName: '', contactPhone: '' },
            cargo: [{ name: '', quantity: 1, weight: 0, volume: 0, specialRequirements: '' }],
        },
    });

    const { control, register, handleSubmit, formState } = form;
    const { errors, isSubmitting } = formState;

    const createMutation = useMutation((payload: ShipmentForm) => createShipmentApi(payload), {
        onSuccess: (data) => {
            toast({ title: 'Заявка создана', description: `Заявка №${data.id} успешно создана` });
            queryClient.invalidateQueries(['shipments']);
            // Reset cargoSpecificationForm.tsx or redirect as needed
        },
        onError: (err: any) => {
            toast({ title: 'Ошибка', description: err.message || 'Не удалось создать заявку' });
        },
    });

    const onSubmit = async (values: ShipmentForm) => {
        // Преобразование данных при необходимости
        await createMutation.mutateAsync(values);
    };

    return (
        <div className="max-w-5xl mx-auto p-6 space-y-6">
            <h1 className="text-2xl font-semibold">Создание заявки на транспортировку</h1>

            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                <Card>
                    <CardHeader>
                        <CardTitle>Параметры отправки</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                            <div className="md:col-span-2">
                                <Controller name="logisticsPointId" control={control} render={({ field }) => <LogisticsPointSelect control={control} />} />
                                {errors?.logisticsPointId && <p className="text-red-600">{(errors as any).logisticsPointId?.message}</p>}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <SenderForm register={register} errors={errors} />
                <ReceiverForm register={register} errors={errors} />
                <CargoSpecificationForm control={control} register={register} errors={errors} />

                <div className="flex items-center justify-end gap-2">
                    <Button type="submit" disabled={isSubmitting || createMutation.isLoading}>Создать заявку</Button>
                </div>
            </form>

            <Separator />

            <div className="text-sm text-muted-foreground">Примечание: данный интерфейс ожидает реализацию backend API по маршрутам <code>/api/logistics-points</code> и <code>/api/shipments</code>. Для разработки можно использовать MSW или JSON-server.</div>
        </div>
    );
}
