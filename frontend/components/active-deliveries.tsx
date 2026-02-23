"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";
import { Truck } from "lucide-react";

interface Shipment { id: string; trackingNumber: string; origin: string; destination: string; status: string; carrier: string; }

export default function ActiveDeliveries() {
  const [shipments, setShipments] = useState<Shipment[]>([]);

  useEffect(() => {
    gql<{ shipments: { items: Shipment[] } }>(`{ shipments(page:1,perPage:5,status:"in_transit") { items { id trackingNumber origin destination status carrier } } }`)
      .then((d) => setShipments(d.shipments.items)).catch(() => {});
  }, []);

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="font-semibold text-gray-900 mb-4">Active Deliveries</h3>
      <div className="space-y-3">
        {shipments.length === 0 && <p className="text-sm text-gray-400">No active deliveries</p>}
        {shipments.map((s) => (
          <div key={s.id} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
            <Truck className="h-5 w-5 text-blue-500 flex-shrink-0" />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900 truncate">{s.trackingNumber}</p>
              <p className="text-xs text-gray-500">{s.origin} → {s.destination}</p>
            </div>
            <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">{s.carrier}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
