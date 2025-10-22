import {Card, CardContent, CardHeader, CardTitle} from "../ui/card.tsx";
import {Label} from "../ui/label.tsx";
import {Input} from "../ui/input.tsx";

export function ReceiverForm({ register, errors }: any) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Данные получателя</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <div>
                    <Label>Наименование организации</Label>
                    <Input {...register('receiver.organization')} />
                    {errors?.receiver?.organization && <p className="text-red-600">{errors.receiver.organization.message}</p>}
                </div>
                <div>
                    <Label>Адрес точки получения</Label>
                    <Input {...register('receiver.address')} />
                    {errors?.receiver?.address && <p className="text-red-600">{errors.receiver.address.message}</p>}
                </div>
                <div>
                    <Label>ФИО ответственного</Label>
                    <Input {...register('receiver.contactName')} />
                    {errors?.receiver?.contactName && <p className="text-red-600">{errors.receiver.contactName.message}</p>}
                </div>
                <div>
                    <Label>Контактный телефон</Label>
                    <Input {...register('receiver.contactPhone')} />
                    {errors?.receiver?.contactPhone && <p className="text-red-600">{errors.receiver.contactPhone.message}</p>}
                </div>
            </CardContent>
        </Card>
    );
}