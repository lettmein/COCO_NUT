import {Card, CardContent, CardHeader, CardTitle} from "../ui/card.tsx";
import {Label} from "../ui/label.tsx";
import {Input} from "../ui/input.tsx";

function SenderForm({ register, errors }: any) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Данные заказчика</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <div>
                    <Label>Наименование организации</Label>
                    <Input {...register('sender.organization')} />
                    {errors?.sender?.organization && <p className="text-red-600">{errors.sender.organization.message}</p>}
                </div>
                <div>
                    <Label>ИНН</Label>
                    <Input {...register('sender.inn')} maxLength={12} />
                    {errors?.sender?.inn && <p className="text-red-600">{errors.sender.inn.message}</p>}
                </div>
                <div>
                    <Label>ФИО ответственного</Label>
                    <Input {...register('sender.contactName')} />
                    {errors?.sender?.contactName && <p className="text-red-600">{errors.sender.contactName.message}</p>}
                </div>
                <div>
                    <Label>Контактный телефон</Label>
                    <Input {...register('sender.contactPhone')} placeholder="+7 (___) ___-__-__" />
                    {errors?.sender?.contactPhone && <p className="text-red-600">{errors.sender.contactPhone.message}</p>}
                </div>
                <div className="md:col-span-2">
                    <Label>Контактный email</Label>
                    <Input {...register('sender.contactEmail')} />
                    {errors?.sender?.contactEmail && <p className="text-red-600">{errors.sender.contactEmail.message}</p>}
                </div>
            </CardContent>
        </Card>
    );
}
