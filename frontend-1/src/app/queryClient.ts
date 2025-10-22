import { useMutation } from "@tanstack/react-query";
import { createShipment } from "../api/shipment.api";

export const useCreateShipment = () =>
    useMutation({
        mutationFn: createShipment,
        onSuccess: () => toast.success("Заявка успешно создана"),
        onError: (err) => toast.error("Ошибка при создании заявки"),
    });


export const createShipment = async (data: ShipmentFormData) => {
    const res = await fetch("/api/shipments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to create shipment");
    return res.json();
};
