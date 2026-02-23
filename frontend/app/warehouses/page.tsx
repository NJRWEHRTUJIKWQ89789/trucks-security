"use client";

import { useEffect, useState } from "react";
import { Building, Loader2, MapPin, Phone, User } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Warehouse { id: string; name: string; location: string; address: string; capacity: number; usedCapacity: number; manager: string; phone: string; status: string; }

export default function WarehousesPage() {
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ warehouses: { items: Warehouse[] } }>(`{ warehouses(page:1,perPage:50) { items { id name location address capacity usedCapacity manager phone status } } }`)
      .then((d) => setWarehouses(d.warehouses.items)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Warehouse Locations</h1>
      <div className="grid grid-cols-2 gap-6">
        {warehouses.map((w) => {
          const pct = w.capacity > 0 ? Math.round((w.usedCapacity / w.capacity) * 100) : 0;
          return (
            <div key={w.id} className="rounded-xl border bg-white p-6 shadow-sm">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3"><div className="h-10 w-10 rounded-lg bg-blue-50 flex items-center justify-center"><Building className="h-5 w-5 text-blue-600" /></div><div><h3 className="font-semibold text-gray-900">{w.name}</h3><span className={`text-xs px-2 py-0.5 rounded-full ${w.status === "active" ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-700"}`}>{w.status}</span></div></div>
              </div>
              <div className="space-y-2 text-sm text-gray-600">
                <div className="flex items-center gap-2"><MapPin className="h-3.5 w-3.5" />{w.location}</div>
                <div className="flex items-center gap-2"><User className="h-3.5 w-3.5" />{w.manager}</div>
                <div className="flex items-center gap-2"><Phone className="h-3.5 w-3.5" />{w.phone}</div>
              </div>
              <div className="mt-4">
                <div className="flex justify-between text-xs mb-1"><span className="text-gray-500">Capacity</span><span className="font-medium">{w.usedCapacity?.toLocaleString()} / {w.capacity?.toLocaleString()}</span></div>
                <div className="h-2 bg-gray-100 rounded-full"><div className={`h-2 rounded-full ${pct > 85 ? "bg-red-500" : pct > 60 ? "bg-amber-500" : "bg-green-500"}`} style={{ width: `${pct}%` }} /></div>
                <p className="text-xs text-gray-400 mt-1">{pct}% utilized</p>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
