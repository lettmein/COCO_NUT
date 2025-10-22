import React from "react";
import { Controller } from "react-hook-form";

import { useQuery } from "@tanstack/react-query";

import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../components/ui/select";
import { Label } from "../components/ui/label";

import { fetchLogisticsPoints } from '../shared/api/api.ts'


export function LogisticsPointSelect({ control }: { control: any }) {
    const [search, setSearch] = React.useState("");
    const { data = [], isLoading } = useQuery(["logisticsPoints", search], () => fetchLogisticsPoints(search), { staleTime: 350000 });

    return (
        <Controller
            control={control}
            name="logisticsPointId"
            render={({ field }) => (
                <div>
                    <Label>Логистическая точка отправки</Label>
                    <Select onValueChange={field.onChange} value={field.value}>
                        <SelectTrigger className="w-full">
                            <SelectValue placeholder={isLoading ? "Загрузка..." : "Выберите логистическую точку"} />
                        </SelectTrigger>
                        <SelectContent>
                            {data.map((pt: any) => (
                                <SelectItem key={pt.id} value={pt.id}>
                                    {pt.name} — {pt.address}
                                </SelectItem>
                            ))}
                            {data.length === 0 && !isLoading && <SelectItem value="">Нет доступных точек</SelectItem>}
                        </SelectContent>
                    </Select>
                    <div className="mt-2 flex gap-2">
                        <Input placeholder="Поиск точки (часть названия)" value={search} onChange={(e) => setSearch(e.target.value)} />
                        <Button variant="outline" onClick={() => setSearch("")}>Очистить</Button>
                    </div>
                </div>
            )}
        />
    );
}


