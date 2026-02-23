"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";

interface Stats { totalVehicles: number; activeVehicles: number; }

export default function FleetStatus() {
  const [stats, setStats] = useState<Stats | null>(null);

  useEffect(() => {
    gql<{ dashboardStats: Stats }>(`{ dashboardStats { totalVehicles activeVehicles } }`)
      .then((d) => setStats(d.dashboardStats)).catch(() => {});
  }, []);

  const total = stats?.totalVehicles || 1;
  const active = stats?.activeVehicles || 0;
  const maintenance = Math.round(total * 0.09);
  const idle = total - active - maintenance;

  const statuses = [
    { label: "Active", count: active, percentage: Math.round((active / total) * 100), color: "bg-blue-500" },
    { label: "Idle", count: idle > 0 ? idle : 0, percentage: Math.round((Math.max(idle, 0) / total) * 100), color: "bg-amber-400" },
    { label: "Maintenance", count: maintenance, percentage: Math.round((maintenance / total) * 100), color: "bg-red-400" },
  ];

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="font-semibold text-gray-900 mb-4">Fleet Status</h3>
      <div className="space-y-4">
        {statuses.map((s) => (
          <div key={s.label}>
            <div className="flex items-center justify-between text-sm mb-1">
              <span className="text-gray-600">{s.label}</span>
              <span className="font-medium text-gray-900">{s.count} ({s.percentage}%)</span>
            </div>
            <div className="h-2 bg-gray-100 rounded-full">
              <div className={`h-2 rounded-full ${s.color}`} style={{ width: `${s.percentage}%` }} />
            </div>
          </div>
        ))}
      </div>
      <p className="text-xs text-gray-400 mt-4">Total: {total} vehicles</p>
    </div>
  );
}
