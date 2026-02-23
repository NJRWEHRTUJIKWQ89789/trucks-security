"use client";

import { useEffect, useState } from "react";
import { Users, Loader2, Star, Phone, Mail } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Driver { id: string; employeeId: string; firstName: string; lastName: string; email: string; phone: string; status: string; rating: number; totalDeliveries: number; licenseNumber: string; }

const statusColor: Record<string, string> = { available: "bg-green-100 text-green-700", on_delivery: "bg-blue-100 text-blue-700", off_duty: "bg-gray-100 text-gray-700", on_leave: "bg-amber-100 text-amber-700" };

export default function DriversPage() {
  const [drivers, setDrivers] = useState<Driver[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ drivers: { items: Driver[] } }>(`{ drivers(page:1,perPage:50) { items { id employeeId firstName lastName email phone status rating totalDeliveries licenseNumber } } }`)
      .then((d) => setDrivers(d.drivers.items)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Driver Assignments</h1>
      </div>
      <div className="grid grid-cols-3 gap-4">
        {drivers.map((d) => (
          <div key={d.id} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3 mb-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-100"><Users className="h-5 w-5 text-blue-600" /></div>
              <div><p className="font-semibold text-gray-900">{d.firstName} {d.lastName}</p><p className="text-sm text-gray-500">{d.employeeId}</p></div>
              <span className={`ml-auto inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColor[d.status] || "bg-gray-100 text-gray-700"}`}>{d.status}</span>
            </div>
            <div className="space-y-1.5 text-sm">
              <div className="flex items-center gap-2 text-gray-600"><Mail className="h-3.5 w-3.5" />{d.email}</div>
              <div className="flex items-center gap-2 text-gray-600"><Phone className="h-3.5 w-3.5" />{d.phone}</div>
              <div className="flex items-center justify-between mt-2 pt-2 border-t">
                <span className="flex items-center gap-1 text-amber-500"><Star className="h-4 w-4 fill-current" />{d.rating?.toFixed(1)}</span>
                <span className="text-gray-500">{d.totalDeliveries} deliveries</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
