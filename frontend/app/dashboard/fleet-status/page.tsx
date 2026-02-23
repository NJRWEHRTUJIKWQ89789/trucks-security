"use client";

import { useEffect, useState } from "react";
import { Truck, Activity, Wrench, Loader2, AlertCircle } from "lucide-react";
import { gql } from "@/lib/graphql";

interface DashboardStats {
  totalVehicles: number;
  activeVehicles: number;
}

interface Vehicle {
  id: string;
  vehicleId: string;
  name: string;
  type: string;
  status: string;
  licensePlate: string;
  fuelLevel: number;
  updatedAt: string;
}

interface Driver {
  id: string;
  firstName: string;
  lastName: string;
  vehicleId: string;
}

const statusBadge: Record<string, string> = {
  active: "bg-green-100 text-green-700",
  maintenance: "bg-amber-100 text-amber-700",
  available: "bg-blue-100 text-blue-700",
  inactive: "bg-gray-100 text-gray-700",
};

function formatStatus(status: string) {
  return status.charAt(0).toUpperCase() + status.slice(1).replace(/_/g, " ");
}

function timeAgo(dateStr: string): string {
  if (!dateStr) return "N/A";
  const diff = Date.now() - new Date(dateStr).getTime();
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return "Just now";
  if (minutes < 60) return `${minutes} min ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hr${hours > 1 ? "s" : ""} ago`;
  const days = Math.floor(hours / 24);
  return `${days} day${days > 1 ? "s" : ""} ago`;
}

export default function FleetStatusPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [drivers, setDrivers] = useState<Driver[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const [statsRes, vehiclesRes, driversRes] = await Promise.all([
          gql<{ dashboardStats: DashboardStats }>(
            `{ dashboardStats { totalVehicles activeVehicles } }`
          ),
          gql<{ vehicles: { items: Vehicle[] } }>(
            `{ vehicles(page:1,perPage:50) { items { id vehicleId name type status licensePlate fuelLevel updatedAt } } }`
          ),
          gql<{ drivers: { items: Driver[] } }>(
            `{ drivers(page:1,perPage:100) { items { id firstName lastName vehicleId } } }`
          ),
        ]);
        setStats(statsRes.dashboardStats);
        setVehicles(vehiclesRes.vehicles.items);
        setDrivers(driversRes.drivers.items);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load fleet data");
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

  const totalVehicles = stats?.totalVehicles ?? 0;
  const activeVehicles = stats?.activeVehicles ?? 0;
  const maintenanceCount = vehicles.filter((v) => v.status === "maintenance").length;

  // Build a map from vehicleId (UUID) to driver name for display
  const driverByVehicle: Record<string, string> = {};
  for (const d of drivers) {
    if (d.vehicleId) {
      driverByVehicle[d.vehicleId] = [d.firstName, d.lastName].filter(Boolean).join(" ") || "Unknown";
    }
  }

  const summaryCards = [
    { label: "Total Vehicles", value: totalVehicles, icon: Truck, color: "text-blue-600", bg: "bg-blue-50" },
    { label: "Active", value: activeVehicles, icon: Activity, color: "text-green-600", bg: "bg-green-50" },
    { label: "In Maintenance", value: maintenanceCount, icon: Wrench, color: "text-amber-600", bg: "bg-amber-50" },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Fleet Status</h1>
        <p className="text-gray-500">Overview of all vehicles and their current status</p>
      </div>

      <div className="grid grid-cols-3 gap-4">
        {summaryCards.map((card) => (
          <div key={card.label} className="bg-white rounded-lg border shadow-sm p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">{card.label}</p>
                <p className="text-2xl font-bold mt-1">{card.value}</p>
              </div>
              <div className={`${card.bg} p-3 rounded-lg`}>
                <card.icon className={`h-6 w-6 ${card.color}`} />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="bg-white rounded-lg border shadow-sm p-4">
        <h3 className="text-lg font-semibold mb-4">Vehicle Overview</h3>
        <table className="w-full">
          <thead>
            <tr className="border-b text-left text-sm text-gray-500">
              <th className="pb-3 font-medium">Vehicle ID</th>
              <th className="pb-3 font-medium">Type</th>
              <th className="pb-3 font-medium">Driver</th>
              <th className="pb-3 font-medium">Status</th>
              <th className="pb-3 font-medium">Plate</th>
              <th className="pb-3 font-medium">Last Updated</th>
            </tr>
          </thead>
          <tbody>
            {vehicles.map((v) => (
              <tr key={v.id} className="border-b last:border-0">
                <td className="py-3 font-medium">{v.vehicleId}</td>
                <td className="py-3 text-gray-600">{formatStatus(v.type || "N/A")}</td>
                <td className="py-3 text-gray-600">{driverByVehicle[v.id] || "Unassigned"}</td>
                <td className="py-3">
                  <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium ${statusBadge[v.status] || "bg-gray-100 text-gray-700"}`}>
                    {formatStatus(v.status)}
                  </span>
                </td>
                <td className="py-3 text-gray-600">{v.licensePlate || "N/A"}</td>
                <td className="py-3 text-gray-600">{timeAgo(v.updatedAt)}</td>
              </tr>
            ))}
            {vehicles.length === 0 && (
              <tr>
                <td colSpan={6} className="py-8 text-center text-gray-400">No vehicles found</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
