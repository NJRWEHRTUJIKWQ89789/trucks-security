"use client";

import { useEffect, useState } from "react";
import {
  Map,
  Truck,
  Package,
  CheckCircle,
  AlertTriangle,
  Loader2,
  AlertCircle,
} from "lucide-react";
import { gql } from "@/lib/graphql";

interface DashboardStats {
  activeVehicles: number;
  activeShipments: number;
  deliveredToday: number;
}

interface Shipment {
  id: string;
  trackingNumber: string;
  origin: string;
  destination: string;
  status: string;
  estimatedDelivery: string;
}

const statusStyles: Record<string, string> = {
  in_transit: "bg-blue-100 text-blue-700",
  delivered: "bg-green-100 text-green-700",
  delayed: "bg-red-100 text-red-700",
  pending: "bg-amber-100 text-amber-700",
};

function formatStatus(status: string): string {
  return status
    .split("_")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

function formatDate(dateStr: string): string {
  if (!dateStr) return "N/A";
  const d = new Date(dateStr);
  return d.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });
}

export default function LiveShipmentMapPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const [statsRes, shipmentsRes] = await Promise.all([
          gql<{ dashboardStats: DashboardStats }>(
            `{ dashboardStats { activeVehicles activeShipments deliveredToday } }`
          ),
          gql<{ shipments: { items: Shipment[] } }>(
            `{ shipments(page:1,perPage:20) { items { id trackingNumber origin destination status estimatedDelivery } } }`
          ),
        ]);
        setStats(statsRes.dashboardStats);
        setShipments(shipmentsRes.shipments.items);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load map data");
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-red-600">
        <AlertCircle className="h-8 w-8 mb-2" />
        <p className="text-sm">{error}</p>
      </div>
    );
  }

  const alertCount = shipments.filter((s) => s.status === "delayed").length;

  const statCards = [
    { label: "Active Vehicles", value: stats?.activeVehicles ?? 0, icon: Truck, color: "text-blue-600", bg: "bg-blue-50" },
    { label: "In Transit", value: stats?.activeShipments ?? 0, icon: Package, color: "text-green-600", bg: "bg-green-50" },
    { label: "Delivered Today", value: stats?.deliveredToday ?? 0, icon: CheckCircle, color: "text-purple-600", bg: "bg-purple-50" },
    { label: "Alerts", value: alertCount, icon: AlertTriangle, color: "text-red-600", bg: "bg-red-50" },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Live Shipment Map</h1>
        <p className="text-gray-500">Real-time tracking of all active shipments</p>
      </div>

      <div className="bg-gray-100 rounded-lg h-[500px] flex flex-col items-center justify-center border-2 border-dashed border-gray-300">
        <Map className="h-16 w-16 text-gray-400" />
        <p className="mt-4 text-gray-500 text-lg font-medium">Interactive Map View</p>
      </div>

      <div className="grid grid-cols-4 gap-4">
        {statCards.map((stat) => (
          <div key={stat.label} className="bg-white rounded-lg border shadow-sm p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">{stat.label}</p>
                <p className="text-2xl font-bold mt-1">{stat.value}</p>
              </div>
              <div className={`${stat.bg} p-3 rounded-lg`}>
                <stat.icon className={`h-6 w-6 ${stat.color}`} />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="bg-white rounded-lg border shadow-sm p-4 mt-6">
        <h3 className="text-lg font-semibold mb-4">Recent Movements</h3>
        <table className="w-full">
          <thead>
            <tr className="border-b text-left text-sm text-gray-500">
              <th className="pb-3 font-medium">Shipment ID</th>
              <th className="pb-3 font-medium">Origin</th>
              <th className="pb-3 font-medium">Destination</th>
              <th className="pb-3 font-medium">Status</th>
              <th className="pb-3 font-medium">ETA</th>
            </tr>
          </thead>
          <tbody>
            {shipments.map((s) => (
              <tr key={s.id} className="border-b last:border-0">
                <td className="py-3 font-medium">{s.trackingNumber}</td>
                <td className="py-3 text-gray-600">{s.origin || "N/A"}</td>
                <td className="py-3 text-gray-600">{s.destination || "N/A"}</td>
                <td className="py-3">
                  <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium ${statusStyles[s.status] || "bg-gray-100 text-gray-700"}`}>
                    {formatStatus(s.status)}
                  </span>
                </td>
                <td className="py-3 text-gray-600">{formatDate(s.estimatedDelivery)}</td>
              </tr>
            ))}
            {shipments.length === 0 && (
              <tr>
                <td colSpan={5} className="py-8 text-center text-gray-400">No shipments found</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
