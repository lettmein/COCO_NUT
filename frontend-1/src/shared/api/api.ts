import type {ShipmentForm} from "../../types/types.ts";

export async function createShipmentApi(payload: ShipmentForm) {
    const res = await fetch('/api/shipments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
    });
    if (!res.ok) {
        const errorBody = await res.json().catch(() => ({}));
        throw new Error(errorBody.message || 'Failed to create shipment');
    }
    return res.json();
}

export async function fetchLogisticsPoints(search = "") {
    const qs = search ? `?search=${encodeURIComponent(search)}` : "";
    const res = await fetch(`/api/logistics-points${qs}`);
    if (!res.ok) throw new Error("Failed to load logistics points");
    return res.json();
}