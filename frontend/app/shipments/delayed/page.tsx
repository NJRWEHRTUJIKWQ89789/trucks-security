"use client";

import { useEffect, useState } from "react";
import { AlertTriangle, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Shipment {
  id: string;
  trackingNumber: string;
  origin: string;
  destination: string;
  status: string;
  carrier: string;
  weight: number;
  estimatedDelivery: string;
  customerName: string;
  notes: string;
  createdAt: string;
}

function formatDate(dateStr: string | null | undefined): string {
  if (!dateStr) return "-";
  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return dateStr;
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

function getDaysOverdue(estimatedDelivery: string | null | undefined): number {
  if (!estimatedDelivery) return 0;
  const eta = new Date(estimatedDelivery);
  if (isNaN(eta.getTime())) return 0;
  const now = new Date();
  const diffMs = now.getTime() - eta.getTime();
  return Math.max(0, Math.ceil(diffMs / (1000 * 60 * 60 * 24)));
}

function getPriority(daysOverdue: number): string {
  if (daysOverdue >= 7) return "High";
  if (daysOverdue >= 3) return "Medium";
  return "Low";
}

const priorityColors: Record<string, string> = {
  High: "bg-red-100 text-red-700",
  Medium: "bg-amber-100 text-amber-700",
  Low: "bg-gray-100 text-gray-700",
};

export default function DelayedShipmentsPage() {
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    gql<{ delayedShipments: Shipment[] }>(
      `{delayedShipments{id trackingNumber origin destination status carrier weight estimatedDelivery customerName notes createdAt}}`
    )
      .then((d) => setShipments(d.delayedShipments))
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load delayed shipments"))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="flex items-center gap-3 mb-8">
        <AlertTriangle className="w-7 h-7 text-red-500" />
        <h1 className="text-2xl font-bold text-gray-900">Delayed Shipments</h1>
        <span className="bg-red-100 text-red-700 text-sm font-medium px-3 py-0.5 rounded-full">
          {shipments.length} shipment{shipments.length !== 1 ? "s" : ""}
        </span>
      </div>

      {error && (
        <div className="flex items-center gap-3 bg-red-50 border border-red-200 rounded-xl p-4 text-red-700 mb-6">
          <AlertTriangle className="w-5 h-5 flex-shrink-0" />
          <p className="text-sm">{error}</p>
        </div>
      )}

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50 border-b border-gray-200">
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Shipment ID</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Route</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Original ETA</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Days Overdue</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Carrier</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Customer</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Priority</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {shipments.map((s) => {
              const daysOverdue = getDaysOverdue(s.estimatedDelivery);
              const priority = getPriority(daysOverdue);
              return (
                <tr key={s.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 text-sm font-medium text-blue-600">{s.trackingNumber}</td>
                  <td className="px-6 py-4 text-sm text-gray-700">
                    {s.origin && s.destination ? `${s.origin} \u2192 ${s.destination}` : s.origin || s.destination || "-"}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-700">{formatDate(s.estimatedDelivery)}</td>
                  <td className="px-6 py-4 text-sm font-medium text-red-600">
                    {daysOverdue > 0 ? `${daysOverdue} day${daysOverdue !== 1 ? "s" : ""}` : "-"}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-700">{s.carrier || "-"}</td>
                  <td className="px-6 py-4 text-sm text-gray-700">{s.customerName || "-"}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium ${priorityColors[priority]}`}>
                      {priority}
                    </span>
                  </td>
                </tr>
              );
            })}
            {shipments.length === 0 && !error && (
              <tr>
                <td colSpan={7} className="px-6 py-8 text-center text-gray-400">
                  No delayed shipments found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
