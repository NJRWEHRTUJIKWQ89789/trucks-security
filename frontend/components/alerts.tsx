"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";
import { AlertTriangle } from "lucide-react";

interface Shipment { id: string; trackingNumber: string; destination: string; estimatedDelivery: string; }

export default function Alerts() {
  const [delayed, setDelayed] = useState<Shipment[]>([]);

  useEffect(() => {
    gql<{ shipments: { items: Shipment[] } }>(`{ shipments(status: "delayed", page: 1, perPage: 5) { items { id trackingNumber destination estimatedDelivery } } }`)
      .then((d) => setDelayed(d.shipments?.items || [])).catch(() => {});
  }, []);

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="font-semibold text-gray-900 mb-4">Alerts</h3>
      <div className="space-y-3">
        {delayed.length === 0 && <p className="text-sm text-gray-400">No alerts</p>}
        {delayed.map((s) => (
          <div key={s.id} className="flex items-start gap-3 p-3 bg-red-50 rounded-lg">
            <AlertTriangle className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
            <div>
              <p className="text-sm font-medium text-red-800">{s.trackingNumber} delayed</p>
              <p className="text-xs text-red-600">To: {s.destination}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
