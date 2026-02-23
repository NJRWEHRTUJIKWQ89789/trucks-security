"use client";

import { useEffect, useState } from "react";
import { Truck, Loader2, CheckCircle, Clock, AlertTriangle } from "lucide-react";
import { gql } from "@/lib/graphql";

interface MonthlyData { month: string; totalDeliveries: number; onTime: number; late: number; onTimeRate: number; }
interface DeliveryReport { totalDeliveries: number; averageOnTimeRate: number; averageDeliveryTime: number; monthlyData: MonthlyData[]; }

export default function DeliveryPage() {
  const [report, setReport] = useState<DeliveryReport | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ deliveryReport: DeliveryReport }>(`{ deliveryReport(year:2026) { totalDeliveries averageOnTimeRate averageDeliveryTime monthlyData { month totalDeliveries onTime late onTimeRate } } }`)
      .then((d) => setReport(d.deliveryReport)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Delivery Performance</h1>
      <div className="grid grid-cols-4 gap-4 mb-6">
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><Truck className="h-5 w-5 text-blue-500" /><span className="text-sm text-gray-500">Total Deliveries</span></div><p className="text-3xl font-bold">{report?.totalDeliveries}</p></div>
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><CheckCircle className="h-5 w-5 text-green-500" /><span className="text-sm text-gray-500">On-Time Rate</span></div><p className="text-3xl font-bold text-green-600">{report?.averageOnTimeRate?.toFixed(1)}%</p></div>
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><Clock className="h-5 w-5 text-purple-500" /><span className="text-sm text-gray-500">Avg Delivery Time</span></div><p className="text-3xl font-bold text-purple-600">{report?.averageDeliveryTime?.toFixed(1)}h</p></div>
        <div className="bg-white rounded-xl border p-5"><div className="flex items-center gap-2 mb-2"><AlertTriangle className="h-5 w-5 text-red-500" /><span className="text-sm text-gray-500">Late (Total)</span></div><p className="text-3xl font-bold text-red-600">{report?.monthlyData?.reduce((s, m) => s + m.late, 0)}</p></div>
      </div>
      <div className="bg-white rounded-xl border p-6">
        <h3 className="font-semibold text-gray-900 mb-4">Monthly Deliveries</h3>
        <div className="space-y-2">
          {report?.monthlyData?.map((m) => (
            <div key={m.month} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-sm font-medium text-gray-700">{m.month}</span>
              <div className="flex items-center gap-4">
                <span className="text-xs text-gray-400">{m.onTime}/{m.totalDeliveries} on time</span>
                <span className="text-sm font-bold text-blue-600">{m.onTimeRate?.toFixed(0)}%</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
