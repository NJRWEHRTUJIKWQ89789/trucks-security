"use client";

import { useEffect, useState } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import { Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface MonthlyDeliveryData {
  month: string;
  totalDeliveries: number;
  onTime: number;
  late: number;
  onTimeRate: number;
}

interface DeliveryReport {
  totalDeliveries: number;
  averageOnTimeRate: number;
  averageDeliveryTime: number;
  monthlyData: MonthlyDeliveryData[];
}

interface ChartRow {
  month: string;
  Shipments: number;
  Deliveries: number;
}

export default function ShipmentTrends() {
  const [data, setData] = useState<ChartRow[] | null>(null);

  useEffect(() => {
    gql<{ deliveryReport: DeliveryReport }>(
      `{ deliveryReport(year: 2026) { monthlyData { month totalDeliveries onTime } } }`
    )
      .then((d) =>
        setData(
          d.deliveryReport.monthlyData.map((m) => ({
            month: m.month,
            Shipments: m.totalDeliveries,
            Deliveries: m.onTime,
          }))
        )
      )
      .catch(() => {});
  }, []);

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <div className="mb-4">
        <h3 className="text-lg font-semibold">Shipment Trends</h3>
        <p className="text-sm text-gray-500">
          Monthly shipment and delivery performance over the past year
        </p>
      </div>
      <div className="h-[300px]">
        {data === null ? (
          <div className="flex items-center justify-center h-full">
            <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
          </div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="Shipments" fill="#3b82f6" />
              <Bar dataKey="Deliveries" fill="#22c55e" />
            </BarChart>
          </ResponsiveContainer>
        )}
      </div>
    </div>
  );
}
