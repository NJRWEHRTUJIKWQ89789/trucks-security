"use client";

import { useEffect, useState } from "react";
import { Truck, Loader2, Wrench, DollarSign } from "lucide-react";
import { gql } from "@/lib/graphql";

interface MonthlyData { month: string; activeVehicles: number; maintenanceCost: number; utilizationRate: number; }
interface FleetReport { totalVehicles: number; averageUtilization: number; totalMaintenanceCost: number; monthlyData: MonthlyData[]; }

export default function FleetReportPage() {
  const [report, setReport] = useState<FleetReport | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ fleetReport: FleetReport }>(`{ fleetReport(year:2026) { totalVehicles averageUtilization totalMaintenanceCost monthlyData { month activeVehicles maintenanceCost utilizationRate } } }`)
      .then((d) => setReport(d.fleetReport)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Fleet Efficiency</h1>
      <div className="grid grid-cols-4 gap-4 mb-6">
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><Truck className="h-5 w-5 text-blue-500" /><span className="text-sm text-gray-500">Total Vehicles</span></div><p className="text-3xl font-bold">{report?.totalVehicles}</p></div>
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><Truck className="h-5 w-5 text-green-500" /><span className="text-sm text-gray-500">Avg Utilization</span></div><p className="text-3xl font-bold text-green-600">{report?.averageUtilization?.toFixed(1)}%</p></div>
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><DollarSign className="h-5 w-5 text-red-500" /><span className="text-sm text-gray-500">Total Maint. Cost</span></div><p className="text-3xl font-bold text-red-600">${report?.totalMaintenanceCost?.toLocaleString()}</p></div>
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><Wrench className="h-5 w-5 text-amber-500" /><span className="text-sm text-gray-500">Monthly Avg</span></div><p className="text-3xl font-bold text-amber-600">${report?.monthlyData?.length ? Math.round(report.totalMaintenanceCost / report.monthlyData.length).toLocaleString() : 0}</p></div>
      </div>
      <div className="bg-white rounded-xl border p-6">
        <h3 className="font-semibold text-gray-900 mb-4">Monthly Maintenance Cost</h3>
        <div className="space-y-2">
          {report?.monthlyData?.map((m) => (
            <div key={m.month} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-sm font-medium text-gray-700">{m.month}</span>
              <div className="flex items-center gap-4"><span className="text-xs text-gray-400">{m.activeVehicles} active, {m.utilizationRate?.toFixed(0)}% util</span><span className="text-sm font-bold text-red-600">${m.maintenanceCost?.toLocaleString()}</span></div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
