"use client";

import { useEffect, useState } from "react";
import { Truck, Plus, Fuel, Calendar, Gauge, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Vehicle { id: string; vehicleId: string; name: string; type: string; status: string; fuelLevel: number; mileage: number; lastService: string; licensePlate: string; }

const statusColor: Record<string, string> = { active: "bg-green-100 text-green-700", maintenance: "bg-amber-100 text-amber-700", available: "bg-blue-100 text-blue-700" };
const typeColor: Record<string, string> = { semi_truck: "bg-slate-100 text-slate-700", box_truck: "bg-purple-100 text-purple-700", flatbed: "bg-orange-100 text-orange-700", van: "bg-cyan-100 text-cyan-700", refrigerated: "bg-blue-100 text-blue-700" };
function fuelBarColor(l: number) { return l >= 70 ? "bg-green-500" : l >= 40 ? "bg-amber-500" : "bg-red-500"; }

export default function VehiclesPage() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ vehicles: { items: Vehicle[] } }>(`{ vehicles(page:1,perPage:50) { items { id vehicleId name type status fuelLevel mileage lastService licensePlate } } }`)
      .then((d) => setVehicles(d.vehicles.items)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="mx-auto max-w-7xl">
        <div className="mb-8 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Vehicle List</h1>
          <button className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-5 py-2.5 text-sm font-medium text-white shadow hover:bg-blue-700 transition"><Plus className="h-4 w-4" /> Add Vehicle</button>
        </div>
        <div className="grid grid-cols-3 gap-4">
          {vehicles.map((v) => (
            <div key={v.id} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm hover:shadow-md transition">
              <div className="mb-4 flex items-start justify-between">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50"><Truck className="h-5 w-5 text-blue-600" /></div>
                  <div><p className="font-semibold text-gray-900">{v.name}</p><p className="text-sm text-gray-500">{v.vehicleId}</p></div>
                </div>
              </div>
              <div className="mb-4 flex items-center gap-2">
                <span className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${typeColor[v.type] || "bg-gray-100 text-gray-700"}`}>{v.type}</span>
                <span className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColor[v.status] || "bg-gray-100 text-gray-700"}`}>{v.status}</span>
              </div>
              <div className="space-y-2 text-sm text-gray-600">
                <div className="flex items-center justify-between"><span className="flex items-center gap-1 text-gray-500"><Gauge className="h-3.5 w-3.5" /> Mileage</span><span className="font-medium text-gray-800">{v.mileage?.toLocaleString()} mi</span></div>
                <div className="flex items-center justify-between"><span className="flex items-center gap-1 text-gray-500"><Calendar className="h-3.5 w-3.5" /> Last Service</span><span className="font-medium text-gray-800">{v.lastService ? new Date(v.lastService).toLocaleDateString() : "N/A"}</span></div>
                <div className="flex items-center justify-between"><span className="text-gray-500">Plate</span><span className="font-medium text-gray-800">{v.licensePlate || "N/A"}</span></div>
              </div>
              <div className="mt-4">
                <div className="mb-1 flex items-center justify-between text-xs"><span className="flex items-center gap-1 text-gray-500"><Fuel className="h-3.5 w-3.5" /> Fuel Level</span><span className="font-medium text-gray-700">{v.fuelLevel}%</span></div>
                <div className="h-2 w-full rounded-full bg-gray-100"><div className={`h-2 rounded-full ${fuelBarColor(v.fuelLevel)}`} style={{ width: `${v.fuelLevel}%` }} /></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
