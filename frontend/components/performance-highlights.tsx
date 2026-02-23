"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";

interface Perf { onTimeDeliveryRate: number; fleetUtilization: number; warehouseUtilization: number; orderFulfillmentRate: number; }

export default function PerformanceHighlights() {
  const [perf, setPerf] = useState<Perf | null>(null);

  useEffect(() => {
    gql<{ dashboardPerformance: Perf }>(`{ dashboardPerformance { onTimeDeliveryRate fleetUtilization warehouseUtilization orderFulfillmentRate } }`)
      .then((d) => setPerf(d.dashboardPerformance)).catch(() => {});
  }, []);

  const metrics = [
    { label: "On-time Delivery Rate", value: `${(perf?.onTimeDeliveryRate ?? 0).toFixed(1)}%`, color: "text-green-600" },
    { label: "Fleet Utilization", value: `${(perf?.fleetUtilization ?? 0).toFixed(1)}%`, color: "text-blue-600" },
    { label: "Warehouse Utilization", value: `${(perf?.warehouseUtilization ?? 0).toFixed(1)}%`, color: "text-amber-600" },
    { label: "Order Fulfillment", value: `${(perf?.orderFulfillmentRate ?? 0).toFixed(1)}%`, color: "text-purple-600" },
  ];

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="font-semibold text-gray-900 mb-4">Performance Highlights</h3>
      <div className="space-y-3">
        {metrics.map((m) => (
          <div key={m.label} className="flex items-center justify-between">
            <span className="text-sm text-gray-600">{m.label}</span>
            <span className={`text-sm font-bold ${m.color}`}>{m.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
