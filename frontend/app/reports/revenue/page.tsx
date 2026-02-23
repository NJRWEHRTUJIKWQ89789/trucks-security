"use client";

import { useEffect, useState } from "react";
import { DollarSign, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface MonthlyData { month: string; revenue: number; orders: number; profit: number; }
interface RevenueReport { totalRevenue: number; totalProfit: number; monthlyData: MonthlyData[]; }

export default function RevenuePage() {
  const [report, setReport] = useState<RevenueReport | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ revenueReport: RevenueReport }>(`{ revenueReport(year:2026) { totalRevenue totalProfit monthlyData { month revenue orders profit } } }`)
      .then((d) => setReport(d.revenueReport)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Revenue Analysis</h1>
      <div className="bg-white rounded-xl border p-6 mb-6">
        <div className="flex items-center gap-3 mb-2"><DollarSign className="h-6 w-6 text-green-600" /><h2 className="text-lg font-semibold">Total Revenue</h2></div>
        <p className="text-4xl font-bold text-gray-900">${report?.totalRevenue?.toLocaleString("en-US", { minimumFractionDigits: 2 })}</p>
      </div>
      <div className="bg-white rounded-xl border p-6">
        <h3 className="font-semibold text-gray-900 mb-4">Monthly Breakdown</h3>
        <div className="space-y-3">
          {report?.monthlyData?.map((m) => (
            <div key={m.month} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-sm font-medium text-gray-700">{m.month}</span>
              <div className="flex items-center gap-4">
                <span className="text-xs text-gray-400">{m.orders} orders</span>
                <span className="text-sm font-bold text-green-600">${m.revenue?.toLocaleString()}</span>
              </div>
            </div>
          ))}
          {(!report?.monthlyData || report.monthlyData.length === 0) && <p className="text-sm text-gray-400">No data available</p>}
        </div>
      </div>
    </div>
  );
}
